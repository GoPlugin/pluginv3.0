package evm

import (
	"github.com/goplugin/plugin-libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/goplugin/plugin-libocr/offchainreporting2plus/types"

	"github.com/goplugin/plugin-automation/pkg/v3/plugin"
)

type AutomationServices interface {
	Keyring() ocr3types.OnchainKeyring[plugin.AutomationReportInfo]
}

func New(keyring ocrtypes.OnchainKeyring) (AutomationServices, error) {
	services := new(automationServices)

	services.keyring = NewOnchainKeyringV3Wrapper(keyring)

	return services, nil
}

type automationServices struct {
	keyring *onchainKeyringV3Wrapper
}

var _ AutomationServices = &automationServices{}

func (f *automationServices) Keyring() ocr3types.OnchainKeyring[plugin.AutomationReportInfo] {
	return f.keyring
}
