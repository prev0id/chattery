import { Check, Trash2 } from "lucide-solid";
import { createSignal, For, Show } from "solid-js";
import Button from "~/components/Button";
import FormTextInput from "~/components/FormTextInput";
import Header from "~/components/Header";
import HeaderItem from "~/components/HeaderItem";
import { UseServerContext } from "~/stores/server";

export default function Edit() {
  const { currentServer } = UseServerContext();

  const [name, setName] = createSignal("");
  const [newTopicName, setNewTopicName] = createSignal("");
  const [newTopicType, setNewTopicType] = createSignal("text");
  const [isLoading, setIsLoading] = createSignal(false);

  return (
    <>
      <Header icon="settings">
        <HeaderItem>{currentServer()?.name}</HeaderItem>
      </Header>
      <div class="m-auto max-w-xl min-w-sm w-full border-3 neo-shadow-lg rounded-xl p-4">
        <h2 class="text-lg tracking-wider font-semibold">Update Server Name</h2>
        <form
          // onSubmit={handleUpdateServer}
          class="flex gap-2"
        >
          <FormTextInput
            value={name}
            onInput={(event) => setName(event.currentTarget.value)}
            required
          />
          <Button type="submit" variant="sky" disabled={isLoading()}>
            <Check class="size-5" />
          </Button>
        </form>
        <hr class="my-4" />
        <h2 class="text-lg tracking-wider font-semibold">Add topic</h2>
        <form
          // onSubmit={handleAddTopic}
          class="flex gap-2"
        >
          <select
            class="rounded-lg neo-shadow border-2 px-2 text-center tracking-wider"
            value={newTopicType()}
            onChange={(event) => setNewTopicType(event.currentTarget.value)}
          >
            <option value="text">Text</option>
            <option value="voice">Voice</option>
          </select>
          <FormTextInput
            value={newTopicName}
            onInput={(event) => setNewTopicName(event.currentTarget.value)}
            required
          />
          <Button type="submit" variant="sky" disabled={isLoading()}>
            <Check class="size-5" />
          </Button>
        </form>
        <Show when={currentServer()?.topics?.length > 0}>
          <hr class="my-4" />
          <h2 class="text-lg tracking-wider font-semibold">Update topics</h2>
          <For each={currentServer()?.topics}>
            {(topic) => (
              <form
                // onSubmit={(e) => handleUpdateTopic(e, topic)}
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
                  // onClick={(e) => handleDeleteTopic(e, topic.id)}
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
          // onClick={handleDeleteServer}
          disabled={isLoading()}
        >
          Delete Server
        </Button>
      </div>
    </>
  );
}
