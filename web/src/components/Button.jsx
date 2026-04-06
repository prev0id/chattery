export default function Button(props) {
  const base =
    "px-2 text-center hover:scale-105 transition-all duration-300 ease-in-out focus:outline-none";

  const colorMap = {
    amber:
      "bg-amber-200 disabled:bg-amber-100 hover:bg-amber-500 focus:border-amber-500 rounded-lg neo-shadow border-2",
    sky: "bg-sky-200 disabled:bg-sky-100 hover:bg-sky-500 focus:border-sky-500 rounded-lg neo-shadow border-2",
    emerald:
      "bg-emerald-200 disabled:bg-emerald-100 hover:bg-emerald-500 focus:border-emerald-500 rounded-lg neo-shadow border-2",
    rose: "bg-rose-200 disabled:bg-rose-100 hover:bg-rose-500 focus:border-rose-500 rounded-lg neo-shadow border-2",
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
