//go:build e2e

package e2e

import (
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"chattery/e2e/client"
	dm_api "chattery/internal/api/dm"
	server_api "chattery/internal/api/server"
	user_api "chattery/internal/api/user"
)

const (
	seededPassword = "password123"
	validPassword  = "Password123"
)

var uniqueUserCounter atomic.Int64

func uniqueCreateUser(t testing.TB, prefix string) *user_api.PostCreateUserRequest {
	t.Helper()

	username := uniqueName(t, prefix, 25)

	return &user_api.PostCreateUserRequest{
		Username: username,
		Login:    username + "@example.com",
		Password: validPassword,
	}
}

func uniqueServerName(t testing.TB, prefix string) string {
	t.Helper()

	return uniqueName(t, prefix, 25)
}

func uniqueTopicName(t testing.TB, prefix string) string {
	t.Helper()

	return uniqueName(t, prefix, 20)
}

func uniqueName(t testing.TB, prefix string, maxLen int) string {
	t.Helper()

	suffix := strconv.FormatInt(time.Now().UnixNano(), 36) +
		strconv.FormatInt(uniqueUserCounter.Add(1), 36)
	name := strings.ToLower(prefix + suffix)
	require.LessOrEqual(t, len(name), maxLen)

	return name
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

func createServer(t testing.TB, session *http.Cookie, name string) int64 {
	t.Helper()

	response := client.MustPostCreateServer(t, &server_api.PostCreateServerRequest{
		Name: name,
	}, session)
	response.RequireStatus(t, http.StatusOK)

	var body server_api.PostCreateServerResponse
	response.RequireJSON(t, &body)
	require.NotZero(t, body.ID)

	return body.ID
}

func createTopic(t testing.TB, session *http.Cookie, serverID int64, name string, topicType string) int64 {
	t.Helper()

	response := client.MustPostCreateTopic(t, &server_api.PostCreateTopicRequest{
		ServerID: serverID,
		Name:     name,
		Type:     topicType,
	}, session)
	response.RequireStatus(t, http.StatusOK)

	var body server_api.PostCreateTopicResponse
	response.RequireJSON(t, &body)
	require.NotZero(t, body.ID)

	return body.ID
}

func createDM(t testing.TB, session *http.Cookie, participantID int64) int64 {
	t.Helper()

	response := client.MustPostCreateDM(t, &dm_api.PostCreateDMRequest{
		ParticipantID: participantID,
	}, session)
	response.RequireStatus(t, http.StatusOK)

	var body dm_api.PostCreateDMResponse
	response.RequireJSON(t, &body)
	require.NotZero(t, body.ID)

	return body.ID
}

func cleanupCreatedServer(t *testing.T, session *http.Cookie, serverID int64) {
	t.Helper()

	t.Cleanup(func() {
		response := client.MustDeleteServer(t, &server_api.DeleteServerRequest{
			ServerID: serverID,
		}, session)
		if response.StatusCode == http.StatusForbidden || response.StatusCode == http.StatusNotFound {
			return
		}

		response.RequireStatus(t, http.StatusOK)
	})
}

func cleanupCreatedTopic(t *testing.T, session *http.Cookie, topicID int64) {
	t.Helper()

	t.Cleanup(func() {
		response := client.MustDeleteTopic(t, &server_api.DeleteTopicRequest{
			TopicID: topicID,
		}, session)
		if response.StatusCode == http.StatusForbidden || response.StatusCode == http.StatusNotFound {
			return
		}

		response.RequireStatus(t, http.StatusOK)
	})
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

	me := getMe(t, session)
	require.Equal(t, username, me.Username)
	require.Equal(t, login, me.Email)
	require.NotZero(t, me.ID)
}

func getMe(t testing.TB, session *http.Cookie) user_api.User {
	t.Helper()
	require.NotNil(t, session)

	response := client.MustGetMe(t, session)
	response.RequireStatus(t, http.StatusOK)

	var body user_api.GetMeResponse
	response.RequireJSON(t, &body)
	require.NotZero(t, body.Me.ID)

	return body.Me
}
