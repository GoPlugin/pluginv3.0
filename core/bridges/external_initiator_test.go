package bridges_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/goplugin/pluginv3.0/v2/core/auth"
	"github.com/goplugin/pluginv3.0/v2/core/bridges"
	"github.com/goplugin/pluginv3.0/v2/core/internal/cltest"

	"github.com/stretchr/testify/assert"
)

func TestNewExternalInitiator(t *testing.T) {
	eia := auth.NewToken()
	assert.Len(t, eia.AccessKey, 32)
	assert.Len(t, eia.Secret, 64)

	url := cltest.WebURL(t, "http://localhost:8888")
	name := uuid.New().String()
	eir := &bridges.ExternalInitiatorRequest{
		Name: name,
		URL:  &url,
	}
	ei, err := bridges.NewExternalInitiator(eia, eir)
	assert.NoError(t, err)
	assert.NotEqual(t, ei.HashedSecret, eia.Secret)
	assert.Equal(t, ei.AccessKey, eia.AccessKey)
}

func TestAuthenticateExternalInitiator(t *testing.T) {
	eia := auth.NewToken()
	ok, err := bridges.AuthenticateExternalInitiator(eia, &bridges.ExternalInitiator{
		Salt:         "salt",
		HashedSecret: "secret",
	})
	require.NoError(t, err)
	require.False(t, ok)

	hs, err := auth.HashedSecret(eia, "salt")
	require.NoError(t, err)
	ok, err = bridges.AuthenticateExternalInitiator(eia, &bridges.ExternalInitiator{
		Salt:         "salt",
		HashedSecret: hs,
	})
	require.NoError(t, err)
	require.True(t, ok)
}
