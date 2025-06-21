import { Burger, Button, Container, Group } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import classes from "./HeaderSimple.module.css";
import Logo from "../logo/Logo";

export function Header() {
  const [opened, { toggle }] = useDisclosure(false);

  return (
    <header className={classes.header}>
      <Container size="xl" className={classes.inner}>
        <Logo />

        <Burger opened={opened} onClick={toggle} hiddenFrom="xs" size="sm" />

        <Group visibleFrom="sm">
          <Button variant="default" radius="xs">
            Log in
          </Button>
          <Button radius="xs">Sign up</Button>
        </Group>
      </Container>
    </header>
  );
}
