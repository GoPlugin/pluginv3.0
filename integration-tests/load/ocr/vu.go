package ocr

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"go.uber.org/ratelimit"

	client2 "github.com/goplugin/plugin-testing-framework/lib/client"
	"github.com/goplugin/plugin-testing-framework/seth"
	"github.com/goplugin/plugin-testing-framework/wasp"

	"github.com/goplugin/pluginv3.0/integration-tests/actions"
	"github.com/goplugin/pluginv3.0/integration-tests/client"
	"github.com/goplugin/pluginv3.0/integration-tests/contracts"
	"github.com/goplugin/pluginv3.0/integration-tests/testconfig/ocr"
)

// VU is a virtual user for the OCR load test
// it creates a feed and triggers new rounds
type VU struct {
	*wasp.VUControl
	rl            ratelimit.Limiter
	rate          int
	rateUnit      time.Duration
	roundNum      atomic.Int64
	seth          *seth.Client
	lta           common.Address
	bootstrapNode *client.PluginK8sClient
	workerNodes   []*client.PluginK8sClient
	msClient      *client2.MockserverClient
	l             zerolog.Logger
	ocrInstances  []contracts.OffchainAggregator
	config        ocr.OffChainAggregatorsConfig
}

func NewVU(
	l zerolog.Logger,
	seth *seth.Client,
	config ocr.OffChainAggregatorsConfig,
	rate int,
	rateUnit time.Duration,
	lta common.Address,
	bootstrapNode *client.PluginK8sClient,
	workerNodes []*client.PluginK8sClient,
	msClient *client2.MockserverClient,
) *VU {
	return &VU{
		VUControl:     wasp.NewVUControl(),
		rl:            ratelimit.New(rate, ratelimit.Per(rateUnit)),
		rate:          rate,
		rateUnit:      rateUnit,
		l:             l,
		seth:          seth,
		lta:           lta,
		msClient:      msClient,
		bootstrapNode: bootstrapNode,
		workerNodes:   workerNodes,
		config:        config,
	}
}

func (m *VU) Clone(_ *wasp.Generator) wasp.VirtualUser {
	return &VU{
		VUControl:     wasp.NewVUControl(),
		rl:            ratelimit.New(m.rate, ratelimit.Per(m.rateUnit)),
		rate:          m.rate,
		rateUnit:      m.rateUnit,
		l:             m.l,
		seth:          m.seth,
		lta:           m.lta,
		msClient:      m.msClient,
		bootstrapNode: m.bootstrapNode,
		workerNodes:   m.workerNodes,
		config:        m.config,
	}
}

func (m *VU) Setup(_ *wasp.Generator) error {
	ocrInstances, err := actions.SetupOCRv1Contracts(m.l, m.seth, m.config, m.lta, contracts.PluginK8sClientToPluginNodeWithKeysAndAddress(m.workerNodes))
	if err != nil {
		return err
	}
	err = actions.CreateOCRJobs(ocrInstances, m.bootstrapNode, m.workerNodes, 5, m.msClient, fmt.Sprint(m.seth.ChainID))
	if err != nil {
		return err
	}
	m.ocrInstances = ocrInstances
	return nil
}

func (m *VU) Teardown(_ *wasp.Generator) error {
	return nil
}

func (m *VU) Call(l *wasp.Generator) {
	m.rl.Take()
	m.roundNum.Add(1)
	requestedRound := m.roundNum.Load()
	m.l.Info().
		Int64("RoundNum", requestedRound).
		Str("FeedID", m.ocrInstances[0].Address()).
		Msg("starting new round")
	err := m.ocrInstances[0].RequestNewRound()
	if err != nil {
		l.ResponsesChan <- &wasp.Response{Error: err.Error(), Failed: true}
	}
	for {
		time.Sleep(5 * time.Second)
		lr, err := m.ocrInstances[0].GetLatestRound(context.Background())
		if err != nil {
			l.ResponsesChan <- &wasp.Response{Error: err.Error(), Failed: true}
		}
		m.l.Info().Interface("LatestRound", lr).Msg("latest round")
		if lr.RoundId.Int64() >= requestedRound {
			l.ResponsesChan <- &wasp.Response{}
		}
	}
}
