export default function ProfilePicture(props) {
  return (
    <img
      {...props}
      class={`border-2 neo-shadow rounded-lg ${props.class ?? ""}`}
    />
  );
}
