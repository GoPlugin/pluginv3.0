package ocr

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"

	"github.com/goplugin/plugin-testing-framework/blockchain"

	"github.com/goplugin/wasp"
	"go.uber.org/ratelimit"

	client2 "github.com/goplugin/plugin-testing-framework/client"

	"github.com/goplugin/pluginv3.0/integration-tests/actions"
	"github.com/goplugin/pluginv3.0/integration-tests/client"
	"github.com/goplugin/pluginv3.0/integration-tests/contracts"
)

// VU is a virtual user for the OCR load test
// it creates a feed and triggers new rounds
type VU struct {
	*wasp.VUControl
	rl            ratelimit.Limiter
	rate          int
	rateUnit      time.Duration
	roundNum      atomic.Int64
	cc            blockchain.EVMClient
	lt            contracts.LinkToken
	cd            contracts.ContractDeployer
	bootstrapNode *client.PluginK8sClient
	workerNodes   []*client.PluginK8sClient
	msClient      *client2.MockserverClient
	l             zerolog.Logger
	ocrInstances  []contracts.OffchainAggregator
}

func NewVU(
	l zerolog.Logger,
	rate int,
	rateUnit time.Duration,
	cc blockchain.EVMClient,
	lt contracts.LinkToken,
	cd contracts.ContractDeployer,
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
		cc:            cc,
		lt:            lt,
		cd:            cd,
		msClient:      msClient,
		bootstrapNode: bootstrapNode,
		workerNodes:   workerNodes,
	}
}

func (m *VU) Clone(_ *wasp.Generator) wasp.VirtualUser {
	return &VU{
		VUControl:     wasp.NewVUControl(),
		rl:            ratelimit.New(m.rate, ratelimit.Per(m.rateUnit)),
		rate:          m.rate,
		rateUnit:      m.rateUnit,
		l:             m.l,
		cc:            m.cc,
		lt:            m.lt,
		cd:            m.cd,
		msClient:      m.msClient,
		bootstrapNode: m.bootstrapNode,
		workerNodes:   m.workerNodes,
	}
}

func (m *VU) Setup(_ *wasp.Generator) error {
	ocrInstances, err := actions.DeployOCRContracts(1, m.lt, m.cd, m.workerNodes, m.cc)
	if err != nil {
		return err
	}
	err = actions.CreateOCRJobs(ocrInstances, m.bootstrapNode, m.workerNodes, 5, m.msClient, m.cc.GetChainID().String())
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