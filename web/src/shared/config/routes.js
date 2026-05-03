export const routes = {
  dm: {
    list: () => "/dm",
    chat: (dmId) => `/dm/${dmId}`,
    search: () => "/dm/search",
  },
  server: {
    list: () => "/server",
    manage: () => "/server/manage",
    create: () => "/server/create",
    edit: (serverId) => `/server/${serverId}/edit`,
    textTopic: (serverId, topicId) =>
      `/server/${serverId}/text/${topicId}`,
    voiceTopic: (serverId, topicId) =>
      `/server/${serverId}/voice/${topicId}`,
  },
};
