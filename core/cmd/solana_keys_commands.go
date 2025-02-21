package cmd

import (
	"github.com/goplugin/plugin-common/pkg/utils"
	"github.com/goplugin/pluginv3.0/v2/core/services/keystore/keys/solkey"
	"github.com/goplugin/pluginv3.0/v2/core/web/presenters"
)

type SolanaKeyPresenter struct {
	JAID
	presenters.SolanaKeyResource
}

// RenderTable implements TableRenderer
func (p SolanaKeyPresenter) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Public key"}
	rows := [][]string{p.ToRow()}

	if _, err := rt.Write([]byte("🔑 Solana Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func (p *SolanaKeyPresenter) ToRow() []string {
	row := []string{
		p.ID,
		p.PubKey,
	}

	return row
}

type SolanaKeyPresenters []SolanaKeyPresenter

// RenderTable implements TableRenderer
func (ps SolanaKeyPresenters) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Public key"}
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	if _, err := rt.Write([]byte("🔑 Solana Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return utils.JustError(rt.Write([]byte("\n")))
}

func NewSolanaKeysClient(s *Shell) KeysClient {
	return newKeysClient[solkey.Key, SolanaKeyPresenter, SolanaKeyPresenters]("Solana", s)
}
