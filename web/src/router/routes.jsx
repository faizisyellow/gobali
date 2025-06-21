import {
  Router,
  Outlet,
  createRoute,
  redirect,
  createRootRouteWithContext,
} from "@tanstack/react-router";
import HomePage from "../pages/HomePage";
import LoginPage from "../pages/LoginPage";
import UserDashboard from "../pages/UserDashboard";
import AdminDashboard from "../pages/AdminDashboard";
import UserLayout from "../layouts/UserLayout";
import AdminLayout from "../layouts/AdminLayout";

const rootRoute = createRootRouteWithContext()({
  component: () => <Outlet />,
});

const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/browse",
  component: HomePage,
});

const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/login",
  component: LoginPage,
  beforeLoad: ({ context }) => {
    const { role, isLoggedIn } = context.auth;

    if (isLoggedIn && role == "user") {
      throw redirect({ to: "/" });
    }
    if (isLoggedIn && role == "admin") {
      throw redirect({ to: "/admin" });
    }

    return;
  },
});

// USER Routes
const userLayoutRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  beforeLoad: ({ context }) => {
    const { role, isLoggedIn } = context.auth;

    if (!isLoggedIn) {
      throw redirect({ to: "/login" });
    }

    if (role !== "user") {
      throw redirect({ to: role === "admin" ? "/admin" : "/browse" });
    }

    return;
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
  beforeLoad: ({ context }) => {
    const { role, isLoggedIn } = context.auth;

    if (!isLoggedIn) {
      throw redirect({ to: "/login" });
    }

    if (role !== "admin") {
      throw redirect({ to: "/browse" });
    }

    return;
  },
  component: AdminLayout,
});

const adminDashboardRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: "/",
  component: AdminDashboard,
});

const routeTree = rootRoute.addChildren([
  indexRoute,
  loginRoute,
  userLayoutRoute.addChildren([userDashboardRoute]),
  adminLayoutRoute.addChildren([adminDashboardRoute]),
]);

export const router = new Router({
  routeTree,
});
