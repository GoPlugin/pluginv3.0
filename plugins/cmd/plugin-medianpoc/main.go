package main

import (
	"github.com/hashicorp/go-plugin"

	"github.com/goplugin/plugin-common/pkg/loop"
	"github.com/goplugin/plugin-common/pkg/loop/reportingplugins"
	"github.com/goplugin/plugin-common/pkg/types"
	"github.com/goplugin/pluginv3.0/v2/plugins/medianpoc"
)

const (
	loggerName = "PluginMedianPoc"
)

func main() {
	s := loop.MustNewStartedServer(loggerName)
	defer s.Stop()

	p := medianpoc.NewPlugin(s.Logger)
	defer s.Logger.ErrorIfFn(p.Close, "Failed to close")

	s.MustRegister(p)

	stop := make(chan struct{})
	defer close(stop)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: reportingplugins.ReportingPluginHandshakeConfig(),
		Plugins: map[string]plugin.Plugin{
			reportingplugins.PluginServiceName: &reportingplugins.GRPCService[types.MedianProvider]{
				PluginServer: p,
				BrokerConfig: loop.BrokerConfig{
					Logger:   s.Logger,
					StopCh:   stop,
					GRPCOpts: s.GRPCOpts,
				},
			},
		},
		GRPCServer: s.GRPCOpts.NewServer,
	})
}
