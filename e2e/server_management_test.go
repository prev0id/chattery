//go:build e2e

package e2e

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"chattery/e2e/client"
	server_api "chattery/internal/api/server"
)

func Test_ServerManagement(t *testing.T) {
	t.Parallel()

	var (
		ownerSession *http.Cookie
		serverID     int64
		textTopicID  int64
		voiceTopicID int64
	)

	serverName := uniqueServerName(t, "server")
	updatedServerName := uniqueServerName(t, "srvupd")
	textTopicName := uniqueTopicName(t, "text")
	updatedTextTopicName := uniqueTopicName(t, "upd")
	voiceTopicName := uniqueTopicName(t, "voice")

	t.Run("create_owner", func(t *testing.T) {
		_, ownerSession = createUser(t, "srvowner")
	})

	require.NotNil(t, ownerSession)
	cleanupCreatedUser(t, ownerSession)

	t.Run("create_server", func(t *testing.T) {
		serverID = createServer(t, ownerSession, serverName)
	})

	require.NotZero(t, serverID)
	cleanupCreatedServer(t, ownerSession, serverID)

	t.Run("check_created_server", func(t *testing.T) {
		server := requireServer(t, ownerSession, serverID, serverName, "owner")
		require.Empty(t, server.Topics)
	})

	t.Run("update_server", func(t *testing.T) {
		response := client.MustPostUpdateServer(t, &server_api.PostUpdateServerRequest{
			ServerID: serverID,
			Name:     updatedServerName,
		}, ownerSession)
		response.RequireStatus(t, http.StatusOK)
		require.Empty(t, response.Body)
	})

	t.Run("check_updated_server", func(t *testing.T) {
		requireServer(t, ownerSession, serverID, updatedServerName, "owner")
	})

	t.Run("create_text_topic", func(t *testing.T) {
		textTopicID = createTopic(t, ownerSession, serverID, textTopicName, "text")
	})

	require.NotZero(t, textTopicID)
	cleanupCreatedTopic(t, ownerSession, textTopicID)

	t.Run("create_voice_topic", func(t *testing.T) {
		voiceTopicID = createTopic(t, ownerSession, serverID, voiceTopicName, "voice")
	})

	require.NotZero(t, voiceTopicID)
	cleanupCreatedTopic(t, ownerSession, voiceTopicID)

	t.Run("check_created_topics", func(t *testing.T) {
		server := requireServer(t, ownerSession, serverID, updatedServerName, "owner")
		requireTopic(t, server, textTopicID, textTopicName, "text")
		requireTopic(t, server, voiceTopicID, voiceTopicName, "voice")
	})

	t.Run("update_topic", func(t *testing.T) {
		response := client.MustPostUpdateTopic(t, &server_api.PostUpdateTopicRequest{
			TopicID: textTopicID,
			Name:    updatedTextTopicName,
		}, ownerSession)
		response.RequireStatus(t, http.StatusOK)
		require.Empty(t, response.Body)
	})

	t.Run("check_updated_topic", func(t *testing.T) {
		server := requireServer(t, ownerSession, serverID, updatedServerName, "owner")
		requireTopic(t, server, textTopicID, updatedTextTopicName, "text")
	})

	t.Run("delete_topic", func(t *testing.T) {
		response := client.MustDeleteTopic(t, &server_api.DeleteTopicRequest{
			TopicID: voiceTopicID,
		}, ownerSession)
		response.RequireStatus(t, http.StatusOK)
		require.Empty(t, response.Body)
	})

	t.Run("check_deleted_topic", func(t *testing.T) {
		server := requireServer(t, ownerSession, serverID, updatedServerName, "owner")
		requireTopic(t, server, textTopicID, updatedTextTopicName, "text")
		requireNoTopic(t, server, voiceTopicID)
	})

	t.Run("delete_server", func(t *testing.T) {
		response := client.MustDeleteServer(t, &server_api.DeleteServerRequest{
			ServerID: serverID,
		}, ownerSession)
		response.RequireStatus(t, http.StatusOK)
		require.Empty(t, response.Body)
	})

	t.Run("check_deleted_server", func(t *testing.T) {
		requireNoServer(t, ownerSession, serverID)
	})
}

func Test_ServerManagementAccess(t *testing.T) {
	t.Parallel()

	var (
		ownerSession  *http.Cookie
		memberSession *http.Cookie
		serverID      int64
		topicID       int64
	)

	serverName := uniqueServerName(t, "access")
	updatedServerName := uniqueServerName(t, "accupd")
	topicName := uniqueTopicName(t, "topic")
	updatedTopicName := uniqueTopicName(t, "tupd")

	t.Run("create_users", func(t *testing.T) {
		_, ownerSession = createUser(t, "srvown")
		_, memberSession = createUser(t, "srvmem")
	})

	require.NotNil(t, ownerSession)
	cleanupCreatedUser(t, ownerSession)
	require.NotNil(t, memberSession)
	cleanupCreatedUser(t, memberSession)

	t.Run("reject_create_server_without_session", func(t *testing.T) {
		response := client.MustPostCreateServer(t, &server_api.PostCreateServerRequest{
			Name: serverName,
		})
		response.RequireStatus(t, http.StatusUnauthorized)
		response.RequireNoSessionCookie(t)
		response.RequireErrorContains(t, "login required")
	})

	t.Run("create_server", func(t *testing.T) {
		serverID = createServer(t, ownerSession, serverName)
	})

	require.NotZero(t, serverID)
	cleanupCreatedServer(t, ownerSession, serverID)

	t.Run("join_member", func(t *testing.T) {
		response := client.MustPostJoinServer(t, &server_api.PostJoinServerRequest{
			ServerID: serverID,
		}, memberSession)
		response.RequireStatus(t, http.StatusOK)
		require.Empty(t, response.Body)
	})

	t.Run("owner_creates_topic", func(t *testing.T) {
		topicID = createTopic(t, ownerSession, serverID, topicName, "text")
	})

	require.NotZero(t, topicID)
	cleanupCreatedTopic(t, ownerSession, topicID)

	t.Run("member_sees_server_as_member", func(t *testing.T) {
		server := requireServer(t, memberSession, serverID, serverName, "member")
		requireTopic(t, server, topicID, topicName, "text")
	})

	t.Run("reject_member_update_server", func(t *testing.T) {
		response := client.MustPostUpdateServer(t, &server_api.PostUpdateServerRequest{
			ServerID: serverID,
			Name:     updatedServerName,
		}, memberSession)
		response.RequireStatus(t, http.StatusForbidden)
		response.RequireErrorContains(t, "only owners can perform this action")
		requireServer(t, ownerSession, serverID, serverName, "owner")
	})

	t.Run("reject_member_create_topic", func(t *testing.T) {
		response := client.MustPostCreateTopic(t, &server_api.PostCreateTopicRequest{
			ServerID: serverID,
			Name:     uniqueTopicName(t, "deny"),
			Type:     "voice",
		}, memberSession)
		response.RequireStatus(t, http.StatusForbidden)
		response.RequireErrorContains(t, "only owners can perform this action")
	})

	t.Run("reject_member_update_topic", func(t *testing.T) {
		response := client.MustPostUpdateTopic(t, &server_api.PostUpdateTopicRequest{
			TopicID: topicID,
			Name:    updatedTopicName,
		}, memberSession)
		response.RequireStatus(t, http.StatusForbidden)
		response.RequireErrorContains(t, "only owners can perform this action")

		server := requireServer(t, ownerSession, serverID, serverName, "owner")
		requireTopic(t, server, topicID, topicName, "text")
	})

	t.Run("reject_member_delete_topic", func(t *testing.T) {
		response := client.MustDeleteTopic(t, &server_api.DeleteTopicRequest{
			TopicID: topicID,
		}, memberSession)
		response.RequireStatus(t, http.StatusForbidden)
		response.RequireErrorContains(t, "only owners can perform this action")

		server := requireServer(t, ownerSession, serverID, serverName, "owner")
		requireTopic(t, server, topicID, topicName, "text")
	})

	t.Run("reject_member_delete_server", func(t *testing.T) {
		response := client.MustDeleteServer(t, &server_api.DeleteServerRequest{
			ServerID: serverID,
		}, memberSession)
		response.RequireStatus(t, http.StatusForbidden)
		response.RequireErrorContains(t, "only owners can perform this action")
		requireServer(t, ownerSession, serverID, serverName, "owner")
	})

	t.Run("reject_update_server_without_session", func(t *testing.T) {
		response := client.MustPostUpdateServer(t, &server_api.PostUpdateServerRequest{
			ServerID: serverID,
			Name:     updatedServerName,
		})
		response.RequireStatus(t, http.StatusUnauthorized)
		response.RequireNoSessionCookie(t)
		response.RequireErrorContains(t, "login required")
	})

	t.Run("reject_delete_topic_without_session", func(t *testing.T) {
		response := client.MustDeleteTopic(t, &server_api.DeleteTopicRequest{
			TopicID: topicID,
		})
		response.RequireStatus(t, http.StatusUnauthorized)
		response.RequireNoSessionCookie(t)
		response.RequireErrorContains(t, "login required")
	})

	t.Run("reject_delete_server_without_session", func(t *testing.T) {
		response := client.MustDeleteServer(t, &server_api.DeleteServerRequest{
			ServerID: serverID,
		})
		response.RequireStatus(t, http.StatusUnauthorized)
		response.RequireNoSessionCookie(t)
		response.RequireErrorContains(t, "login required")
	})
}

func requireServer(t testing.TB, session *http.Cookie, serverID int64, name string, role string) server_api.ServerResponse {
	t.Helper()

	response := client.MustGetServers(t, session)
	response.RequireStatus(t, http.StatusOK)

	var body server_api.GetServersResponse
	response.RequireJSON(t, &body)

	for _, server := range body.Servers {
		if server.ID == serverID {
			require.Equal(t, name, server.Name)
			require.Equal(t, role, server.Role)
			return server
		}
	}

	require.FailNowf(t, "server not found", "server_id=%d response=%+v", serverID, body.Servers)
	return server_api.ServerResponse{}
}

func requireNoServer(t testing.TB, session *http.Cookie, serverID int64) {
	t.Helper()

	response := client.MustGetServers(t, session)
	response.RequireStatus(t, http.StatusOK)

	var body server_api.GetServersResponse
	response.RequireJSON(t, &body)

	for _, server := range body.Servers {
		require.NotEqual(t, serverID, server.ID)
	}
}

func requireTopic(t testing.TB, server server_api.ServerResponse, topicID int64, name string, topicType string) {
	t.Helper()

	for _, topic := range server.Topics {
		if topic.ID == topicID {
			require.Equal(t, name, topic.Name)
			require.Equal(t, topicType, topic.Type)
			return
		}
	}

	require.FailNowf(t, "topic not found", "topic_id=%d response=%+v", topicID, server.Topics)
}

func requireNoTopic(t testing.TB, server server_api.ServerResponse, topicID int64) {
	t.Helper()

	for _, topic := range server.Topics {
		require.NotEqual(t, topicID, topic.ID)
	}
}
