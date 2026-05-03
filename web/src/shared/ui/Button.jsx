import { splitProps } from "solid-js";
import { cn } from "~/shared/lib/cn";

const buttonVariantClass = {
  amber:
    "bg-amber-200 disabled:bg-amber-100 hover:bg-amber-500 focus:border-amber-500 rounded-lg neo-shadow border-2",
  sky: "bg-sky-200 disabled:bg-sky-100 hover:bg-sky-500 focus:border-sky-500 rounded-lg neo-shadow border-2",
  emerald:
    "bg-emerald-200 disabled:bg-emerald-100 hover:bg-emerald-500 focus:border-emerald-500 rounded-lg neo-shadow border-2",
  rose: "bg-rose-200 disabled:bg-rose-100 hover:bg-rose-500 focus:border-rose-500 rounded-lg neo-shadow border-2",
};

export default function Button(props) {
  const [local, rest] = splitProps(props, [
    "variant",
    "smallText",
    "sideways",
    "class",
  ]);

  const base =
    "px-2 text-center transition-all duration-300 ease-in-out focus:outline-none";

  return (
    <button
      {...rest}
      type={rest.type ?? "button"}
      class={cn(
        base,
        buttonVariantClass[local.variant ?? "sky"],
        !local.smallText && "text-lg font-semibold tracking-widest",
        local.sideways && "[writing-mode:sideways-lr]",
        local.class,
      )}
    >
      {props.children}
    </button>
  );
}
