import { useState } from "react";
import { axiosQueryPublic } from "../services/axios/public/public";
import { useAuth } from "../context/auth/auth";
import { useNavigate } from "@tanstack/react-router";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const { setCredentials } = useAuth();
  const navigate = useNavigate();

  async function handleLogin(e) {
    e.preventDefault();
    const payload = { email, password };

    try {
      const response = await axiosQueryPublic.login(payload);
      setCredentials(response?.data?.data);

      navigate({ to: "/" });
    } catch (err) {
      console.log(err);
    }
  }

  return (
    <div style={{ padding: "2rem" }}>
      <h1 className="text-xl font-bold mb-4">Login Page</h1>
      <form onSubmit={handleLogin} className="flex flex-col gap-4 max-w-sm">
        <label className="flex flex-col">
          Email:
          <input
            type="email"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="border p-2"
          />
        </label>

        <label className="flex flex-col">
          Password:
          <input
            type="password"
            required
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="border p-2"
          />
        </label>

        <button type="submit" className="bg-blue-600 text-white py-2 rounded">
          Login
        </button>
      </form>
    </div>
  );
}
