import { createUniqueId, Show } from "solid-js";

export default function TextField(props) {
  const id = props.id ?? createUniqueId();

  return (
    <>
      <Show when={props.label}>
        <label class="block font-semibold tracking-wider" for={id}>
          {props.label}
        </label>
      </Show>
      <input
        class="px-2 border-2 neo-shadow rounded-lg focus:outline-none focus:border-sky-500 w-full"
        id={id}
        type={props.type ?? "text"}
        value={props.value()}
        onInput={props.onInput}
        required={props.required}
        autocomplete={props.autocomplete}
      />
    </>
  );
}
