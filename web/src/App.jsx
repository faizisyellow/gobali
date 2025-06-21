import { AuthProvider, useAuth } from "./context/auth/auth.jsx";
import { router } from "./router/routes.jsx";
import { RouterProvider } from "@tanstack/react-router";
import { MantineProvider } from "@mantine/core";
import "@mantine/core/styles.css";

function App() {
  return (
    <MantineProvider>
      <AuthProvider>
        <InnerApp />
      </AuthProvider>
    </MantineProvider>
  );
}

function InnerApp() {
  const auth = useAuth();
  return <RouterProvider router={router} context={{ auth }} />;
}

export default App;
