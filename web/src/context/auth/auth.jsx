import { createContext, useContext, useEffect, useState } from "react";
import { decodeJWTClaims } from "../../lib/jwt";

const AuthContext = createContext();

export function AuthProvider({ children }) {
  const [token, setToken] = useState(null);
  const [role, setRole] = useState(null);

  const decodeToken = (jwt) => {
    try {
      const payload = decodeJWTClaims(jwt);

      return payload?.role || null;
    } catch (err) {
      console.error("Invalid token:", err);
      return null;
    }
  };

  useEffect(() => {
    const storedToken = localStorage.getItem("auth_token");

    // if there's a token set token and role.
    if (storedToken) {
      setToken(storedToken);
      setRole(decodeToken(storedToken));
    }
  }, []);

  const setCredentials = (newToken) => {
    localStorage.setItem("auth_token", newToken);
    setToken(newToken);
    setRole(decodeToken(newToken));
  };

  const logout = () => {
    localStorage.removeItem("auth_token");
    setToken(null);
    setRole(null);
  };

  return (
    <AuthContext.Provider
      value={{
        token,
        role,
        isLoggedIn: !!token,
        setCredentials,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}
