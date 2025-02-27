package managed

import (
	"context"

	"github.com/goplugin/plugin-libocr/commontypes"
	"github.com/goplugin/plugin-libocr/internal/loghelper"
	"github.com/goplugin/plugin-libocr/offchainreporting/types"
)

func loadConfigFromDatabase(ctx context.Context, database types.Database, logger loghelper.LoggerWithContext) *types.ContractConfig {
	cc, err := database.ReadConfig(ctx)
	if err != nil {
		logger.ErrorIfNotCanceled("loadConfigFromDatabase: Error during Database.ReadConfig", ctx, commontypes.LogFields{
			"error": err,
		})
		return nil
	}

	if cc == nil {
		logger.Info("loadConfigFromDatabase: Database.ReadConfig returned nil, no configuration to restore", nil)
		return nil
	}

	return cc
}
