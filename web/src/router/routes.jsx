import {
  Router,
  Outlet,
  createRootRoute,
  createRoute,
  redirect,
} from "@tanstack/react-router";
import HomePage from "../pages/HomePage";
import LoginPage from "../pages/LoginPage";
import UserDashboard from "../pages/UserDashboard";
import AdminDashboard from "../pages/AdminDashboard";
import UserLayout from "../layouts/UserLayout";
import AdminLayout from "../layouts/AdminLayout";
import { authLoader } from "../components/AuthLoader";

const rootRoute = createRootRoute({
  component: () => (
    <>
      <Outlet />
    </>
  ),
});

const publicRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  component: HomePage,
  beforeLoad: async () => {
    const user = await authLoader();

    if (user && user.role === "admin") {
      throw redirect({ to: "/admin" });
    }

    return {};
  },
});

// Login route - redirect if already authenticated
const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/login",
  component: LoginPage,
  beforeLoad: async () => {
    const user = await authLoader();
    if (user) {
      if (user.role === "admin") {
        throw redirect({ to: "/admin" });
      } else {
        throw redirect({ to: "/" });
      }
    }
    return {};
  },
});

// USER Routes
const userLayoutRoute = createRoute({
  getParentRoute: () => rootRoute,
  id: "user-layout",
  beforeLoad: async () => {
    const user = await authLoader();

    if (!user) {
      throw redirect({ to: "/login" });
    }

    if (user.role === "admin") {
      throw redirect({ to: "/admin" });
    }

    if (user.role !== "user") {
      throw redirect({ to: "/login" });
    }

    return { user };
  },
  component: UserLayout,
});

const userDashboardRoute = createRoute({
  getParentRoute: () => userLayoutRoute,
  path: "/profile",
  component: UserDashboard,
});

// ADMIN Routes
const adminLayoutRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/admin",
  beforeLoad: async () => {
    const user = await authLoader();
    if (!user) {
      throw redirect({ to: "/login" });
    }
    if (user.role !== "admin") {
      throw redirect({ to: "/" });
    }
    return { user };
  },
  component: AdminLayout,
});

const adminDashboardRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: "/",
  component: AdminDashboard,
});


const routeTree = rootRoute.addChildren([
  publicRoute,
  loginRoute,
  userLayoutRoute.addChildren([userDashboardRoute]),
  adminLayoutRoute.addChildren([adminDashboardRoute]),
]);

export const router = new Router({
  routeTree,
  defaultPreload: "intent",
});
