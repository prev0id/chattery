//go:build e2e

package e2e

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"chattery/e2e/client"
	user_api "chattery/internal/api/user"
)

func Test_CreateUser(t *testing.T) {
	t.Parallel()

	user := uniqueCreateUser(t, "create")
	var session *http.Cookie

	t.Run("create_user", func(t *testing.T) {
		response := client.MustPostCreateUser(t, user)
		response.RequireStatus(t, http.StatusOK)
		session = response.RequireSessionCookie(t)
		require.Empty(t, response.Body)
	})

	require.NotNil(t, session)
	cleanupCreatedUser(t, session)

	t.Run("check_session", func(t *testing.T) {
		requireMe(t, session, user.Username, user.Login)
	})

	t.Run("reject_invalid_login", func(t *testing.T) {
		request := uniqueCreateUser(t, "badlogin")
		request.Login = "not-email"

		response := client.MustPostCreateUser(t, request)
		response.RequireStatus(t, http.StatusBadRequest)
		response.RequireNoSessionCookie(t)
		response.RequireErrorContains(t, "login must be a valid email address")
	})

	t.Run("reject_weak_password", func(t *testing.T) {
		request := uniqueCreateUser(t, "badpass")
		request.Password = seededPassword

		response := client.MustPostCreateUser(t, request)
		response.RequireStatus(t, http.StatusBadRequest)
		response.RequireNoSessionCookie(t)
		response.RequireErrorContains(t, "password must contain at least one uppercase letter")
	})

	t.Run("reject_duplicate_login", func(t *testing.T) {
		request := uniqueCreateUser(t, "duplogin")
		request.Login = user.Login

		response := client.MustPostCreateUser(t, request)
		response.RequireStatus(t, http.StatusConflict)
		response.RequireNoSessionCookie(t)
		response.RequireErrorContains(t, "already exists")
		response.RequireErrorContains(t, user.Login)
	})

	t.Run("reject_duplicate_username", func(t *testing.T) {
		request := uniqueCreateUser(t, "dupname")
		request.Username = user.Username

		response := client.MustPostCreateUser(t, request)
		response.RequireStatus(t, http.StatusConflict)
		response.RequireNoSessionCookie(t)
		response.RequireErrorContains(t, "already exists")
		response.RequireErrorContains(t, user.Username)
	})
}

func Test_Login(t *testing.T) {
	t.Parallel()

	var session *http.Cookie

	t.Run("login_existing_user", func(t *testing.T) {
		response := client.MustPostLogin(t, &user_api.PostLoginRequest{
			Login:    "alex@example.com",
			Password: seededPassword,
		})
		response.RequireStatus(t, http.StatusOK)
		session = response.RequireSessionCookie(t)
		require.Empty(t, response.Body)
	})

	require.NotNil(t, session)
	cleanupSession(t, session)

	t.Run("check_session", func(t *testing.T) {
		requireMe(t, session, "alex", "alex@example.com")
	})

	t.Run("reject_wrong_password", func(t *testing.T) {
		response := client.MustPostLogin(t, &user_api.PostLoginRequest{
			Login:    "alex@example.com",
			Password: "wrongPassword123",
		})
		response.RequireStatus(t, http.StatusForbidden)
		response.RequireNoSessionCookie(t)
		response.RequireErrorContains(t, "invalid password")
	})

	t.Run("reject_unknown_user", func(t *testing.T) {
		response := client.MustPostLogin(t, &user_api.PostLoginRequest{
			Login:    uniqueCreateUser(t, "unknown").Login,
			Password: seededPassword,
		})
		response.RequireStatus(t, http.StatusForbidden)
		response.RequireNoSessionCookie(t)
		response.RequireErrorContains(t, "not found")
	})
}

func Test_Logout(t *testing.T) {
	t.Parallel()

	var (
		session       *http.Cookie
		closedSession *http.Cookie
	)

	t.Run("login_user", func(t *testing.T) {
		response := client.MustPostLogin(t, &user_api.PostLoginRequest{
			Login:    "bob@example.com",
			Password: seededPassword,
		})
		response.RequireStatus(t, http.StatusOK)
		session = response.RequireSessionCookie(t)
	})

	require.NotNil(t, session)
	cleanupSession(t, session)

	t.Run("logout_user", func(t *testing.T) {
		require.NotNil(t, session)

		response := client.MustPostLogout(t, session)
		response.RequireStatus(t, http.StatusSeeOther)
		response.RequireClearedSessionCookie(t)
		require.Equal(t, "/login", response.Header.Get("Location"))

		closedSession = session
	})

	t.Run("check_session_closed", func(t *testing.T) {
		require.NotNil(t, closedSession)

		response := client.MustGetMe(t, closedSession)
		response.RequireStatus(t, http.StatusUnauthorized)
		response.RequireClearedSessionCookie(t)
		response.RequireErrorContains(t, "login required")
	})

	t.Run("reject_logout_without_session", func(t *testing.T) {
		response := client.MustPostLogout(t)
		response.RequireStatus(t, http.StatusUnauthorized)
		response.RequireNoSessionCookie(t)
		response.RequireErrorContains(t, "login required")
	})
}
