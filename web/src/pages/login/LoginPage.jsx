import { useState } from "react";
import { useFormik } from "formik";
import * as yup from "yup";
import {
  Button,
  TextField,
  Container,
  Typography,
  Box,
  Paper,
  Link,
  Snackbar,
  Alert,
} from "@mui/material";
import { styled } from "@mui/material/styles";
import { axiosQueryPublic } from "../../services/axios/public/public";
import { useNavigate } from "@tanstack/react-router";
import { useAuth } from "../../context/auth/auth";

const validationSchema = yup.object({
  email: yup
    .string("Enter your email")
    .email("Enter a valid email")
    .required("Email is required"),
  password: yup
    .string("Enter your password")
    .min(6, "Password should be of minimum 6 characters length")
    .required("Password is required"),
});

export default function LoginPage() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [snackbar, setSnackbar] = useState({
    open: false,
    vertical: "bottom",
    horizontal: "right",
  });
  const { vertical, horizontal, open} = snackbar;

  const handleClose = (event, reason) => {
    if (reason === "clickaway") {
      return;
    }

    setSnackbar({ ...snackbar, open: false });
  };

  const { setCredentials } = useAuth();
  const navigate = useNavigate();

  const formik = useFormik({
    initialValues: {
      email: "",
      password: "",
    },
    validationSchema: validationSchema,
    onSubmit: async (values) => {
      setLoading(true);
      const payload = { email: values.email, password: values.password };

      try {
        const response = await axiosQueryPublic.login(payload);
        setCredentials(response?.data?.data);
        navigate({ to: "/" });
      } catch (err) {
        setError(err.toString());
        setSnackbar({ ...snackbar, open: true });
      } finally {
        setLoading(false);
      }
    },
  });

  return (
    <>
      <LoginContainer>
        <Container maxWidth="xs" sx={{ width: "100%", maxWidth: "400px" }}>
          <LoginPaper elevation={0}>
            <Box textAlign="center">
              <LoginTitle variant="h4" component="h1">
                Welcome back
              </LoginTitle>
              <LoginSubtitle variant="body2">
                Sign in to your account to continue
              </LoginSubtitle>
            </Box>

            <Box component="form" onSubmit={formik.handleSubmit} noValidate>
              <TextField
                fullWidth
                id="email"
                name="email"
                label="Email address"
                placeholder="Enter your email"
                variant="outlined"
                value={formik.values.email}
                onChange={formik.handleChange}
                onBlur={formik.handleBlur}
                error={formik.touched.email && Boolean(formik.errors.email)}
                helperText={formik.touched.email && formik.errors.email}
              />

              <TextField
                fullWidth
                id="password"
                name="password"
                label="Password"
                placeholder="Enter your password"
                type="password"
                variant="outlined"
                value={formik.values.password}
                onChange={formik.handleChange}
                onBlur={formik.handleBlur}
                error={
                  formik.touched.password && Boolean(formik.errors.password)
                }
                helperText={formik.touched.password && formik.errors.password}
                sx={{ mb: 1 }}
              />

              <Box display="flex" justifyContent="flex-end" mb={2}>
                <ForgotPasswordLink
                  component="button"
                  type="button"
                  variant="body2"
                  onClick={() => console.log("Forgot password clicked")}
                >
                  Forgot password?
                </ForgotPasswordLink>
              </Box>

              <LoginButton
                type="submit"
                fullWidth
                variant="contained"
                disabled={loading}
              >
                {loading ? "Signing in..." : "Sign in"}
              </LoginButton>

              <Box textAlign="center">
                <Typography variant="body2" color="text.secondary">
                  Don't have an account?{" "}
                  <SignupLink
                    component="button"
                    type="button"
                    variant="body2"
                    onClick={() => console.log("Sign up clicked")}
                  >
                    Sign up
                  </SignupLink>
                </Typography>
              </Box>
            </Box>
          </LoginPaper>
        </Container>
      </LoginContainer>

      {error && (
        <Snackbar
          anchorOrigin={{ vertical, horizontal }}
          open={open}
          onClose={handleClose}
          autoHideDuration={2000}
        >
          <Alert
            onClose={handleClose}
            severity="error"
            variant="filled"
            sx={{ width: "100%" }}
          >
            {error.includes("Network")
              ? "Server Error Try Another Time"
              : "Email or Password maybe incorrect"}
          </Alert>
        </Snackbar>
      )}
    </>
  );
}

const LoginContainer = styled(Box)(({ theme }) => ({
  minHeight: "100vh",
  backgroundColor: "#fafafa",
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
  padding: theme.spacing(1),
  [theme.breakpoints.down("md")]: {
    alignItems: "flex-start",
    paddingTop: theme.spacing(4),
  },
  [theme.breakpoints.down("sm")]: {
    padding: theme.spacing(0.5),
    paddingTop: theme.spacing(2),
  },
}));

const LoginPaper = styled(Paper)(({ theme }) => ({
  padding: theme.spacing(5),
  border: "1px solid #e9ecef",
  borderRadius: 0,
  backgroundColor: "white",
  boxShadow: "none",
  [theme.breakpoints.down("md")]: {
    padding: theme.spacing(3),
    border: "none",
    backgroundColor: "transparent",
  },
  [theme.breakpoints.down("sm")]: {
    padding: theme.spacing(2.5, 2),
  },
}));

const LoginTitle = styled(Typography)(({ theme }) => ({
  fontWeight: 600,
  color: "#212529",
  marginBottom: theme.spacing(1),
  letterSpacing: "-0.02em",
  fontSize: "1.75rem",
  [theme.breakpoints.down("sm")]: {
    fontSize: "1.375rem",
  },
}));

const LoginSubtitle = styled(Typography)(({ theme }) => ({
  color: theme.palette.text.secondary,
  marginBottom: theme.spacing(4),
  fontWeight: 400,
}));


const LoginButton = styled(Button)(({ theme }) => ({
  marginTop: theme.spacing(2),
  marginBottom: theme.spacing(3),
  borderRadius: 0,
  textTransform: "none",
  fontWeight: 500,
  fontSize: "14px",
  height: "48px",
  backgroundColor: "#212529",
  color: "white",
  "&:hover": {
    backgroundColor: "#343a40",
  },
  "&:disabled": {
    backgroundColor: "#6c757d",
  },
  [theme.breakpoints.down("sm")]: {
    fontSize: "16px",
    height: "52px",
  },
}));

const ForgotPasswordLink = styled(Link)(({ theme }) => ({
  color: theme.palette.text.secondary,
  textDecoration: "none",
  fontWeight: 400,
  fontSize: "0.875rem",
  cursor: "pointer",
  "&:hover": {
    textDecoration: "underline",
  },
}));

const SignupLink = styled(Link)(() => ({
  color: "#212529",
  textDecoration: "none",
  fontWeight: 500,
  cursor: "pointer",
  "&:hover": {
    textDecoration: "underline",
  },
}));
