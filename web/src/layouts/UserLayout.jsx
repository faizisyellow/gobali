import { Outlet } from "@tanstack/react-router"

export default function UserLayout() {
  return (
    <div>
      <h2>User Layout</h2>
      <Outlet />
    </div>
  )
}