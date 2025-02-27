package mailboxtest

import (
	"testing"

	"github.com/goplugin/plugin-common/pkg/logger"
	"github.com/goplugin/plugin-common/pkg/utils/mailbox"
)

func NewMonitor(t testing.TB) *mailbox.Monitor {
	return mailbox.NewMonitor(t.Name(), logger.Named(logger.Test(t), "Mailbox"))
}
