package dummy

import (
	"context"
	"errors"

	ocrtypes "github.com/goplugin/plugin-libocr/offchainreporting2plus/types"

	"github.com/goplugin/plugin-common/pkg/services"
	commontypes "github.com/goplugin/plugin-common/pkg/types"
	relaytypes "github.com/goplugin/plugin-common/pkg/types"
	llotypes "github.com/goplugin/plugin-common/pkg/types/llo"
	"github.com/goplugin/pluginv3.0/v2/core/logger"
)

var _ commontypes.LLOProvider = (*lloProvider)(nil)

type Transmitter interface {
	services.Service
	llotypes.Transmitter
}

type lloProvider struct {
	cp                     commontypes.ConfigProvider
	transmitter            Transmitter
	logger                 logger.Logger
	channelDefinitionCache llotypes.ChannelDefinitionCache

	ms services.MultiStart
}

func NewLLOProvider(
	lggr logger.Logger,
	cp commontypes.ConfigProvider,
	transmitter Transmitter,
	channelDefinitionCache llotypes.ChannelDefinitionCache,
) relaytypes.LLOProvider {
	return &lloProvider{
		cp,
		transmitter,
		lggr.Named("LLOProvider"),
		channelDefinitionCache,
		services.MultiStart{},
	}
}

func (p *lloProvider) Start(ctx context.Context) error {
	err := p.ms.Start(ctx, p.cp, p.transmitter, p.channelDefinitionCache)
	return err
}

func (p *lloProvider) Close() error {
	return p.ms.Close()
}

func (p *lloProvider) Ready() error {
	return errors.Join(p.cp.Ready(), p.transmitter.Ready(), p.channelDefinitionCache.Ready())
}

func (p *lloProvider) Name() string {
	return p.logger.Name()
}

func (p *lloProvider) HealthReport() map[string]error {
	report := map[string]error{p.Name(): nil}
	services.CopyHealth(report, p.cp.HealthReport())
	services.CopyHealth(report, p.transmitter.HealthReport())
	services.CopyHealth(report, p.channelDefinitionCache.HealthReport())
	return report
}

func (p *lloProvider) ContractConfigTracker() ocrtypes.ContractConfigTracker {
	return p.cp.ContractConfigTracker()
}

func (p *lloProvider) OffchainConfigDigester() ocrtypes.OffchainConfigDigester {
	return p.cp.OffchainConfigDigester()
}

func (p *lloProvider) ContractTransmitter() llotypes.Transmitter {
	return p.transmitter
}

func (p *lloProvider) ChannelDefinitionCache() llotypes.ChannelDefinitionCache {
	return p.channelDefinitionCache
}
