import { Check, Trash2 } from "lucide-solid";
import Button from "~/components/Button";
import Header from "~/components/Header";
import HeaderItem from "~/components/HeaderItem";
import { useNavigate, useSubmission } from "@solidjs/router";
import { createServerAction } from "~/lib/create_server";

export default function Create() {
  return (
    <>
      <Header icon="settings">
        <HeaderItem>Create Your Server!</HeaderItem>
      </Header>
      <div class="m-auto max-w-xl min-w-sm w-full border-3 neo-shadow-lg rounded-xl p-4">
        <SectionHeader>Name</SectionHeader>
        <CreateServerForm />
      </div>
    </>
  );
}

function SectionHeader(props) {
  return (
    <h2 class="text-lg tracking-wider font-semibold mb-2">{props.children}</h2>
  );
}

function CreateServerForm() {
  const submission = useSubmission(createServerAction);

  const error = () => {
    const result = submission.result;
    if (result && !result.ok) {
      return result.error;
    }
  };

  return (
    <>
      <form
        action={createServerAction}
        class="flex flex-col gap-4 mt-4"
        method="post"
      >
        <input
          class="px-2 border-2 neo-shadow rounded-lg focus:outline-none focus:border-sky-500 flex-1"
          type="text"
          name="name"
          required
        />
        <Button
          type="submit"
          class="block ml-auto"
          variant="sky"
          disabled={submission.pending}
        >
          Create
        </Button>
      </form>
      <Error>{error()}</Error>
    </>
  );
}

function Error(props) {
  return (
    <p class="bg-red-200 text-red-700 rounded-lg mt-4 px-4">{props.children}</p>
  );
}
