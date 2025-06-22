import Logo from "../logo/Logo";
import { useNavigate } from "@tanstack/react-router";

export function Header() {
  const navigate = useNavigate();

  function goLogin() {
    navigate({ to: "/login" });
  }

  function goSignup() {
    navigate({ to: "/signup" });
  }

  return <p>header</p>;
}
