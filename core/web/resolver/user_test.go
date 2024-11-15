package resolver

import (
	"context"
	"testing"

	gqlerrors "github.com/graph-gophers/graphql-go/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/goplugin/pluginv3.0/v2/core/utils"
	"github.com/goplugin/pluginv3.0/v2/core/web/auth"
)

func TestResolver_UpdateUserPassword(t *testing.T) {
	t.Parallel()

	mutation := `
		mutation UpdateUserPassword($input: UpdatePasswordInput!) {
			updateUserPassword(input: $input) {
				... on UpdatePasswordSuccess {
					user {
						email
					}
				}
				... on InputErrors {
					errors {
						path
						message
						code
					}
				}
			}
		}`
	oldPassword := "old"
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"newPassword": "new",
			"oldPassword": oldPassword,
		},
	}

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: mutation, variables: variables}, "updateUserPassword"),
		{
			name:          "success",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				session, ok := auth.GetGQLAuthenticatedSession(ctx)
				require.True(t, ok)
				require.NotNil(t, session)

				pwd, err := utils.HashPassword(oldPassword)
				require.NoError(t, err)

				session.User.HashedPassword = pwd

				f.Mocks.authProvider.On("FindUser", mock.Anything, session.User.Email).Return(*session.User, nil)
				f.Mocks.authProvider.On("SetPassword", mock.Anything, session.User, "new").Return(nil)
				f.Mocks.authProvider.On("ClearNonCurrentSessions", mock.Anything, session.SessionID).Return(nil)
				f.App.On("AuthenticationProvider").Return(f.Mocks.authProvider)
			},
			query:     mutation,
			variables: variables,
			result: `
				{
					"updateUserPassword": {
						"user": {
							"email": "gqltester@chain.link"
						}
					}
				}`,
		},
		{
			name:          "update password match error",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				session, ok := auth.GetGQLAuthenticatedSession(ctx)
				require.True(t, ok)
				require.NotNil(t, session)

				session.User.HashedPassword = "random-string"

				f.Mocks.authProvider.On("FindUser", mock.Anything, session.User.Email).Return(*session.User, nil)
				f.App.On("AuthenticationProvider").Return(f.Mocks.authProvider)
			},
			query:     mutation,
			variables: variables,
			result: `
				{
					"updateUserPassword": {
						"errors": [{
							"path": "oldPassword",
							"message": "old password does not match",
							"code": "INVALID_INPUT"
						}]
					}
				}`,
		},
		{
			name:          "failed to clear session error",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				session, ok := auth.GetGQLAuthenticatedSession(ctx)
				require.True(t, ok)
				require.NotNil(t, session)

				pwd, err := utils.HashPassword(oldPassword)
				require.NoError(t, err)

				session.User.HashedPassword = pwd

				f.Mocks.authProvider.On("FindUser", mock.Anything, session.User.Email).Return(*session.User, nil)
				f.Mocks.authProvider.On("ClearNonCurrentSessions", mock.Anything, session.SessionID).Return(
					clearSessionsError{},
				)
				f.App.On("AuthenticationProvider").Return(f.Mocks.authProvider)
			},
			query:     mutation,
			variables: variables,
			result:    `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: clearSessionsError{},
					Path:          []interface{}{"updateUserPassword"},
					Message:       "failed to clear non current user sessions",
				},
			},
		},
		{
			name:          "failed to update current user password error",
			authenticated: true,
			before: func(ctx context.Context, f *gqlTestFramework) {
				session, ok := auth.GetGQLAuthenticatedSession(ctx)
				require.True(t, ok)
				require.NotNil(t, session)

				pwd, err := utils.HashPassword(oldPassword)
				require.NoError(t, err)

				session.User.HashedPassword = pwd

				f.Mocks.authProvider.On("FindUser", mock.Anything, session.User.Email).Return(*session.User, nil)
				f.Mocks.authProvider.On("ClearNonCurrentSessions", mock.Anything, session.SessionID).Return(nil)
				f.Mocks.authProvider.On("SetPassword", mock.Anything, session.User, "new").Return(failedPasswordUpdateError{})
				f.App.On("AuthenticationProvider").Return(f.Mocks.authProvider)
			},
			query:     mutation,
			variables: variables,
			result:    `null`,
			errors: []*gqlerrors.QueryError{
				{
					Extensions:    nil,
					ResolverError: failedPasswordUpdateError{},
					Path:          []interface{}{"updateUserPassword"},
					Message:       "failed to update current user password",
				},
			},
		},
	}

	RunGQLTests(t, testCases)
}
