import { Burger, Button, Container, Group } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import classes from "./HeaderSimple.module.css";
import Logo from "../logo/Logo";
import { useNavigate } from "@tanstack/react-router";

export function Header() {
  const [opened, { toggle }] = useDisclosure(false);
  const navigate = useNavigate();

  function goLogin() {
    navigate({ to: "/login" });
  }

  function goSignup() {
    navigate({ to: "/signup" });
  }

  return (
    <header className={classes.header}>
      <Container size="xl" className={classes.inner}>
        <Logo />

        <Burger opened={opened} onClick={toggle} hiddenFrom="xs" size="sm" />

        <Group visibleFrom="sm">
          <Button variant="default" radius="xs" onClick={goLogin}>
            Log in
          </Button>
          <Button radius="xs" onClick={goSignup}>
            Sign up
          </Button>
        </Group>
      </Container>
    </header>
  );
}
