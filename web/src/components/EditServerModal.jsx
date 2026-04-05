import { createSignal, createEffect, For, Show, onMount } from "solid-js";
import Modal from "./Modal";
import FormTextInput from "./FormTextInput";
import Button from "./Button";
import {
  selectedServerForEdit,
  setSelectedServerForEdit,
} from "../stores/serverState";
import {
  updateServer,
  createTopic,
  updateTopic,
  deleteTopic,
  deleteServer,
} from "../lib/api";

import { Trash2, Check } from "lucide-solid";

export default function EditServerModal(props) {
  const [name, setName] = createSignal("");
  const [newTopicName, setNewTopicName] = createSignal("");
  const [newTopicType, setNewTopicType] = createSignal("text");
  const [isLoading, setIsLoading] = createSignal(false);

  createEffect(() => {
    const server = selectedServerForEdit();
    if (server) {
      setName(server.name);
    }
  });

  const handleUpdateServer = async (e) => {
    e.preventDefault();
    const server = selectedServerForEdit();
    if (!server) return;
    setIsLoading(true);
    const result = await updateServer(server.id, name());
    if (result) {
      setSelectedServerForEdit({ ...server, name: name() });
    }
    setIsLoading(false);
  };

  const handleAddTopic = async (e) => {
    e.preventDefault();
    const server = selectedServerForEdit();
    if (!server || !newTopicName()) return;
    setIsLoading(true);
    const result = await createTopic(server.id, newTopicName(), newTopicType());
    if (result) {
      setNewTopicName("");
      const updatedTopics = [
        { id: result.id, name: newTopicName(), type: newTopicType() },
        ...(server.topics || []),
      ];
      setSelectedServerForEdit({ ...server, topics: updatedTopics });
    }
    setIsLoading(false);
  };

  const handleUpdateTopic = async (e, topic) => {
    e.preventDefault();
    const input = e.currentTarget.querySelector('input[type="text"]');
    if (!input) return;
    const newName = input.value;
    if (!newName || newName === topic.name) return;
    setIsLoading(true);
    await updateTopic(topic.id, newName);
    const server = selectedServerForEdit();
    if (server) {
      const updatedTopics = server.topics.map((t) =>
        t.id === topic.id ? { ...t, name: newName } : t,
      );
      setSelectedServerForEdit({ ...server, topics: updatedTopics });
    }
    setIsLoading(false);
  };

  const handleDeleteTopic = async (e, topicId) => {
    e.preventDefault();
    if (!confirm("Are you sure you want to delete this topic?")) return;
    setIsLoading(true);
    await deleteTopic(topicId);
    const server = selectedServerForEdit();
    if (server) {
      const updatedTopics = server.topics.filter((t) => t.id !== topicId);
      setSelectedServerForEdit({ ...server, topics: updatedTopics });
    }
    setIsLoading(false);
  };

  const handleDeleteServer = async (e) => {
    e.preventDefault();
    const server = selectedServerForEdit();
    if (!server) return;
    if (!confirm("Are you sure you want to delete this server?")) return;
    setIsLoading(true);
    const result = await deleteServer(server.id);
    if (result) {
      setSelectedServerForEdit(null);
      document.getElementById(props.id)?.hidePopover();
      window.location.reload();
    }
    setIsLoading(false);
  };

  return (
    <Modal id={props.id} name="Edit Server">
      <div class="mt-4">
        <Show when={selectedServerForEdit()}>
          <h2 class="text-lg tracking-wider font-semibold">
            Update Server Name
          </h2>
          <form onSubmit={handleUpdateServer} class="flex gap-2">
            <FormTextInput
              value={name}
              onInput={(e) => setName(e.currentTarget.value)}
              required
            />
            <Button type="submit" variant="sky" disabled={isLoading()}>
              <Check class="size-5" />
            </Button>
          </form>
          <hr class="my-4" />
          <h2 class="text-lg tracking-wider font-semibold">Add topic</h2>
          <form onSubmit={handleAddTopic} class="flex gap-2">
            <select
              class="rounded-lg neo-shadow border-2 px-2 text-center tracking-wider"
              value={newTopicType()}
              onChange={(e) => setNewTopicType(e.currentTarget.value)}
            >
              <option value="text">Text</option>
              <option value="voice">Voice</option>
            </select>
            <FormTextInput
              value={newTopicName}
              onInput={(e) => setNewTopicName(e.currentTarget.value)}
              required
            />
            <Button type="submit" variant="sky" disabled={isLoading()}>
              <Check class="size-5" />
            </Button>
          </form>
          <Show when={selectedServerForEdit().topics?.length > 0}>
            <hr class="my-4" />
            <h2 class="text-lg tracking-wider font-semibold">Update topics</h2>
            <For each={selectedServerForEdit().topics}>
              {(topic) => (
                <form
                  onSubmit={(e) => handleUpdateTopic(e, topic)}
                  class="flex gap-2 mb-2"
                >
                  <div
                    class="rounded-lg neo-shadow border-2 px-2 text-center tracking-wider w-16 flex items-center justify-center"
                    classList={{
                      "bg-sky-200": topic.type === "text",
                      "bg-amber-200": topic.type === "voice",
                    }}
                  >
                    {topic.type}
                  </div>

                  <input
                    class="px-2 border-2 neo-shadow rounded-lg focus:outline-none focus:border-sky-500 flex-1"
                    type="text"
                    value={topic.name}
                  />
                  <Button type="submit" variant="sky" disabled={isLoading()}>
                    <Check class="size-5" />
                  </Button>
                  <Button
                    type="button"
                    variant="rose"
                    disabled={isLoading()}
                    onClick={(e) => handleDeleteTopic(e, topic.id)}
                  >
                    <Trash2 class="size-5" />
                  </Button>
                </form>
              )}
            </For>
          </Show>

          <hr class="my-4" />
          <Button
            variant="rose"
            onClick={handleDeleteServer}
            disabled={isLoading()}
          >
            Delete Server
          </Button>
        </Show>
      </div>
    </Modal>
  );
}
