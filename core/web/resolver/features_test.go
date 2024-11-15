package resolver

import (
	"context"
	"testing"

	"github.com/goplugin/pluginv3.0/v2/core/internal/testutils/configtest"
	"github.com/goplugin/pluginv3.0/v2/core/services/plugin"
)

func Test_ToFeatures(t *testing.T) {
	query := `
	{
		features {
			... on Features {
				csa
				feedsManager
				multiFeedsManagers
			}	
		}
	}`

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "features"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				f.App.On("GetConfig").Return(configtest.NewGeneralConfig(t, func(c *plugin.Config, s *plugin.Secrets) {
					t, f := true, false
					c.Feature.UICSAKeys = &f
					c.Feature.FeedsManager = &t
					c.Feature.MultiFeedsManagers = &f
				}))
			},
			query: query,
			result: `
			{
				"features": {
					"csa": false,
					"feedsManager": true,
					"multiFeedsManagers": false
				}
			}`,
		},
	}

	RunGQLTests(t, testCases)
}
