export default function LoginPage() {
  function handleLogin(role) {
    localStorage.setItem("user", JSON.stringify({ role }))
    location.reload()
  }

  return (
    <div>
      <h1>Login Page</h1>
      <button onClick={() => handleLogin("user")}>Login as User</button>
      <button onClick={() => handleLogin("admin")}>Login as Admin</button>
    </div>
  )
}