import { useState } from "react";
import {
  TextInput,
  PasswordInput,
  Button,
  Paper,
  Title,
  Container,
  Stack,
  Text,
  Box,
  Group,
  Anchor,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { IconMail, IconLock } from "@tabler/icons-react";
import { axiosQueryPublic } from "../services/axios/public/public";
import { useAuth } from "../context/auth/auth";
import { useNavigate } from "@tanstack/react-router";
import styles from "./../module/loginPage.module.css";

export default function LoginPage() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const { setCredentials } = useAuth();
  const navigate = useNavigate();

  const form = useForm({
    initialValues: {
      email: "",
      password: "",
    },
    validate: {
      email: (value) => (/^\S+@\S+$/.test(value) ? null : "Invalid email"),
      password: (value) =>
        value.length < 6 ? "Password must be at least 6 characters" : null,
    },
  });

  async function handleLogin(values) {
    setLoading(true);
    const payload = { email: values.email, password: values.password };

    try {
      const response = await axiosQueryPublic.login(payload);
      setCredentials(response?.data?.data);

      navigate({ to: "/browse" });
    } catch (err) {
      // TODO: ADD TOAST ERROR
      setError("Email or password may be incorrect");
      console.log(err);
    } finally {
      setLoading(false);
    }
  }

  return (
    <>
      <Box className={styles.container}>
        <Container size="xs" className={styles.formContainer}>
          <Paper shadow="none" className={styles.paper}>
            <Stack gap={32}>
              <Box className={styles.header}>
                <Title order={1} size="h2" className={styles.title}>
                  Welcome back
                </Title>
                <Text c="dimmed" size="sm" className={styles.subtitle}>
                  Sign in to your account to continue
                </Text>
              </Box>

              <Stack
                gap={24}
                component="form"
                onSubmit={form.onSubmit(handleLogin)}
              >
                <TextInput
                  label="Email address"
                  placeholder="Enter your email"
                  type="email"
                  required
                  leftSection={<IconMail size={18} />}
                  className={styles.input}
                  classNames={{
                    input: styles.inputField,
                    label: styles.inputLabel,
                  }}
                  {...form.getInputProps("email")}
                />

                <PasswordInput
                  label="Password"
                  placeholder="Enter your password"
                  required
                  leftSection={<IconLock size={18} />}
                  className={styles.input}
                  classNames={{
                    input: styles.inputField,
                    label: styles.inputLabel,
                    innerInput: styles.passwordInput,
                  }}
                  {...form.getInputProps("password")}
                />

                <Group justify="flex-end" className={styles.forgotPassword}>
                  <Anchor
                    size="sm"
                    c="dimmed"
                    className={styles.link}
                    onClick={() => console.log("Forgot password clicked")}
                  >
                    Forgot password?
                  </Anchor>
                </Group>

                <Button
                  type="submit"
                  fullWidth
                  loading={loading}
                  className={styles.submitButton}
                >
                  Sign in
                </Button>
              </Stack>

              <Box className={styles.signupSection}>
                <Text size="sm" c="dimmed" className={styles.signupText}>
                  Don't have an account?{" "}
                  <Anchor
                    size="sm"
                    className={styles.signupLink}
                    onClick={() => console.log("Sign up clicked")}
                  >
                    Sign up
                  </Anchor>
                </Text>
              </Box>
            </Stack>
          </Paper>
        </Container>
      </Box>
    </>
  );
}
