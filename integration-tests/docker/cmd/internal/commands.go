package internal

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	ctf_test_env "github.com/goplugin/plugin-testing-framework/lib/docker/test_env"
	"github.com/goplugin/plugin-testing-framework/lib/logging"
	"github.com/goplugin/pluginv3.0/integration-tests/docker/test_env"
	"github.com/goplugin/pluginv3.0/integration-tests/testconfig"
)

var NodeCountFlag = "node-count"

var StartNodesCmd = &cobra.Command{
	Use:   "start-nodes",
	Short: "Start Plugin nodes",
	RunE: func(cmd *cobra.Command, _ []string) error {
		nodeCount, err := cmd.Flags().GetInt(NodeCountFlag)
		if err != nil {
			return err
		}

		log.Logger = logging.GetLogger(nil, "CORE_DOCKER_ENV_LOG_LEVEL")
		log.Info().Msg("Starting docker test env with Plugin nodes..")

		config, err := testconfig.GetConfig([]string{"Smoke"}, testconfig.OCR2)
		if err != nil {
			return err
		}

		ethBuilder := ctf_test_env.NewEthereumNetworkBuilder()
		network, err := ethBuilder.
			WithExistingConfig(*config.GetPrivateEthereumNetworkConfig()).
			Build()
		if err != nil {
			return err
		}

		_, err = test_env.NewCLTestEnvBuilder().
			WithTestConfig(&config).
			WithPrivateEthereumNetwork(network.EthereumNetworkConfig).
			WithMockAdapter().
			WithCLNodes(nodeCount).
			WithoutCleanup().
			Build()
		if err != nil {
			return err
		}

		log.Info().Msg("Test env is ready")

		handleExitSignal()

		return nil
	},
}

func init() {
	StartNodesCmd.PersistentFlags().Int(
		NodeCountFlag,
		6,
		"Number of Plugin nodes",
	)
}

func handleExitSignal() {
	// Create a channel to receive exit signals
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, syscall.SIGTERM)

	log.Info().Msg("Press Ctrl+C to destroy the test environment")

	// Block until an exit signal is received
	<-exitChan
}
