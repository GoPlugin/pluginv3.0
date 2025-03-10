package factory

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	cciptypes "github.com/goplugin/plugin-common/pkg/types/ccip"

	"github.com/goplugin/plugin-common/pkg/logger"

	"github.com/goplugin/pluginv3.0/v2/core/chains/evm/logpoller"
	mocks2 "github.com/goplugin/pluginv3.0/v2/core/chains/evm/logpoller/mocks"
	"github.com/goplugin/pluginv3.0/v2/core/chains/evm/utils"
	ccipconfig "github.com/goplugin/pluginv3.0/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/goplugin/pluginv3.0/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/goplugin/pluginv3.0/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
)

func TestCommitStore(t *testing.T) {
	for _, versionStr := range []string{ccipdata.V1_2_0} {
		lggr := logger.Test(t)
		addr := cciptypes.Address(utils.RandomAddress().String())
		lp := mocks2.NewLogPoller(t)

		lp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil)
		versionFinder := newMockVersionFinder(ccipconfig.CommitStore, *semver.MustParse(versionStr), nil)
		_, err := NewCommitStoreReader(lggr, versionFinder, addr, nil, lp)
		assert.NoError(t, err)

		expFilterName := logpoller.FilterName(v1_2_0.ExecReportAccepts, addr)
		lp.On("UnregisterFilter", mock.Anything, expFilterName).Return(nil)
		err = CloseCommitStoreReader(lggr, versionFinder, addr, nil, lp)
		assert.NoError(t, err)
	}
}
