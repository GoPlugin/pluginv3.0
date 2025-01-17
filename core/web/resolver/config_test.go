package resolver

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/goplugin/pluginv3.0/v2/core/services/plugin"
)

var (
	//go:embed testdata/config-empty-effective.toml
	configEmptyEffective string
	//go:embed testdata/config-full.toml
	configFull string
	//go:embed testdata/config-multi-chain.toml
	configMulti string
	//go:embed testdata/config-multi-chain-effective.toml
	configMultiEffective string
)

func TestResolver_ConfigV2(t *testing.T) {
	t.Parallel()

	query := `
		query FetchConfigV2 {
			configv2 {
				user
				effective
			}
	  	}`

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "configv2"),
		{
			name:          "empty",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				opts := plugin.GeneralConfigOpts{}
				cfg, err := opts.New()
				require.NoError(t, err)
				f.App.On("GetConfig").Return(cfg)
			},
			query:  query,
			result: fmt.Sprintf(`{"configv2":{"user":"","effective":%s}}`, mustJSONMarshal(t, configEmptyEffective)),
		},
		{
			name:          "full",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				opts := plugin.GeneralConfigOpts{
					ConfigStrings:  []string{configFull},
					SecretsStrings: []string{},
				}
				cfg, err := opts.New()
				require.NoError(t, err)
				f.App.On("GetConfig").Return(cfg)
			},
			query:  query,
			result: fmt.Sprintf(`{"configv2":{"user":%s,"effective":%s}}`, mustJSONMarshal(t, configFull), mustJSONMarshal(t, configFull)),
		},
		{
			name:          "partial",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				opts := plugin.GeneralConfigOpts{
					ConfigStrings:  []string{configMulti},
					SecretsStrings: []string{},
				}
				cfg, err := opts.New()
				require.NoError(t, err)
				f.App.On("GetConfig").Return(cfg)
			},
			query:  query,
			result: fmt.Sprintf(`{"configv2":{"user":%s,"effective":%s}}`, mustJSONMarshal(t, configMulti), mustJSONMarshal(t, configMultiEffective)),
		},
	}

	RunGQLTests(t, testCases)
}

func mustJSONMarshal(t *testing.T, s string) string {
	b, err := json.Marshal(s)
	require.NoError(t, err)
	return string(b)
}
