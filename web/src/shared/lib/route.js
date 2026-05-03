export function parseRouteId(value) {
  const id = Number(value);
  return Number.isInteger(id) && id > 0 ? id : 0;
}
