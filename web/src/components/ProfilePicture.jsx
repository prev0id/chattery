export default function ProfilePicture(props) {
  const { src, class: extraClass, ...rest } = props;

  return (
    <img
      src={src}
      class={`border-2 neo-shadow rounded-lg ${extraClass || ""}`}
      alt="profile_picture"
      {...rest}
    />
  );
}
