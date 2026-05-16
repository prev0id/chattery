//go:build e2e

package e2e

import (
	"net/http"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"chattery/e2e/client"
	user_api "chattery/internal/api/user"
)

const (
	seededPassword = "password123"
	validPassword  = "Password123"
)

var uniqueUserCounter atomic.Int64

func uniqueCreateUser(t testing.TB, prefix string) *user_api.PostCreateUserRequest {
	t.Helper()

	suffix := strconv.FormatInt(time.Now().UnixNano(), 36) +
		strconv.FormatInt(uniqueUserCounter.Add(1), 36)
	username := prefix + suffix
	require.LessOrEqual(t, len(username), 25)

	return &user_api.PostCreateUserRequest{
		Username: username,
		Login:    username + "@example.com",
		Password: validPassword,
	}
}

func createUser(t testing.TB, prefix string) (*user_api.PostCreateUserRequest, *http.Cookie) {
	t.Helper()

	user := uniqueCreateUser(t, prefix)
	response := client.MustPostCreateUser(t, user)
	response.RequireStatus(t, http.StatusOK)
	require.Empty(t, response.Body)

	return user, response.RequireSessionCookie(t)
}

func loginUser(t testing.TB, login, password string) *http.Cookie {
	t.Helper()

	response := client.MustPostLogin(t, &user_api.PostLoginRequest{
		Login:    login,
		Password: password,
	})
	response.RequireStatus(t, http.StatusOK)
	require.Empty(t, response.Body)

	return response.RequireSessionCookie(t)
}

func cleanupCreatedUser(t *testing.T, session *http.Cookie) {
	t.Helper()

	t.Cleanup(func() {
		response := client.MustDeleteMe(t, session)
		if response.StatusCode == http.StatusUnauthorized {
			return
		}

		response.RequireStatus(t, http.StatusOK)
		response.RequireClearedSessionCookie(t)
	})
}

func cleanupSession(t *testing.T, session *http.Cookie) {
	t.Helper()

	t.Cleanup(func() {
		response := client.MustPostLogout(t, session)
		if response.StatusCode == http.StatusUnauthorized {
			return
		}

		response.RequireStatus(t, http.StatusSeeOther)
		response.RequireClearedSessionCookie(t)
	})
}

func requireMe(t testing.TB, session *http.Cookie, username, login string) {
	t.Helper()
	require.NotNil(t, session)

	response := client.MustGetMe(t, session)
	response.RequireStatus(t, http.StatusOK)

	var body user_api.GetMeResponse
	response.RequireJSON(t, &body)
	require.Equal(t, username, body.Me.Username)
	require.Equal(t, login, body.Me.Email)
	require.NotZero(t, body.Me.ID)
}
