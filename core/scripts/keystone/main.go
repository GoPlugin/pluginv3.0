package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/goplugin/pluginv3.0/core/scripts/keystone/src"
)

type command interface {
	Run([]string)
	Name() string
}

func main() {
	commands := []command{
		src.NewDeployContractsCommand(),
		src.NewDeployJobSpecsCommand(),
		src.NewGenerateCribClusterOverridesCommand(),
		src.NewDeleteJobsCommand(),
		src.NewDeployAndInitializeCapabilitiesRegistryCommand(),
		src.NewDeployWorkflowsCommand(),
		src.NewDeleteWorkflowsCommand(),
	}

	commandsList := func(commands []command) string {
		var scs []string
		for _, command := range commands {
			scs = append(scs, command.Name())
		}
		return strings.Join(scs, ", ")
	}(commands)

	if len(os.Args) >= 2 {
		requestedCommand := os.Args[1]

		for _, command := range commands {
			if command.Name() == requestedCommand {
				command.Run(os.Args[2:])
				return
			}
		}
		fmt.Println("Unknown command:", requestedCommand)
	} else {
		fmt.Println("No command specified")
	}

	fmt.Println("Supported commands:", commandsList)
	os.Exit(1)
}
