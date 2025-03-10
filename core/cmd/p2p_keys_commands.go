package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.uber.org/multierr"

	cutils "github.com/goplugin/plugin-common/pkg/utils"
	"github.com/goplugin/pluginv3.0/v2/core/utils"
	"github.com/goplugin/pluginv3.0/v2/core/web/presenters"
)

func initP2PKeysSubCmd(s *Shell) cli.Command {
	return cli.Command{
		Name:  "p2p",
		Usage: "Remote commands for administering the node's p2p keys",
		Subcommands: cli.Commands{
			{
				Name:   "create",
				Usage:  format(`Create a p2p key, encrypted with password from the password file, and store it in the database.`),
				Action: s.CreateP2PKey,
			},
			{
				Name:  "delete",
				Usage: format(`Delete the encrypted P2P key by id`),
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "yes, y",
						Usage: "skip the confirmation prompt",
					},
					cli.BoolFlag{
						Name:  "hard",
						Usage: "hard-delete the key instead of archiving (irreversible!)",
					},
				},
				Action: s.DeleteP2PKey,
			},
			{
				Name:   "list",
				Usage:  format(`List available P2P keys`),
				Action: s.ListP2PKeys,
			},
			{
				Name:  "import",
				Usage: format(`Imports a P2P key from a JSON file`),
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "old-password, oldpassword, p",
						Usage: "`FILE` containing the password used to encrypt the key in the JSON file",
					},
				},
				Action: s.ImportP2PKey,
			},
			{
				Name:  "export",
				Usage: format(`Exports a P2P key to a JSON file`),
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "new-password, newpassword, p",
						Usage: "`FILE` containing the password to encrypt the key (required)",
					},
					cli.StringFlag{
						Name:  "output, o",
						Usage: "`FILE` where the JSON file will be saved (required)",
					},
				},
				Action: s.ExportP2PKey,
			},
		},
	}
}

type P2PKeyPresenter struct {
	JAID
	presenters.P2PKeyResource
}

// RenderTable implements TableRenderer
func (p *P2PKeyPresenter) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Peer ID", "Public key"}
	rows := [][]string{p.ToRow()}

	if _, err := rt.Write([]byte("🔑 P2P Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return cutils.JustError(rt.Write([]byte("\n")))
}

func (p *P2PKeyPresenter) ToRow() []string {
	row := []string{
		p.ID,
		p.PeerID,
		p.PubKey,
	}

	return row
}

type P2PKeyPresenters []P2PKeyPresenter

// RenderTable implements TableRenderer
func (ps P2PKeyPresenters) RenderTable(rt RendererTable) error {
	headers := []string{"ID", "Peer ID", "Public key"}
	rows := [][]string{}

	for _, p := range ps {
		rows = append(rows, p.ToRow())
	}

	if _, err := rt.Write([]byte("🔑 P2P Keys\n")); err != nil {
		return err
	}
	renderList(headers, rows, rt.Writer)

	return cutils.JustError(rt.Write([]byte("\n")))
}

// ListP2PKeys retrieves a list of all P2P keys
func (s *Shell) ListP2PKeys(_ *cli.Context) (err error) {
	resp, err := s.HTTP.Get(s.ctx(), "/v2/keys/p2p", nil)
	if err != nil {
		return s.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return s.renderAPIResponse(resp, &P2PKeyPresenters{})
}

// CreateP2PKey creates a new P2P key
func (s *Shell) CreateP2PKey(_ *cli.Context) (err error) {
	resp, err := s.HTTP.Post(s.ctx(), "/v2/keys/p2p", nil)
	if err != nil {
		return s.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return s.renderAPIResponse(resp, &P2PKeyPresenter{}, "Created P2P keypair")
}

// DeleteP2PKey deletes a P2P key,
// key ID must be passed
func (s *Shell) DeleteP2PKey(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return s.errorOut(errors.New("Must pass the key ID to be deleted"))
	}
	id := c.Args().Get(0)

	if !confirmAction(c) {
		return nil
	}

	var queryStr string
	if c.Bool("hard") {
		queryStr = "?hard=true"
	}

	resp, err := s.HTTP.Delete(s.ctx(), fmt.Sprintf("/v2/keys/p2p/%s%s", id, queryStr))
	if err != nil {
		return s.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return s.renderAPIResponse(resp, &P2PKeyPresenter{}, "P2P key deleted")
}

// ImportP2PKey imports and stores a P2P key,
// path to key must be passed
func (s *Shell) ImportP2PKey(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return s.errorOut(errors.New("Must pass the filepath of the key to be imported"))
	}

	oldPasswordFile := c.String("old-password")
	if len(oldPasswordFile) == 0 {
		return s.errorOut(errors.New("Must specify --old-password/-p flag"))
	}
	oldPassword, err := os.ReadFile(oldPasswordFile)
	if err != nil {
		return s.errorOut(errors.Wrap(err, "Could not read password file"))
	}

	filepath := c.Args().Get(0)
	keyJSON, err := os.ReadFile(filepath)
	if err != nil {
		return s.errorOut(err)
	}

	normalizedPassword := normalizePassword(string(oldPassword))
	resp, err := s.HTTP.Post(s.ctx(), "/v2/keys/p2p/import?oldpassword="+normalizedPassword, bytes.NewReader(keyJSON))
	if err != nil {
		return s.errorOut(err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	return s.renderAPIResponse(resp, &P2PKeyPresenter{}, "🔑 Imported P2P key")
}

// ExportP2PKey exports a P2P key,
// key ID must be passed
func (s *Shell) ExportP2PKey(c *cli.Context) (err error) {
	if !c.Args().Present() {
		return s.errorOut(errors.New("Must pass the ID of the key to export"))
	}

	newPasswordFile := c.String("new-password")
	if len(newPasswordFile) == 0 {
		return s.errorOut(errors.New("Must specify --new-password/-p flag"))
	}
	newPassword, err := os.ReadFile(newPasswordFile)
	if err != nil {
		return s.errorOut(errors.Wrap(err, "Could not read password file"))
	}

	filepath := c.String("output")
	if len(filepath) == 0 {
		return s.errorOut(errors.New("Must specify --output/-o flag"))
	}

	ID := c.Args().Get(0)

	normalizedPassword := normalizePassword(string(newPassword))
	resp, err := s.HTTP.Post(s.ctx(), "/v2/keys/p2p/export/"+ID+"?newpassword="+normalizedPassword, nil)
	if err != nil {
		return s.errorOut(errors.Wrap(err, "Could not make HTTP request"))
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			err = multierr.Append(err, cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return s.errorOut(fmt.Errorf("error exporting: %w", httpError(resp)))
	}

	keyJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		return s.errorOut(errors.Wrap(err, "Could not read response body"))
	}

	err = utils.WriteFileWithMaxPerms(filepath, keyJSON, 0o600)
	if err != nil {
		return s.errorOut(errors.Wrapf(err, "Could not write %v", filepath))
	}

	_, err = os.Stderr.WriteString(fmt.Sprintf("🔑 Exported P2P key %s to %s\n", ID, filepath))
	if err != nil {
		return s.errorOut(err)
	}

	return nil
}
