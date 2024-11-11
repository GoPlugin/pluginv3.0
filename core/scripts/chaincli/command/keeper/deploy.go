package keeper

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/goplugin/pluginv3.0/core/scripts/chaincli/config"
	"github.com/goplugin/pluginv3.0/core/scripts/chaincli/handler"
)

// deployCmd represents the command to run the service
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy keepers",
	Long:  `This command deploys keepers (keeper registry + upkeeps).`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		if err := cfg.Validate(); err != nil {
			log.Fatal(err)
		}

		hdlr := handler.NewKeeper(cfg)
		hdlr.DeployKeepers(cmd.Context())
	},
}
