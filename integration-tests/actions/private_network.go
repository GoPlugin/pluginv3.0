package actions

import (
	"github.com/rs/zerolog"

	ctf_test_env "github.com/goplugin/plugin-testing-framework/docker/test_env"
	tc "github.com/goplugin/pluginv3.0/integration-tests/testconfig"
)

func EthereumNetworkConfigFromConfig(l zerolog.Logger, config tc.GlobalTestConfig) (network ctf_test_env.EthereumNetwork, err error) {
	if config.GetPrivateEthereumNetworkConfig() == nil {
		l.Warn().Msg("No TOML private ethereum network config found, will use old geth")
		ethBuilder := ctf_test_env.NewEthereumNetworkBuilder()
		network, err = ethBuilder.
			WithConsensusType(ctf_test_env.ConsensusType_PoW).
			WithExecutionLayer(ctf_test_env.ExecutionLayer_Geth).
			Build()

		return
	}

	ethBuilder := ctf_test_env.NewEthereumNetworkBuilder()
	network, err = ethBuilder.
		WithExistingConfig(*config.GetPrivateEthereumNetworkConfig()).
		Build()

	return
}
