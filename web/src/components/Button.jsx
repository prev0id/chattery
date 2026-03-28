export default function Button(props) {
  const base =
    "neo-shadow border-2 rounded-lg px-2 text-center hover:scale-105 transition-all duration-300 ease-in-out focus:outline-none";

  const colorMap = {
    amber: "bg-amber-200 hover:bg-amber-500 focus:border-amber-500",
    sky: "bg-sky-200 hover:bg-sky-500 focus:border-sky-500",
    emerald: "bg-emerald-200 hover:bg-emerald-500 focus:border-emerald-500",
  };

  return (
    <button
      {...props}
      class={`${base} ${colorMap[props.variant ?? "sky"]} ${props.class ?? ""}`}
      classList={{
        "[writing-mode:sideways-lr]": props.sideways,
        "text-lg font-semibold tracking-widest": !props.smallText,
      }}
    >
      {props.children}
    </button>
  );
}
