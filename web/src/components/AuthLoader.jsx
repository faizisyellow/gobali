export async function authLoader() {
  const user = JSON.parse(localStorage.getItem("user"));

  if (!user) {
    throw new Error("Unauthorized");
  }

  return user;
}
