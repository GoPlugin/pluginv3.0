package ldapauth

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"

	"github.com/goplugin/pluginv3.0/v2/core/config"
)

type ldapClient struct {
	config config.LDAP
}

// Wrapper for creating a handle to a *ldap.Conn/LDAPConn interface
type LDAPClient interface {
	CreateEphemeralConnection() (LDAPConn, error)
}

// Wrapper for ldap connection and mock testing, implemented by *ldap.Conn
type LDAPConn interface {
	Search(searchRequest *ldap.SearchRequest) (*ldap.SearchResult, error)
	Bind(username string, password string) error
	Close() (err error)
}

func newLDAPClient(config config.LDAP) LDAPClient {
	return &ldapClient{config}
}

// CreateEphemeralConnection returns a valid, active LDAP connection for upstream Search and Bind queries
func (l *ldapClient) CreateEphemeralConnection() (LDAPConn, error) {
	conn, err := ldap.DialURL(l.config.ServerAddress())
	if err != nil {
		return nil, fmt.Errorf("failed to Dial LDAP Server: %w", err)
	}
	// Root level root user auth with credentials provided from config
	bindStr := l.config.BaseUserAttr() + "=" + l.config.ReadOnlyUserLogin() + "," + l.config.BaseDN()
	if err := conn.Bind(bindStr, l.config.ReadOnlyUserPass()); err != nil {
		return nil, fmt.Errorf("unable to login as initial root LDAP user: %w", err)
	}
	return conn, nil
}
