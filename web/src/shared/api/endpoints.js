export const API_ENDPOINTS = {
  user: {
    me: "/v1/user/me",
    login: "/v1/user/login",
    create: "/v1/user/create",
  },
  dm: {
    list: "/v1/dm/list",
    searchUsers: "/v1/dm/search",
    create: "/v1/dm/create",
    messages: "/v1/dm/messages",
    sendMessage: "/v1/dm/message",
    markRead: "/v1/dm/read",
  },
  server: {
    list: "/v1/server/list",
    search: "/v1/server/search",
    join: "/v1/server/join",
    leave: "/v1/server/leave",
    create: "/v1/server/create",
    update: "/v1/server/update",
    delete: "/v1/server/delete",
    createTopic: "/v1/server/topic/create",
    updateTopic: "/v1/server/topic/update",
    deleteTopic: "/v1/server/topic/delete",
    topicMessages: "/v1/server/topic/messages",
    sendTopicMessage: "/v1/server/topic/message",
  },
};
