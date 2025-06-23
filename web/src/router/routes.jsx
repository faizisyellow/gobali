import {
  Router,
  Outlet,
  createRoute,
  redirect,
  createRootRouteWithContext,
} from "@tanstack/react-router";
import HomePage from "../pages/browse/HomePage";
import LoginPage from "../pages/login/LoginPage";
import UserDashboard from "../pages/users/UserDashboard";
import VillaView from "../pages/admin/villa/View";
import VillaAdd from "../pages/admin/villa/Add";
import VillaUpdate from "../pages/admin/villa/Update";
import UserLayout from "../layouts/users/UserLayout";
import AdminLayout from "../layouts/admin/AdminLayout";

const rootRoute = createRootRouteWithContext()({
  component: () => <Outlet />,
});

const browseRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/browse",
  component: HomePage,
  beforeLoad: ({ context }) => {
    const { role, isLoggedIn } = context.auth;

    if (isLoggedIn && role == "admin") {
      throw redirect({ to: "/" });
    }

    return;
  },
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

const villaManagementViewRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: "/",
  component: VillaView,
});

const villaManagementNewRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: "/villas-management/new",
  component: VillaAdd,
});

const villaManagementUpdateRoute = createRoute({
  getParentRoute: () => adminLayoutRoute,
  path: "/villas-management/$id",
  component: VillaUpdate,
});

const routeTree = rootRoute.addChildren([
  browseRoute,
  loginRoute,
  userLayoutRoute.addChildren([userDashboardRoute]),
  adminLayoutRoute.addChildren([
    villaManagementViewRoute,
    villaManagementNewRoute,
    villaManagementUpdateRoute,
  ]),
]);

export const router = new Router({
  routeTree,
 });
