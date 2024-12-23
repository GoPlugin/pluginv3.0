package aptoskey

import (
	"testing"

	"github.com/goplugin/pluginv3.0/v2/core/services/keystore/keys"
)

func TestAptosKeys_ExportImport(t *testing.T) {
	keys.RunKeyExportImportTestcase(t, createKey, decryptKey)
}

func createKey() (keys.KeyType, error) {
	return New()
}

func decryptKey(keyJSON []byte, password string) (keys.KeyType, error) {
	return FromEncryptedJSON(keyJSON, password)
}
