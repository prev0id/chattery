export default function Button(props) {
  const {
    children,
    variant = "sky",
    sideways = false,
    class: extraClass,
    ...rest
  } = props;

  const base =
    "neo-shadow border-2 rounded-lg px-2 text-lg text-center font-semibold hover:scale-105 transition-all duration-300 ease-in-out tracking-widest focus:outline-none";

  const colorMap = {
    amber: "bg-amber-200 hover:bg-amber-500 focus:border-amber-500",
    sky: "bg-sky-200 hover:bg-sky-500 focus:border-sky-500",
    emerald: "bg-emerald-200 hover:bg-emerald-500 focus:border-emerald-500",
  };

  return (
    <button
      class={`${base} ${colorMap[props.variant ?? "sky"]} ${props.class ?? ""}`}
      classList={{
        "[writing-mode:sideways-lr]": props.sideways ?? false,
      }}
      {...rest}
    >
      {children}
    </button>
  );
}
