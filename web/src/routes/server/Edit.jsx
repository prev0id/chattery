import { Check, Trash2, X } from "lucide-solid";
import { createSignal, For, Show } from "solid-js";
import Button from "~/components/Button";
import FormTextInput from "~/components/FormTextInput";
import Header from "~/components/Header";
import HeaderItem from "~/components/HeaderItem";
import { UseServerContext } from "~/stores/server";
import {
  updateServerAction,
  addTopicAction,
  updateTopicAction,
  deleteTopicAction,
  deleteServerAction,
} from "~/lib/api";
import { useAction, useSubmission, useSubmissions } from "@solidjs/router";

export default function Edit() {
  const { currentServer } = UseServerContext();

  const serverName = () => currentServer()?.name;

  return (
    <>
      <Header icon="settings">
        <HeaderItem>{serverName()}</HeaderItem>
      </Header>
      <div class="m-auto max-w-xl min-w-sm w-full border-3 neo-shadow-lg rounded-xl p-4">
        <SectionHeader>Update Server Name</SectionHeader>
        <UpdateServerForm server={currentServer} />

        <SectionSeparator />

        <SectionHeader>Add topic</SectionHeader>
        <AddTopicForm server={currentServer} />

        <Show when={currentServer()?.topics?.length > 0}>
          <SectionSeparator />
          <SectionHeader>Update topics</SectionHeader>
          <For each={currentServer()?.topics} key={(topic) => topic.id}>
            {(topic) => <UpdateTopicForm topic={topic} />}
          </For>
        </Show>
        <SectionSeparator />
        <DeleteServerForm server={currentServer()} />
      </div>
    </>
  );
}

function Error(props) {
  return (
    <p class="bg-red-200 text-red-700 rounded-lg mt-4 px-4">{props.children}</p>
  );
}

function SectionHeader(props) {
  return (
    <h2 class="text-lg tracking-wider font-semibold mb-2">{props.children}</h2>
  );
}

function SectionSeparator() {
  return <hr class="my-4" />;
}

function UpdateServerForm(props) {
  const serverName = () => props.server()?.name;
  const serverID = () => props.server()?.id;

  const submission = useSubmission(updateServerAction);

  const error = () => {
    const result = submission.result;
    if (result && !result.ok) {
      return result.error;
    }
  };

  return (
    <>
      <form
        action={updateServerAction.with(serverID(), serverName())}
        class="flex gap-2"
        method="post"
      >
        <input
          class="px-2 border-2 neo-shadow rounded-lg focus:outline-none focus:border-sky-500 w-full"
          name="name"
          value={serverName()}
          type="text"
          required
        />
        <Button type="submit" variant="sky" disabled={submission.pending}>
          <Check class="size-5" />
        </Button>
      </form>
      <Error>{error()}</Error>
    </>
  );
}

function AddTopicForm(props) {
  const serverID = () => props.server()?.id;

  const submission = useSubmission(addTopicAction);

  const error = () => {
    const result = submission.result;
    if (result && !result.ok) {
      return result.error;
    }
  };

  return (
    <>
      <form
        action={addTopicAction.with(serverID())}
        class="flex gap-2"
        method="post"
      >
        <select
          class="rounded-lg neo-shadow border-2 px-2 text-center tracking-wider"
          name="type"
        >
          <option value="text">Text</option>
          <option value="voice">Voice</option>
        </select>
        <input
          class="px-2 border-2 neo-shadow rounded-lg focus:outline-none focus:border-sky-500 w-full"
          name="name"
          type="text"
          required
        />
        <Button type="submit" variant="sky" disabled={submission.pending}>
          <Check class="size-5" />
        </Button>
      </form>
      <Error>{error()}</Error>
    </>
  );
}

function UpdateTopicForm(props) {
  const update = useAction(updateTopicAction);
  const del = useAction(deleteTopicAction);

  const updateSubmission = useSubmission(updateTopicAction);
  const deleteSubmission = useSubmission(deleteTopicAction);

  const handleUpdate = async (e) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    await update(props.topic.id, formData);
  };

  const updateError = () => {
    const result = updateSubmission.result;
    if (
      updateSubmission.input?.[0] === props.topic.id &&
      result &&
      !result.ok
    ) {
      return result.error;
    }
  };

  const deleteError = () => {
    const result = deleteSubmission.result;
    if (
      deleteSubmission.input?.[0] === props.topic.id &&
      result &&
      !result.ok
    ) {
      return result.error;
    }
  };

  return (
    <>
      <form onSubmit={handleUpdate} class="flex gap-2 mb-2">
        <div
          class="rounded-lg neo-shadow border-2 px-2 text-center tracking-wider w-16 flex items-center justify-center"
          classList={{
            "bg-sky-200": props.topic.type === "text",
            "bg-amber-200": props.topic.type === "voice",
          }}
        >
          {props.topic.type}
        </div>

        <input
          class="px-2 border-2 neo-shadow rounded-lg focus:outline-none focus:border-sky-500 flex-1"
          type="text"
          name="name"
          value={props.topic.name}
          required
        />

        <Button type="submit" variant="sky" disabled={updateSubmission.pending}>
          <Check class="size-5" />
        </Button>

        <Button
          type="button"
          variant="rose"
          disabled={deleteSubmission.pending}
          onClick={() => del(props.topic.id)}
        >
          <Trash2 class="size-5" />
        </Button>
      </form>

      <Error>{updateError()}</Error>
      <Error>{deleteError()}</Error>
    </>
  );
}

function DeleteServerForm(props) {
  const submission = useSubmission(deleteServerAction);

  const deleteError = () => {
    const result = submission.result;
    if (result && !result.ok) {
      return result.error;
    }
  };

  return (
    <form action={deleteServerAction.with(props.server?.id)} method="post">
      <Button type="submit" variant="rose" disabled={submission.pending}>
        Delete Server
      </Button>
      <Error>{deleteError()}</Error>
    </form>
  );
}
