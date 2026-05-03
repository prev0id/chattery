import { action, query, redirect, revalidate } from "@solidjs/router";
import {
  createServer,
  createTopic,
  deleteServer,
  deleteTopic,
  getServers,
  updateServer,
  updateTopic,
} from "~/features/server/api";
import { SERVER_MESSAGES, SERVER_QUERY_KEY } from "~/features/server/constants";
import { getUserErrorMessage } from "~/shared/api/errors";
import { routes } from "~/shared/config/routes";

export const getServersQuery = query(getServers, SERVER_QUERY_KEY);

function formString(formData, name) {
  return String(formData.get(name) ?? "").trim();
}

function actionError(error, fallback) {
  return {
    ok: false,
    error: getUserErrorMessage(error, fallback),
  };
}

export const createServerAction = action(async (formData) => {
  const name = formString(formData, "name");

  try {
    const server = await createServer(name);
    await revalidate(getServersQuery.key);
    return redirect(routes.server.edit(server.id));
  } catch (error) {
    return actionError(error, SERVER_MESSAGES.createFailed);
  }
}, "server.create");

export const updateServerAction = action(
  async (serverId, currentName, formData) => {
    const nextName = formString(formData, "name");

    if (nextName === currentName) {
      return { ok: false, error: SERVER_MESSAGES.unchangedName };
    }

    try {
      await updateServer(serverId, nextName);
      await revalidate(getServersQuery.key);
      return { ok: true };
    } catch (error) {
      return actionError(error, SERVER_MESSAGES.updateFailed);
    }
  },
  "server.update",
);

export const addTopicAction = action(async (serverId, formData) => {
  const name = formString(formData, "name");
  const type = formString(formData, "type");

  try {
    await createTopic(serverId, name, type);
    await revalidate(getServersQuery.key);
    return { ok: true };
  } catch (error) {
    return actionError(error, SERVER_MESSAGES.addTopicFailed);
  }
}, "server.topic.create");

export const updateTopicAction = action(async (topicId, formData) => {
  const name = formString(formData, "name");

  try {
    await updateTopic(topicId, name);
    await revalidate(getServersQuery.key);
    return { ok: true };
  } catch (error) {
    return actionError(error, SERVER_MESSAGES.updateTopicFailed);
  }
}, "server.topic.update");

export const deleteTopicAction = action(async (topicId) => {
  try {
    await deleteTopic(topicId);
    await revalidate(getServersQuery.key);
    return { ok: true };
  } catch (error) {
    return actionError(error, SERVER_MESSAGES.deleteTopicFailed);
  }
}, "server.topic.delete");

export const deleteServerAction = action(async (serverId) => {
  try {
    await deleteServer(serverId);
    await revalidate(getServersQuery.key);
    return redirect(routes.server.list());
  } catch (error) {
    return actionError(error, SERVER_MESSAGES.deleteFailed);
  }
}, "server.delete");
