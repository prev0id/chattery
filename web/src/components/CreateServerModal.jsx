import { createSignal } from "solid-js";
import Modal from "./Modal";
import FormTextInput from "./FormTextInput";
import Button from "./Button";
import { createServer } from "../stores/app";

export default function CreateServerModal(props) {
  const [name, setName] = createSignal("");
  const [isLoading, setIsLoading] = createSignal(false);

  const handleSubmit = async (event) => {
    event.preventDefault();
    setIsLoading(true);

    const success = await createServer(name());
    if (success) {
      setName("");
      document.getElementById(props.id).hidePopover();
    }

    setIsLoading(false);
  };

  return (
    <Modal id={props.id} name="Create Server">
      <form onSubmit={handleSubmit} class="flex flex-col gap-4 mt-4">
        <FormTextInput
          label="Server Name"
          value={name}
          onInput={(e) => setName(e.currentTarget.value)}
          required
        />
        <Button type="submit" variant="sky" disabled={isLoading()}>
          Create
        </Button>
      </form>
    </Modal>
  );
}
