// Package client enables interaction with APIs of test components like the mockserver and Plugin nodes
package client

import (
	"os"
	"regexp"

	"github.com/rs/zerolog/log"

	"github.com/goplugin/plugin-testing-framework/k8s/environment"
)

type PluginK8sClient struct {
	ChartName string
	PodName   string
	*PluginClient
}

// NewPlugin creates a new Plugin model using a provided config
func NewPluginK8sClient(c *PluginConfig, podName, chartName string) (*PluginK8sClient, error) {
	rc, err := initRestyClient(c.URL, c.Email, c.Password, c.HTTPTimeout)
	if err != nil {
		return nil, err
	}
	_, isSet := os.LookupEnv("CL_CLIENT_DEBUG")
	if isSet {
		rc.SetDebug(true)
	}
	return &PluginK8sClient{
		PluginClient: &PluginClient{
			APIClient: rc,
			pageSize:  25,
			Config:    c,
		},
		ChartName: chartName,
		PodName:   podName,
	}, nil
}

// UpgradeVersion upgrades the plugin node to the new version
// Note: You need to call Run() on the test environment for changes to take effect
// Note: This function is not thread safe, call from a single thread
func (c *PluginK8sClient) UpgradeVersion(testEnvironment *environment.Environment, newImage, newVersion string) error {
	log.Info().
		Str("Chart Name", c.ChartName).
		Str("New Image", newImage).
		Str("New Version", newVersion).
		Msg("Upgrading Plugin Node")
	upgradeVals := map[string]any{
		"plugin": map[string]any{
			"image": map[string]any{
				"image":   newImage,
				"version": newVersion,
			},
		},
	}
	_, err := testEnvironment.UpdateHelm(c.ChartName, upgradeVals)
	return err
}

// Name Plugin instance chart/service name
func (c *PluginK8sClient) Name() string {
	return c.ChartName
}

func parseHostname(s string) string {
	r := regexp.MustCompile(`://(?P<Host>.*):`)
	return r.FindStringSubmatch(s)[1]
}

// ConnectPluginNodes creates new Plugin clients
func ConnectPluginNodes(e *environment.Environment) ([]*PluginK8sClient, error) {
	var clients []*PluginK8sClient
	for _, nodeDetails := range e.PluginNodeDetails {
		c, err := NewPluginK8sClient(&PluginConfig{
			URL:        nodeDetails.LocalIP,
			Email:      "notreal@fakeemail.ch",
			Password:   "fj293fbBnlQ!f9vNs",
			InternalIP: parseHostname(nodeDetails.InternalIP),
		}, nodeDetails.PodName, nodeDetails.ChartName)
		if err != nil {
			return nil, err
		}
		log.Debug().
			Str("URL", c.Config.URL).
			Str("Internal IP", c.Config.InternalIP).
			Str("Chart Name", nodeDetails.ChartName).
			Str("Pod Name", nodeDetails.PodName).
			Msg("Connected to Plugin node")
		clients = append(clients, c)
	}
	return clients, nil
}

// ReconnectPluginNodes reconnects to Plugin nodes after they have been modified, say through a Helm upgrade
// Note: Experimental as of now, will likely not work predictably.
func ReconnectPluginNodes(testEnvironment *environment.Environment, nodes []*PluginK8sClient) (err error) {
	for _, node := range nodes {
		for _, details := range testEnvironment.PluginNodeDetails {
			if details.ChartName == node.ChartName { // Make the link from client to pod consistent
				node, err = NewPluginK8sClient(&PluginConfig{
					URL:        details.LocalIP,
					Email:      "notreal@fakeemail.ch",
					Password:   "fj293fbBnlQ!f9vNs",
					InternalIP: parseHostname(details.InternalIP),
				}, details.PodName, details.ChartName)
				if err != nil {
					return err
				}
				log.Debug().
					Str("URL", node.Config.URL).
					Str("Internal IP", node.Config.InternalIP).
					Str("Chart Name", node.ChartName).
					Str("Pod Name", node.PodName).
					Msg("Reconnected to Plugin node")
			}
		}
	}
	return nil
}

// ConnectPluginNodeURLs creates new Plugin clients based on just URLs, should only be used inside K8s tests
func ConnectPluginNodeURLs(urls []string) ([]*PluginK8sClient, error) {
	var clients []*PluginK8sClient
	for _, url := range urls {
		c, err := ConnectPluginNodeURL(url)
		if err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}
	return clients, nil
}

// ConnectPluginNodeURL creates a new Plugin client based on just a URL, should only be used inside K8s tests
func ConnectPluginNodeURL(url string) (*PluginK8sClient, error) {
	return NewPluginK8sClient(&PluginConfig{
		URL:        url,
		Email:      "notreal@fakeemail.ch",
		Password:   "fj293fbBnlQ!f9vNs",
		InternalIP: parseHostname(url),
	},
		parseHostname(url),   // a decent guess
		"connectedNodeByURL", // an intentionally bad decision
	)
}
