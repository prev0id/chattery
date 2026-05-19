//go:build e2e

package e2e

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/stretchr/testify/require"

	"chattery/e2e/client"
	dm_api "chattery/internal/api/dm"
	server_api "chattery/internal/api/server"
	user_api "chattery/internal/api/user"
	"chattery/internal/api/websocket/event_desc"
)

const (
	wsReadTimeout    = 2 * time.Second
	wsNoEventTimeout = 500 * time.Millisecond
	wsJoinWait       = 200 * time.Millisecond
)

func Test_WebSocketConnection(t *testing.T) {
	t.Parallel()

	var session *http.Cookie

	t.Run("create_user", func(t *testing.T) {
		_, session = createUser(t, "wsconn")
	})

	require.NotNil(t, session)
	cleanupCreatedUser(t, session)

	t.Run("reject_without_session", func(t *testing.T) {
		response := client.MustRejectWebsocket(t, http.StatusUnauthorized)
		response.RequireNoSessionCookie(t)
		response.RequireErrorContains(t, "login required")
	})

	t.Run("connect_with_session", func(t *testing.T) {
		conn := client.MustDialWebsocket(t, session)
		client.MustCloseWebsocket(t, conn)
	})
}

func Test_WebSocketJoinLeave(t *testing.T) {
	t.Parallel()

	var (
		owner           user_api.User
		member          user_api.User
		ownerSession    *http.Cookie
		memberSession   *http.Cookie
		outsiderSession *http.Cookie
		serverID        int64
		topicID         int64
		ownerConn       *websocket.Conn
	)

	serverName := uniqueServerName(t, "wssrv")
	topicName := uniqueTopicName(t, "ws")
	messageText := uniqueMessageTexts(t, "wsjl", 1)[0]
	messageAfterLeave := uniqueMessageTexts(t, "wslv", 1)[0]

	t.Run("create_users", func(t *testing.T) {
		_, ownerSession = createUser(t, "wsjo")
		_, memberSession = createUser(t, "wsjm")
		_, outsiderSession = createUser(t, "wsjx")
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

	t.Run("create_server_and_topic", func(t *testing.T) {
		serverID = createServer(t, ownerSession, serverName)
		topicID = createTopic(t, ownerSession, serverID, topicName, "text")
	})

	require.NotZero(t, serverID)
	cleanupCreatedServer(t, ownerSession, serverID)
	require.NotZero(t, topicID)
	cleanupCreatedTopic(t, ownerSession, topicID)

	t.Run("join_member_to_server", func(t *testing.T) {
		response := client.MustPostJoinServer(t, &server_api.PostJoinServerRequest{
			ServerID: serverID,
		}, memberSession)
		response.RequireStatus(t, http.StatusOK)
		require.Empty(t, response.Body)
	})

	t.Run("owner_joins_topic", func(t *testing.T) {
		ownerConn = client.MustDialWebsocket(t, ownerSession)
		wsJoinTextTopic(t, ownerConn, topicID)
	})

	require.NotNil(t, ownerConn)
	t.Cleanup(func() { client.MustCloseWebsocket(t, ownerConn) })

	t.Run("owner_receives_topic_message", func(t *testing.T) {
		postTopicMessage(t, memberSession, topicID, messageText)
		requireWSMessage(t, ownerConn, textTopicChannel(topicID), messageText, member)
	})

	t.Run("owner_leaves_topic", func(t *testing.T) {
		wsLeave(t, ownerConn)
	})

	t.Run("owner_does_not_receive_after_leave", func(t *testing.T) {
		postTopicMessage(t, memberSession, topicID, messageAfterLeave)
		client.RequireNoWSEvent(t, ownerConn, wsNoEventTimeout)
	})

	t.Run("outsider_cannot_join_topic", func(t *testing.T) {
		outsiderConn := client.MustDialWebsocket(t, outsiderSession)
		t.Cleanup(func() { client.MustCloseWebsocket(t, outsiderConn) })

		wsJoinTextTopic(t, outsiderConn, topicID)
		requireWSError(t, outsiderConn, "no access to topic")
	})
}

func Test_WebSocketDMMessages(t *testing.T) {
	t.Parallel()

	var (
		userA           user_api.User
		userB           user_api.User
		userASession    *http.Cookie
		userBSession    *http.Cookie
		outsiderSession *http.Cookie
		dmID            int64
	)

	messageText := uniqueMessageTexts(t, "wsdm", 1)[0]

	t.Run("create_users", func(t *testing.T) {
		_, userASession = createUser(t, "wsda")
		_, userBSession = createUser(t, "wsdb")
		_, outsiderSession = createUser(t, "wsdx")
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
		require.NotZero(t, dmID)
	})

	var userAConn *websocket.Conn
	t.Run("first_user_joins_dm", func(t *testing.T) {
		userAConn = client.MustDialWebsocket(t, userASession)
		wsJoinDM(t, userAConn, dmID)
	})

	require.NotNil(t, userAConn)
	t.Cleanup(func() { client.MustCloseWebsocket(t, userAConn) })

	t.Run("first_user_receives_dm_message", func(t *testing.T) {
		postDMMessage(t, userBSession, dmID, messageText)
		requireWSMessage(t, userAConn, dmChannel(dmID), messageText, userB)
	})

	t.Run("outsider_cannot_join_dm", func(t *testing.T) {
		outsiderConn := client.MustDialWebsocket(t, outsiderSession)
		t.Cleanup(func() { client.MustCloseWebsocket(t, outsiderConn) })

		wsJoinDM(t, outsiderConn, dmID)
		requireWSError(t, outsiderConn, "no access to dm")
	})
}

func Test_WebSocketDMNotificationFromAnotherChat(t *testing.T) {
	t.Parallel()

	var (
		userA        user_api.User
		userB        user_api.User
		userC        user_api.User
		userASession *http.Cookie
		userBSession *http.Cookie
		userCSession *http.Cookie
		dmABID       int64
		dmACID       int64
		userAConn    *websocket.Conn
	)
	messageText := uniqueMessageTexts(t, "wsdn", 1)[0]

	t.Run("create_users", func(t *testing.T) {
		_, userASession = createUser(t, "wsna")
		_, userBSession = createUser(t, "wsnb")
		_, userCSession = createUser(t, "wsnc")
		userA = getMe(t, userASession)
		userB = getMe(t, userBSession)
		userC = getMe(t, userCSession)
	})

	require.NotNil(t, userASession)
	cleanupCreatedUser(t, userASession)
	require.NotNil(t, userBSession)
	cleanupCreatedUser(t, userBSession)
	require.NotNil(t, userCSession)
	cleanupCreatedUser(t, userCSession)
	require.NotZero(t, userA.ID)
	require.NotZero(t, userB.ID)
	require.NotZero(t, userC.ID)

	t.Run("create_dms", func(t *testing.T) {
		dmABID = createDM(t, userASession, userB.ID)
		dmACID = createDM(t, userASession, userC.ID)
	})

	require.NotZero(t, dmABID)
	require.NotZero(t, dmACID)

	t.Run("first_user_joins_first_dm", func(t *testing.T) {
		userAConn = client.MustDialWebsocket(t, userASession)
		wsJoinDM(t, userAConn, dmABID)
	})

	require.NotNil(t, userAConn)
	t.Cleanup(func() { client.MustCloseWebsocket(t, userAConn) })

	t.Run("first_user_receives_second_dm_notification", func(t *testing.T) {
		postDMMessage(t, userCSession, dmACID, messageText)
		requireWSMessage(t, userAConn, dmChannel(dmACID), messageText, userC)
	})
}

func Test_WebSocketTextTopicFiltering(t *testing.T) {
	t.Parallel()

	var (
		owner         user_api.User
		member        user_api.User
		ownerSession  *http.Cookie
		memberSession *http.Cookie
		serverID      int64
		topicAID      int64
		topicBID      int64
		ownerConn     *websocket.Conn
	)

	serverName := uniqueServerName(t, "wsflt")
	topicAName := uniqueTopicName(t, "wa")
	topicBName := uniqueTopicName(t, "wb")
	topicAText := uniqueMessageTexts(t, "wsta", 1)[0]
	topicBText := uniqueMessageTexts(t, "wstb", 1)[0]

	t.Run("create_users", func(t *testing.T) {
		_, ownerSession = createUser(t, "wsfo")
		_, memberSession = createUser(t, "wsfm")
		owner = getMe(t, ownerSession)
		member = getMe(t, memberSession)
	})

	require.NotNil(t, ownerSession)
	cleanupCreatedUser(t, ownerSession)
	require.NotNil(t, memberSession)
	cleanupCreatedUser(t, memberSession)
	require.NotZero(t, owner.ID)
	require.NotZero(t, member.ID)

	t.Run("create_server_and_topics", func(t *testing.T) {
		serverID = createServer(t, ownerSession, serverName)
		topicAID = createTopic(t, ownerSession, serverID, topicAName, "text")
		topicBID = createTopic(t, ownerSession, serverID, topicBName, "text")
	})

	require.NotZero(t, serverID)
	cleanupCreatedServer(t, ownerSession, serverID)
	require.NotZero(t, topicAID)
	cleanupCreatedTopic(t, ownerSession, topicAID)
	require.NotZero(t, topicBID)
	cleanupCreatedTopic(t, ownerSession, topicBID)

	t.Run("join_member_to_server", func(t *testing.T) {
		response := client.MustPostJoinServer(t, &server_api.PostJoinServerRequest{
			ServerID: serverID,
		}, memberSession)
		response.RequireStatus(t, http.StatusOK)
		require.Empty(t, response.Body)
	})

	t.Run("owner_joins_first_topic", func(t *testing.T) {
		ownerConn = client.MustDialWebsocket(t, ownerSession)
		wsJoinTextTopic(t, ownerConn, topicAID)
	})

	require.NotNil(t, ownerConn)
	t.Cleanup(func() { client.MustCloseWebsocket(t, ownerConn) })

	t.Run("owner_receives_joined_topic_message", func(t *testing.T) {
		postTopicMessage(t, memberSession, topicAID, topicAText)
		requireWSMessage(t, ownerConn, textTopicChannel(topicAID), topicAText, member)
	})

	t.Run("owner_does_not_receive_other_topic_message", func(t *testing.T) {
		postTopicMessage(t, memberSession, topicBID, topicBText)
		client.RequireNoWSEvent(t, ownerConn, wsNoEventTimeout)
	})
}

func wsJoinDM(t testing.TB, conn *websocket.Conn, dmID int64) {
	t.Helper()

	wsJoin(t, conn, dmChannel(dmID))
}

func wsJoinTextTopic(t testing.TB, conn *websocket.Conn, topicID int64) {
	t.Helper()

	wsJoin(t, conn, textTopicChannel(topicID))
}

func wsJoin(t testing.TB, conn *websocket.Conn, channel event_desc.Channel) {
	t.Helper()

	client.MustWriteWSEvent(t, conn, &event_desc.Event{
		Type:    event_desc.TypeJoin,
		Payload: wsPayload(t, channel),
	})
	time.Sleep(wsJoinWait)
}

func wsLeave(t testing.TB, conn *websocket.Conn) {
	t.Helper()

	client.MustWriteWSEvent(t, conn, &event_desc.Event{
		Type: event_desc.TypeLeave,
	})
	time.Sleep(wsJoinWait)
}

func postDMMessage(t testing.TB, session *http.Cookie, dmID int64, text string) {
	t.Helper()

	response := client.MustPostDMMessage(t, &dm_api.PostMessageRequest{
		DMID: dmID,
		Text: text,
	}, session)
	response.RequireStatus(t, http.StatusOK)
	require.Empty(t, response.Body)
}

func postTopicMessage(t testing.TB, session *http.Cookie, topicID int64, text string) {
	t.Helper()

	response := client.MustPostTopicMessage(t, &server_api.PostMessageRequest{
		TopicID: topicID,
		Text:    text,
	}, session)
	response.RequireStatus(t, http.StatusOK)
	require.Empty(t, response.Body)
}

func requireWSMessage(
	t testing.TB,
	conn *websocket.Conn,
	channel event_desc.Channel,
	text string,
	sender user_api.User,
) event_desc.MessagePayload {
	t.Helper()

	event := client.MustReadWSEvent(t, conn, wsReadTimeout)
	require.Equal(t, event_desc.TypeMessage, event.Type)
	require.Equal(t, channel, event.Channel)

	var payload event_desc.MessagePayload
	require.NoError(t, json.Unmarshal(event.Payload, &payload))
	require.NotZero(t, payload.ID)
	require.NotEmpty(t, payload.CreatedAt)
	require.Equal(t, text, payload.Text)
	require.Equal(t, sender.ID, payload.Sender.ID)
	require.Equal(t, sender.Username, payload.Sender.Username)

	return payload
}

func requireWSError(t testing.TB, conn *websocket.Conn, expected string) {
	t.Helper()

	event := client.MustReadWSEvent(t, conn, wsReadTimeout)
	require.Equal(t, event_desc.TypeError, event.Type)

	var payload event_desc.ErrorPayload
	require.NoError(t, json.Unmarshal(event.Payload, &payload))
	require.Contains(t, payload.Message, expected)
}

func wsPayload(t testing.TB, value any) json.RawMessage {
	t.Helper()

	raw, err := json.Marshal(value)
	require.NoError(t, err)

	return raw
}

func dmChannel(dmID int64) event_desc.Channel {
	return event_desc.Channel{
		Type: event_desc.ChannelDM,
		ID:   dmID,
	}
}

func textTopicChannel(topicID int64) event_desc.Channel {
	return event_desc.Channel{
		Type: event_desc.ChannelTextTopic,
		ID:   topicID,
	}
}
