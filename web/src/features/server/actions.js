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
import {
  SERVER_MESSAGES,
  SERVER_QUERY_KEY,
  SERVER_TOPIC_TYPE,
} from "~/features/server/constants";
import { getUserErrorMessage } from "~/shared/api/errors";
import { routes } from "~/shared/config/routes";

export const getServersQuery = query(getServers, SERVER_QUERY_KEY);

const MIN_NAME_LENGTH = 1;

function formString(formData, name) {
  return String(formData.get(name) ?? "").trim();
}

function actionError(error, fallback) {
  return {
    ok: false,
    error: getUserErrorMessage(error, fallback),
  };
}

function validateId(id, fieldName) {
  const value = Number(id);
  if (!Number.isInteger(value) || value <= 0) {
    return `${fieldName} is invalid`;
  }
  return "";
}

function validateName(name) {
  return name.length >= MIN_NAME_LENGTH ? "" : "Name is required";
}

function validateTopicType(type) {
  return Object.values(SERVER_TOPIC_TYPE).includes(type)
    ? ""
    : "Topic type is invalid";
}

/**
 * @typedef {Object} ActionResult
 * @property {boolean} ok
 * @property {string=} error
 */

export const createServerAction = action(async (formData) => {
  const name = formString(formData, "name");
  const error = validateName(name);
  if (error) return { ok: false, error };

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
    const idError = validateId(serverId, "Server");
    if (idError) return { ok: false, error: idError };

    const nameError = validateName(nextName);
    if (nameError) return { ok: false, error: nameError };

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
  const idError = validateId(serverId, "Server");
  if (idError) return { ok: false, error: idError };

  const nameError = validateName(name);
  if (nameError) return { ok: false, error: nameError };

  const typeError = validateTopicType(type);
  if (typeError) return { ok: false, error: typeError };

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
  const idError = validateId(topicId, "Topic");
  if (idError) return { ok: false, error: idError };

  const nameError = validateName(name);
  if (nameError) return { ok: false, error: nameError };

  try {
    await updateTopic(topicId, name);
    await revalidate(getServersQuery.key);
    return { ok: true };
  } catch (error) {
    return actionError(error, SERVER_MESSAGES.updateTopicFailed);
  }
}, "server.topic.update");

export const deleteTopicAction = action(async (topicId) => {
  const idError = validateId(topicId, "Topic");
  if (idError) return { ok: false, error: idError };

  try {
    await deleteTopic(topicId);
    await revalidate(getServersQuery.key);
    return { ok: true };
  } catch (error) {
    return actionError(error, SERVER_MESSAGES.deleteTopicFailed);
  }
}, "server.topic.delete");

export const deleteServerAction = action(async (serverId) => {
  const idError = validateId(serverId, "Server");
  if (idError) return { ok: false, error: idError };

  try {
    await deleteServer(serverId);
    await revalidate(getServersQuery.key);
    return redirect(routes.server.list());
  } catch (error) {
    return actionError(error, SERVER_MESSAGES.deleteFailed);
  }
}, "server.delete");
