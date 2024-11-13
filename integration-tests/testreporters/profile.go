package testreporters

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
	"golang.org/x/sync/errgroup"

	"github.com/goplugin/plugin-testing-framework/lib/testreporters"
	"github.com/goplugin/pluginv3.0/integration-tests/client"
)

type PluginProfileTestReporter struct {
	Results   []*client.PluginProfileResults
	namespace string
}

// SetNamespace sets the namespace of the report for clean reports
func (c *PluginProfileTestReporter) SetNamespace(namespace string) {
	c.namespace = namespace
}

// WriteReport create the profile files
func (c *PluginProfileTestReporter) WriteReport(folderLocation string) error {
	profFiles := new(errgroup.Group)
	for _, res := range c.Results {
		result := res
		profFiles.Go(func() error {
			filePath := filepath.Join(folderLocation, fmt.Sprintf("plugin-node-%d-profiles", result.NodeIndex))
			if err := testreporters.MkdirIfNotExists(filePath); err != nil {
				return err
			}
			for _, rep := range result.Reports {
				report := rep
				reportFile, err := os.Create(filepath.Join(filePath, report.Type))
				if err != nil {
					return err
				}
				if _, err = reportFile.Write(report.Data); err != nil {
					return err
				}
				if err = reportFile.Close(); err != nil {
					return err
				}
			}
			return nil
		})
	}
	return profFiles.Wait()
}

// SendNotification hasn't been implemented for this test
func (c *PluginProfileTestReporter) SendSlackNotification(_ *testing.T, _ *slack.Client, _ testreporters.GrafanaURLProvider) error {
	log.Warn().Msg("No Slack notification integration for Plugin profile tests")
	return nil
}
