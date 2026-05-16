//go:build e2e

package e2e

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"chattery/e2e/client"
	dm_api "chattery/internal/api/dm"
	server_api "chattery/internal/api/server"
	user_api "chattery/internal/api/user"
)

const messagesPageLimit = 20

func Test_TextTopicMessages(t *testing.T) {
	t.Parallel()

	var owner user_api.User
	var member user_api.User
	var ownerSession *http.Cookie
	var memberSession *http.Cookie
	var outsiderSession *http.Cookie
	var serverID int64
	var textTopicID int64
	var voiceTopicID int64
	var firstPage server_api.PostTopicMessagesResponse
	serverName := uniqueServerName(t, "msgsrv")
	textTopicName := uniqueTopicName(t, "mt")
	voiceTopicName := uniqueTopicName(t, "mv")
	messageTexts := uniqueMessageTexts(t, "tm", messagesPageLimit+4)

	t.Run("create_users", func(t *testing.T) {
		_, ownerSession = createUser(t, "msgown")
		_, memberSession = createUser(t, "msgmem")
		_, outsiderSession = createUser(t, "msgout")
		owner = getMe(t, ownerSession)
		member = getMe(t, memberSession)
	})
	require.NotNil(t, ownerSession)
	cleanupCreatedUser(t, ownerSession)
	require.NotNil(t, memberSession)
	cleanupCreatedUser(t, memberSession)
	require.NotNil(t, outsiderSession)
	cleanupCreatedUser(t, outsiderSession)
	require.NotZero(t, owner.ID)
	require.NotZero(t, member.ID)

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

	t.Run("create_topics", func(t *testing.T) {
		textTopicID = createTopic(t, ownerSession, serverID, textTopicName, "text")
		voiceTopicID = createTopic(t, ownerSession, serverID, voiceTopicName, "voice")
	})
	require.NotZero(t, textTopicID)
	cleanupCreatedTopic(t, ownerSession, textTopicID)
	require.NotZero(t, voiceTopicID)
	cleanupCreatedTopic(t, ownerSession, voiceTopicID)

	t.Run("reject_send_without_session", func(t *testing.T) {
		response := client.MustPostTopicMessage(t, &server_api.PostMessageRequest{
			TopicID: textTopicID,
			Text:    "without session",
		})
		response.RequireStatus(t, http.StatusUnauthorized)
		response.RequireNoSessionCookie(t)
		response.RequireErrorContains(t, "login required")
	})

	t.Run("reject_read_without_session", func(t *testing.T) {
		response := client.MustPostTopicMessages(t, &server_api.PostTopicMessagesRequest{
			Cursor: &server_api.Cursor{
				TopicID: textTopicID,
			},
		})
		response.RequireStatus(t, http.StatusUnauthorized)
		response.RequireNoSessionCookie(t)
		response.RequireErrorContains(t, "login required")
	})

	t.Run("reject_outsider_send", func(t *testing.T) {
		response := client.MustPostTopicMessage(t, &server_api.PostMessageRequest{
			TopicID: textTopicID,
			Text:    "outsider message",
		}, outsiderSession)
		response.RequireStatus(t, http.StatusNotFound)
		response.RequireErrorContains(t, "not a participant")
	})

	t.Run("reject_outsider_read", func(t *testing.T) {
		response := client.MustPostTopicMessages(t, &server_api.PostTopicMessagesRequest{
			Cursor: &server_api.Cursor{
				TopicID: textTopicID,
			},
		}, outsiderSession)
		response.RequireStatus(t, http.StatusNotFound)
		response.RequireErrorContains(t, "not a participant")
	})

	t.Run("reject_voice_topic_send", func(t *testing.T) {
		response := client.MustPostTopicMessage(t, &server_api.PostMessageRequest{
			TopicID: voiceTopicID,
			Text:    "voice topic message",
		}, ownerSession)
		response.RequireStatus(t, http.StatusBadRequest)
		response.RequireErrorContains(t, "type must be text")
	})

	t.Run("reject_voice_topic_read", func(t *testing.T) {
		response := client.MustPostTopicMessages(t, &server_api.PostTopicMessagesRequest{
			Cursor: &server_api.Cursor{
				TopicID: voiceTopicID,
			},
		}, ownerSession)
		response.RequireStatus(t, http.StatusBadRequest)
		response.RequireErrorContains(t, "type must be text")
	})

	t.Run("owner_sends_messages", func(t *testing.T) {
		for _, text := range messageTexts[:len(messageTexts)-1] {
			response := client.MustPostTopicMessage(t, &server_api.PostMessageRequest{
				TopicID: textTopicID,
				Text:    text,
			}, ownerSession)
			response.RequireStatus(t, http.StatusOK)
			require.Empty(t, response.Body)
		}
	})

	t.Run("member_sends_message", func(t *testing.T) {
		response := client.MustPostTopicMessage(t, &server_api.PostMessageRequest{
			TopicID: textTopicID,
			Text:    messageTexts[len(messageTexts)-1],
		}, memberSession)
		response.RequireStatus(t, http.StatusOK)
		require.Empty(t, response.Body)
	})

	t.Run("owner_reads_first_page", func(t *testing.T) {
		response := client.MustPostTopicMessages(t, &server_api.PostTopicMessagesRequest{
			Cursor: &server_api.Cursor{
				TopicID: textTopicID,
			},
		}, ownerSession)
		response.RequireStatus(t, http.StatusOK)
		response.RequireJSON(t, &firstPage)
		require.NotNil(t, firstPage.Cursor)
		require.Equal(t, textTopicID, firstPage.Cursor.TopicID)
		require.NotZero(t, firstPage.Cursor.MessageID)
		require.False(t, firstPage.Cursor.Timestamp.IsZero())
		requireTopicMessageTexts(t, firstPage.Messages, firstPageExpectedTexts(messageTexts))
		require.Equal(t, member.ID, firstPage.Messages[0].SenderID)
		require.Equal(t, member.Username, firstPage.Messages[0].Sender.Username)
	})
	require.NotNil(t, firstPage.Cursor)

	t.Run("member_reads_same_first_page", func(t *testing.T) {
		response := client.MustPostTopicMessages(t, &server_api.PostTopicMessagesRequest{
			Cursor: &server_api.Cursor{
				TopicID: textTopicID,
			},
		}, memberSession)
		response.RequireStatus(t, http.StatusOK)

		var body server_api.PostTopicMessagesResponse
		response.RequireJSON(t, &body)
		requireTopicMessageTexts(t, body.Messages, firstPageExpectedTexts(messageTexts))
	})

	t.Run("owner_reads_next_page", func(t *testing.T) {
		response := client.MustPostTopicMessages(t, &server_api.PostTopicMessagesRequest{
			Cursor: firstPage.Cursor,
		}, ownerSession)
		response.RequireStatus(t, http.StatusOK)

		var body server_api.PostTopicMessagesResponse
		response.RequireJSON(t, &body)
		require.Nil(t, body.Cursor)
		requireTopicMessageTexts(t, body.Messages, nextPageExpectedTexts(messageTexts))
		require.Equal(t, owner.ID, body.Messages[0].SenderID)
		require.Equal(t, owner.Username, body.Messages[0].Sender.Username)
	})
}

func Test_DMMessages(t *testing.T) {
	t.Parallel()

	var userA user_api.User
	var userB user_api.User
	var userASession *http.Cookie
	var userBSession *http.Cookie
	var outsiderSession *http.Cookie
	var dmID int64
	var firstPage dm_api.PostMessagesResponse
	messageTexts := uniqueMessageTexts(t, "dm", messagesPageLimit+4)

	t.Run("create_users", func(t *testing.T) {
		_, userASession = createUser(t, "dmusera")
		_, userBSession = createUser(t, "dmuserb")
		_, outsiderSession = createUser(t, "dmout")
		userA = getMe(t, userASession)
		userB = getMe(t, userBSession)
	})
	require.NotNil(t, userASession)
	cleanupCreatedUser(t, userASession)
	require.NotNil(t, userBSession)
	cleanupCreatedUser(t, userBSession)
	require.NotNil(t, outsiderSession)
	cleanupCreatedUser(t, outsiderSession)
	require.NotZero(t, userA.ID)
	require.NotZero(t, userB.ID)

	t.Run("create_dm", func(t *testing.T) {
		dmID = createDM(t, userASession, userB.ID)
	})
	require.NotZero(t, dmID)

	t.Run("reject_send_without_session", func(t *testing.T) {
		response := client.MustPostDMMessage(t, &dm_api.PostMessageRequest{
			DMID: dmID,
			Text: "without session",
		})
		response.RequireStatus(t, http.StatusUnauthorized)
		response.RequireNoSessionCookie(t)
		response.RequireErrorContains(t, "login required")
	})

	t.Run("reject_read_without_session", func(t *testing.T) {
		response := client.MustPostDMMessages(t, &dm_api.PostMessagesRequest{
			Cursor: &dm_api.Cursor{
				DMID: dmID,
			},
		})
		response.RequireStatus(t, http.StatusUnauthorized)
		response.RequireNoSessionCookie(t)
		response.RequireErrorContains(t, "login required")
	})

	t.Run("reject_outsider_send", func(t *testing.T) {
		response := client.MustPostDMMessage(t, &dm_api.PostMessageRequest{
			DMID: dmID,
			Text: "outsider message",
		}, outsiderSession)
		response.RequireStatus(t, http.StatusNotFound)
		response.RequireErrorContains(t, "not a participant")
	})

	t.Run("reject_outsider_read", func(t *testing.T) {
		response := client.MustPostDMMessages(t, &dm_api.PostMessagesRequest{
			Cursor: &dm_api.Cursor{
				DMID: dmID,
			},
		}, outsiderSession)
		response.RequireStatus(t, http.StatusNotFound)
		response.RequireErrorContains(t, "not a participant")
	})

	t.Run("first_user_sends_messages", func(t *testing.T) {
		for _, text := range messageTexts[:len(messageTexts)-1] {
			response := client.MustPostDMMessage(t, &dm_api.PostMessageRequest{
				DMID: dmID,
				Text: text,
			}, userASession)
			response.RequireStatus(t, http.StatusOK)
			require.Empty(t, response.Body)
		}
	})

	t.Run("second_user_sends_message", func(t *testing.T) {
		response := client.MustPostDMMessage(t, &dm_api.PostMessageRequest{
			DMID: dmID,
			Text: messageTexts[len(messageTexts)-1],
		}, userBSession)
		response.RequireStatus(t, http.StatusOK)
		require.Empty(t, response.Body)
	})

	t.Run("second_user_reads_first_page", func(t *testing.T) {
		response := client.MustPostDMMessages(t, &dm_api.PostMessagesRequest{
			Cursor: &dm_api.Cursor{
				DMID: dmID,
			},
		}, userBSession)
		response.RequireStatus(t, http.StatusOK)
		response.RequireJSON(t, &firstPage)
		require.NotNil(t, firstPage.Cursor)
		require.Equal(t, dmID, firstPage.Cursor.DMID)
		require.NotZero(t, firstPage.Cursor.MessageID)
		require.False(t, firstPage.Cursor.Timestamp.IsZero())
		requireDMMessageTexts(t, firstPage.Messages, firstPageExpectedTexts(messageTexts))
		require.Equal(t, userB.ID, firstPage.Messages[0].Sender.ID)
		require.Equal(t, userB.Username, firstPage.Messages[0].Sender.Username)
	})
	require.NotNil(t, firstPage.Cursor)

	t.Run("first_user_reads_same_first_page", func(t *testing.T) {
		response := client.MustPostDMMessages(t, &dm_api.PostMessagesRequest{
			Cursor: &dm_api.Cursor{
				DMID: dmID,
			},
		}, userASession)
		response.RequireStatus(t, http.StatusOK)

		var body dm_api.PostMessagesResponse
		response.RequireJSON(t, &body)
		requireDMMessageTexts(t, body.Messages, firstPageExpectedTexts(messageTexts))
	})

	t.Run("first_user_reads_next_page", func(t *testing.T) {
		response := client.MustPostDMMessages(t, &dm_api.PostMessagesRequest{
			Cursor: firstPage.Cursor,
		}, userASession)
		response.RequireStatus(t, http.StatusOK)

		var body dm_api.PostMessagesResponse
		response.RequireJSON(t, &body)
		require.Nil(t, body.Cursor)
		requireDMMessageTexts(t, body.Messages, nextPageExpectedTexts(messageTexts))
		require.Equal(t, userA.ID, body.Messages[0].Sender.ID)
		require.Equal(t, userA.Username, body.Messages[0].Sender.Username)
	})
}

func uniqueMessageTexts(t testing.TB, prefix string, count int) []string {
	t.Helper()

	batch := uniqueName(t, prefix, 20)
	texts := make([]string, 0, count)
	for i := 1; i <= count; i++ {
		texts = append(texts, fmt.Sprintf("%s message %02d", batch, i))
	}

	return texts
}

func firstPageExpectedTexts(texts []string) []string {
	return reversedStrings(texts[len(texts)-messagesPageLimit:])
}

func nextPageExpectedTexts(texts []string) []string {
	return reversedStrings(texts[:len(texts)-messagesPageLimit])
}

func reversedStrings(values []string) []string {
	result := make([]string, 0, len(values))
	for i := len(values) - 1; i >= 0; i-- {
		result = append(result, values[i])
	}

	return result
}

func requireTopicMessageTexts(t testing.TB, messages []server_api.Message, expected []string) {
	t.Helper()

	require.Len(t, messages, len(expected))
	for i, expectedText := range expected {
		require.NotZero(t, messages[i].ID)
		require.NotZero(t, messages[i].SenderID)
		require.NotEmpty(t, messages[i].CreatedAt)
		require.Equal(t, expectedText, messages[i].Text)
	}
}

func requireDMMessageTexts(t testing.TB, messages []*dm_api.Message, expected []string) {
	t.Helper()

	require.Len(t, messages, len(expected))
	for i, expectedText := range expected {
		require.NotZero(t, messages[i].ID)
		require.NotZero(t, messages[i].Sender.ID)
		require.NotEmpty(t, messages[i].CreatedAt)
		require.Equal(t, expectedText, messages[i].Text)
	}
}
