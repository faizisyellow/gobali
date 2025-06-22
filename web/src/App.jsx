import { AuthProvider, useAuth } from "./context/auth/auth.jsx";
import { router } from "./router/routes.jsx";
import { RouterProvider } from "@tanstack/react-router";

function App() {
  return (
      <AuthProvider>
        <InnerApp />
      </AuthProvider>
  );
}

function InnerApp() {
  const auth = useAuth();
  return <RouterProvider router={router} context={{ auth }} />;
}

export default App;
