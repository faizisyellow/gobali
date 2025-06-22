import { Outlet } from "@tanstack/react-router"

export default function AdminLayout() {
  return (
    <div>
      <h2>Admin Layout</h2>
      <Outlet />
    </div>
  )
}
