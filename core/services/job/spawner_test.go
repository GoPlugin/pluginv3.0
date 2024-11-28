package job_test

import (
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/jmoiron/sqlx"

	"github.com/goplugin/plugin-common/pkg/loop"
	"github.com/goplugin/plugin-common/pkg/services"
	"github.com/goplugin/plugin-common/pkg/services/servicetest"
	"github.com/goplugin/plugin-common/pkg/utils"
	"github.com/goplugin/plugin-common/pkg/utils/mailbox/mailboxtest"
	"github.com/goplugin/pluginv3.0/v2/core/capabilities"

	"github.com/goplugin/pluginv3.0/v2/core/bridges"
	mocklp "github.com/goplugin/pluginv3.0/v2/core/chains/evm/logpoller/mocks"
	evmtypes "github.com/goplugin/pluginv3.0/v2/core/chains/evm/types"
	"github.com/goplugin/pluginv3.0/v2/core/internal/cltest"
	"github.com/goplugin/pluginv3.0/v2/core/internal/testutils"
	"github.com/goplugin/pluginv3.0/v2/core/internal/testutils/configtest"
	"github.com/goplugin/pluginv3.0/v2/core/internal/testutils/evmtest"
	"github.com/goplugin/pluginv3.0/v2/core/internal/testutils/pgtest"
	"github.com/goplugin/pluginv3.0/v2/core/logger"
	"github.com/goplugin/pluginv3.0/v2/core/services/plugin"
	"github.com/goplugin/pluginv3.0/v2/core/services/job"
	"github.com/goplugin/pluginv3.0/v2/core/services/job/mocks"
	"github.com/goplugin/pluginv3.0/v2/core/services/ocr"
	"github.com/goplugin/pluginv3.0/v2/core/services/ocr2"
	"github.com/goplugin/pluginv3.0/v2/core/services/pipeline"
	"github.com/goplugin/pluginv3.0/v2/core/services/relay"
	evmrelay "github.com/goplugin/pluginv3.0/v2/core/services/relay/evm"
	evmrelayer "github.com/goplugin/pluginv3.0/v2/core/services/relay/evm"
	"github.com/goplugin/pluginv3.0/v2/plugins"
)

type delegate struct {
	jobType                    job.Type
	services                   []job.ServiceCtx
	jobID                      int32
	chContinueCreatingServices chan struct{}
	job.Delegate
}

func (d delegate) JobType() job.Type {
	return d.jobType
}

// ServicesForSpec satisfies the job.Delegate interface.
func (d delegate) ServicesForSpec(js job.Job) ([]job.ServiceCtx, error) {
	if js.Type != d.jobType {
		return nil, nil
	}
	return d.services, nil
}

func clearDB(t *testing.T, db *sqlx.DB) {
	cltest.ClearDBTables(t, db, "jobs", "pipeline_runs", "pipeline_specs", "pipeline_task_runs")
}

type relayGetter struct {
	e evmrelay.EVMChainRelayerExtender
	r *evmrelayer.Relayer
}

func (g *relayGetter) Get(id relay.ID) (loop.Relayer, error) {
	return evmrelayer.NewLoopRelayServerAdapter(g.r, g.e), nil
}

func TestSpawner_CreateJobDeleteJob(t *testing.T) {
	t.Parallel()

	config := configtest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, config.Database())
	ethKeyStore := keyStore.Eth()
	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))
	require.NoError(t, keyStore.P2P().Add(cltest.DefaultP2PKey))
	require.NoError(t, keyStore.OCR2().Add(cltest.DefaultOCR2Key))

	_, address := cltest.MustInsertRandomKey(t, ethKeyStore)
	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config.Database())
	_, bridge2 := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config.Database())

	ethClient := cltest.NewEthMocksWithDefaultChain(t)
	ethClient.On("CallContext", mock.Anything, mock.Anything, "eth_getBlockByNumber", mock.Anything, false).
		Run(func(args mock.Arguments) {
			head := args.Get(1).(**evmtypes.Head)
			*head = cltest.Head(10)
		}).
		Return(nil).Maybe()

	relayExtenders := evmtest.NewChainRelayExtenders(t, evmtest.TestChainOpts{DB: db, Client: ethClient, GeneralConfig: config, KeyStore: ethKeyStore})
	legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)
	t.Run("should respect its dependents", func(t *testing.T) {
		lggr := logger.TestLogger(t)
		orm := NewTestORM(t, db, pipeline.NewORM(db, lggr, config.Database(), config.JobPipeline().MaxSuccessfulRuns()), bridges.NewORM(db, lggr, config.Database()), keyStore, config.Database())
		a := utils.NewDependentAwaiter()
		a.AddDependents(1)
		spawner := job.NewSpawner(orm, config.Database(), noopChecker{}, map[job.Type]job.Delegate{}, db, lggr, []utils.DependentAwaiter{a})
		// Starting the spawner should signal to the dependents
		result := make(chan bool)
		go func() {
			select {
			case <-a.AwaitDependents():
				result <- true
			case <-time.After(2 * time.Second):
				result <- false
			}
		}()
		require.NoError(t, spawner.Start(testutils.Context(t)))
		assert.True(t, <-result, "failed to signal to dependents")
	})

	t.Run("starts and stops job services when jobs are added and removed", func(t *testing.T) {
		jobA := cltest.MakeDirectRequestJobSpec(t)
		jobB := makeOCRJobSpec(t, address, bridge.Name.String(), bridge2.Name.String())

		lggr := logger.TestLogger(t)
		orm := NewTestORM(t, db, pipeline.NewORM(db, lggr, config.Database(), config.JobPipeline().MaxSuccessfulRuns()), bridges.NewORM(db, lggr, config.Database()), keyStore, config.Database())

		eventuallyA := cltest.NewAwaiter()
		serviceA1 := mocks.NewServiceCtx(t)
		serviceA2 := mocks.NewServiceCtx(t)
		serviceA1.On("Start", mock.Anything).Return(nil).Once()
		serviceA2.On("Start", mock.Anything).Return(nil).Once().Run(func(mock.Arguments) { eventuallyA.ItHappened() })
		mailMon := servicetest.Run(t, mailboxtest.NewMonitor(t))
		dA := ocr.NewDelegate(nil, orm, nil, nil, nil, monitoringEndpoint, legacyChains, logger.TestLogger(t), config.Database(), mailMon)
		delegateA := &delegate{jobA.Type, []job.ServiceCtx{serviceA1, serviceA2}, 0, make(chan struct{}), dA}

		eventuallyB := cltest.NewAwaiter()
		serviceB1 := mocks.NewServiceCtx(t)
		serviceB2 := mocks.NewServiceCtx(t)
		serviceB1.On("Start", mock.Anything).Return(nil).Once()
		serviceB2.On("Start", mock.Anything).Return(nil).Once().Run(func(mock.Arguments) { eventuallyB.ItHappened() })
		dB := ocr.NewDelegate(nil, orm, nil, nil, nil, monitoringEndpoint, legacyChains, logger.TestLogger(t), config.Database(), mailMon)
		delegateB := &delegate{jobB.Type, []job.ServiceCtx{serviceB1, serviceB2}, 0, make(chan struct{}), dB}

		spawner := job.NewSpawner(orm, config.Database(), noopChecker{}, map[job.Type]job.Delegate{
			jobA.Type: delegateA,
			jobB.Type: delegateB,
		}, db, lggr, nil)
		require.NoError(t, spawner.Start(testutils.Context(t)))
		err := spawner.CreateJob(jobA)
		require.NoError(t, err)
		jobSpecIDA := jobA.ID
		delegateA.jobID = jobSpecIDA
		close(delegateA.chContinueCreatingServices)

		eventuallyA.AwaitOrFail(t, 20*time.Second)

		err = spawner.CreateJob(jobB)
		require.NoError(t, err)
		jobSpecIDB := jobB.ID
		delegateB.jobID = jobSpecIDB
		close(delegateB.chContinueCreatingServices)

		eventuallyB.AwaitOrFail(t, 20*time.Second)

		serviceA1.On("Close").Return(nil).Once()
		serviceA2.On("Close").Return(nil).Once()
		err = spawner.DeleteJob(jobSpecIDA)
		require.NoError(t, err)

		serviceB1.On("Close").Return(nil).Once()
		serviceB2.On("Close").Return(nil).Once()
		err = spawner.DeleteJob(jobSpecIDB)
		require.NoError(t, err)

		require.NoError(t, spawner.Close())
	})

	clearDB(t, db)

	t.Run("starts and stops job services from the DB when .Start()/.Stop() is called", func(t *testing.T) {
		jobA := makeOCRJobSpec(t, address, bridge.Name.String(), bridge2.Name.String())

		eventually := cltest.NewAwaiter()
		serviceA1 := mocks.NewServiceCtx(t)
		serviceA2 := mocks.NewServiceCtx(t)
		serviceA1.On("Start", mock.Anything).Return(nil).Once()
		serviceA2.On("Start", mock.Anything).Return(nil).Once().Run(func(mock.Arguments) { eventually.ItHappened() })

		lggr := logger.TestLogger(t)
		orm := NewTestORM(t, db, pipeline.NewORM(db, lggr, config.Database(), config.JobPipeline().MaxSuccessfulRuns()), bridges.NewORM(db, lggr, config.Database()), keyStore, config.Database())
		mailMon := servicetest.Run(t, mailboxtest.NewMonitor(t))
		d := ocr.NewDelegate(nil, orm, nil, nil, nil, monitoringEndpoint, legacyChains, logger.TestLogger(t), config.Database(), mailMon)
		delegateA := &delegate{jobA.Type, []job.ServiceCtx{serviceA1, serviceA2}, 0, nil, d}
		spawner := job.NewSpawner(orm, config.Database(), noopChecker{}, map[job.Type]job.Delegate{
			jobA.Type: delegateA,
		}, db, lggr, nil)

		err := orm.CreateJob(jobA)
		require.NoError(t, err)
		delegateA.jobID = jobA.ID

		require.NoError(t, spawner.Start(testutils.Context(t)))

		eventually.AwaitOrFail(t)

		serviceA1.On("Close").Return(nil).Once()
		serviceA2.On("Close").Return(nil).Once()

		require.NoError(t, spawner.Close())
	})

	clearDB(t, db)

	t.Run("closes job services on 'DeleteJob()'", func(t *testing.T) {
		jobA := makeOCRJobSpec(t, address, bridge.Name.String(), bridge2.Name.String())

		eventuallyStart := cltest.NewAwaiter()
		serviceA1 := mocks.NewServiceCtx(t)
		serviceA2 := mocks.NewServiceCtx(t)
		serviceA1.On("Start", mock.Anything).Return(nil).Once()
		serviceA2.On("Start", mock.Anything).Return(nil).Once().Run(func(mock.Arguments) { eventuallyStart.ItHappened() })

		lggr := logger.TestLogger(t)
		orm := NewTestORM(t, db, pipeline.NewORM(db, lggr, config.Database(), config.JobPipeline().MaxSuccessfulRuns()), bridges.NewORM(db, lggr, config.Database()), keyStore, config.Database())
		mailMon := servicetest.Run(t, mailboxtest.NewMonitor(t))
		d := ocr.NewDelegate(nil, orm, nil, nil, nil, monitoringEndpoint, legacyChains, logger.TestLogger(t), config.Database(), mailMon)
		delegateA := &delegate{jobA.Type, []job.ServiceCtx{serviceA1, serviceA2}, 0, nil, d}
		spawner := job.NewSpawner(orm, config.Database(), noopChecker{}, map[job.Type]job.Delegate{
			jobA.Type: delegateA,
		}, db, lggr, nil)

		err := orm.CreateJob(jobA)
		require.NoError(t, err)
		jobSpecIDA := jobA.ID
		delegateA.jobID = jobSpecIDA

		require.NoError(t, spawner.Start(testutils.Context(t)))
		defer func() { assert.NoError(t, spawner.Close()) }()

		eventuallyStart.AwaitOrFail(t)

		// Wait for the claim lock to be taken
		gomega.NewWithT(t).Eventually(func() bool {
			jobs := spawner.ActiveJobs()
			_, exists := jobs[jobSpecIDA]
			return exists
		}, testutils.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal(true))

		eventuallyClose := cltest.NewAwaiter()
		serviceA1.On("Close").Return(nil).Once()
		serviceA2.On("Close").Return(nil).Once().Run(func(mock.Arguments) { eventuallyClose.ItHappened() })

		err = spawner.DeleteJob(jobSpecIDA)
		require.NoError(t, err)

		eventuallyClose.AwaitOrFail(t)

		// Wait for the claim lock to be released
		gomega.NewWithT(t).Eventually(func() bool {
			jobs := spawner.ActiveJobs()
			_, exists := jobs[jobSpecIDA]
			return exists
		}, testutils.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal(false))

		clearDB(t, db)
	})

	t.Run("Unregisters filters on 'DeleteJob()'", func(t *testing.T) {
		config = configtest.NewGeneralConfig(t, func(c *plugin.Config, s *plugin.Secrets) {
			c.Feature.LogPoller = func(b bool) *bool { return &b }(true)
		})
		lp := &mocklp.LogPoller{}
		testopts := evmtest.TestChainOpts{
			DB:            db,
			Client:        ethClient,
			GeneralConfig: config,
			LogPoller:     lp,
			KeyStore:      ethKeyStore,
		}

		lggr := logger.TestLogger(t)
		relayExtenders := evmtest.NewChainRelayExtenders(t, testopts)
		assert.Equal(t, relayExtenders.Len(), 1)
		legacyChains := evmrelay.NewLegacyChainsFromRelayerExtenders(relayExtenders)
		chain := evmtest.MustGetDefaultChain(t, legacyChains)

		evmRelayer, err := evmrelayer.NewRelayer(lggr, chain, evmrelayer.RelayerOpts{
			DB:             db,
			QConfig:        testopts.GeneralConfig.Database(),
			CSAETHKeystore: keyStore,
		})
		assert.NoError(t, err)

		testRelayGetter := &relayGetter{
			e: relayExtenders.Slice()[0],
			r: evmRelayer,
		}

		jobOCR2VRF := makeOCR2VRFJobSpec(t, keyStore, config, address, chain.ID(), 2)

		orm := NewTestORM(t, db, pipeline.NewORM(db, lggr, config.Database(), config.JobPipeline().MaxSuccessfulRuns()), bridges.NewORM(db, lggr, config.Database()), keyStore, config.Database())
		mailMon := servicetest.Run(t, mailboxtest.NewMonitor(t))

		processConfig := plugins.NewRegistrarConfig(loop.GRPCOpts{}, func(name string) (*plugins.RegisteredLoop, error) { return nil, nil })
		ocr2DelegateConfig := ocr2.NewDelegateConfig(config.OCR2(), config.Mercury(), config.Threshold(), config.Insecure(), config.JobPipeline(), config.Database(), processConfig)

		d := ocr2.NewDelegate(nil, orm, nil, nil, nil, nil, monitoringEndpoint, legacyChains, lggr, ocr2DelegateConfig,
			keyStore.OCR2(), keyStore.DKGSign(), keyStore.DKGEncrypt(), ethKeyStore, testRelayGetter, mailMon, capabilities.NewRegistry())
		delegateOCR2 := &delegate{jobOCR2VRF.Type, []job.ServiceCtx{}, 0, nil, d}

		spawner := job.NewSpawner(orm, config.Database(), noopChecker{}, map[job.Type]job.Delegate{
			jobOCR2VRF.Type: delegateOCR2,
		}, db, lggr, nil)

		err = spawner.CreateJob(jobOCR2VRF)
		require.NoError(t, err)
		jobSpecID := jobOCR2VRF.ID
		delegateOCR2.jobID = jobOCR2VRF.ID

		lp.On("UnregisterFilter", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			lggr.Debugf("Got here, with args %v", args)
		})

		err = spawner.DeleteJob(jobSpecID)
		require.NoError(t, err)

		lp.AssertNumberOfCalls(t, "UnregisterFilter", 3)

		lp.On("Close").Return(nil).Once()
		spawner.Close()
	})
}

type noopChecker struct{}

func (n noopChecker) Register(service services.HealthReporter) error { return nil }

func (n noopChecker) Unregister(name string) error { return nil }

func (n noopChecker) IsReady() (ready bool, errors map[string]error) { return true, nil }

func (n noopChecker) IsHealthy() (healthy bool, errors map[string]error) { return true, nil }

func (n noopChecker) Start() error { return nil }

func (n noopChecker) Close() error { return nil }
