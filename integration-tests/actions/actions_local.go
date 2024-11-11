// Package actions enables common plugin interactions
package actions

import (
	"github.com/goplugin/pluginv3.0/integration-tests/docker/test_env"
)

// UpgradePluginNodeVersions upgrades all Plugin nodes to a new version, and then runs the test environment
// to apply the upgrades
func UpgradePluginNodeVersionsLocal(
	newImage, newVersion string,
	nodes ...*test_env.ClNode,
) error {
	for _, node := range nodes {
		if err := node.UpgradeVersion(newImage, newVersion); err != nil {
			return err
		}
	}
	return nil
}
