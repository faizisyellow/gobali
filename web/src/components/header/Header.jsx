import Logo from "../logo/Logo";
import { useNavigate } from "@tanstack/react-router";
import {
  AppBar,
  Box,
  Toolbar,
  Button,
  CssBaseline,
  Typography,
} from "@mui/material";
import { styled, useMediaQuery, useTheme } from "@mui/system";

const SlideContainer = styled(Box)(({ theme }) => ({
  width: "100%",
  height: "100vh",
  backgroundSize: "cover",
  backgroundPosition: "center",
  backgroundImage: `url(
  https://images.unsplash.com/photo-1688653802629-5360086bf632?q=80&w=1332&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D
)`,
  display: "flex",
  flexDirection: "column",
  justifyContent: "center",
  alignItems: "center",
  position: "relative",
  color: "#fff",
  textAlign: "center",
  padding: theme.spacing(4),
  "&::before": {
    content: '""',
    position: "absolute",
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    background: "linear-gradient(rgba(0,0,0,0.3), rgba(0,0,0,0.6))",
    zIndex: 1,
  },
}));

const TextContent = styled(Box)(({ theme }) => ({
  position: "relative",
  zIndex: 2,
  maxWidth: "600px",
}));

export function Header() {
  const navigate = useNavigate();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("sm"));

  return (
    <>
      <CssBaseline />
      <AppBar position="fixed" color="default" elevation={0}>
        <Toolbar>
          <Box sx={{ flexGrow: 1 }}>
            <Logo />
          </Box>
          <Button
            variant="outlined"
            sx={{ mr: 2 }}
            onClick={() => navigate({ to: "/login" })}
          >
            Log in
          </Button>
          <Button
            variant="contained"
            onClick={() => navigate({ to: "/signup" })}
          >
            Sign up
          </Button>
        </Toolbar>
      </AppBar>

      <Toolbar />

      <SlideContainer>
        <TextContent>
          <Typography
            variant={isMobile ? "h4" : "h3"}
            sx={{
              fontWeight: "bold",
              textShadow: "2px 2px 4px rgba(0,0,0,0.3)",
            }}
            gutterBottom
          >
            Discover Your Perfect Bali Villa with{" "}
            <span style={{ fontStyle: "italic", color: "#5DC9E2" }}>
              Gobali
            </span>
          </Typography>
          <Typography
            variant="body1"
            sx={{
              fontSize: isMobile ? "1rem" : "1.2rem",
              textShadow: "1px 1px 2px rgba(0,0,0,0.3)",
            }}
          >
            Experience luxury and tranquility with{" "}
            <span style={{ fontStyle: "italic", color: "#5DC9E2" }}>
              Gobali
            </span>{" "}
            — your trusted platform for booking the most stunning villas across
            Bali. Whether you're planning a romantic getaway, a family vacation,
            or a solo retreat,{" "}
            <span style={{ fontStyle: "italic", color: "#5DC9E2" }}>
              Gobali
            </span>{" "}
            connects you to curated villas that combine comfort, privacy, and
            authentic Balinese charm.
          </Typography>
        </TextContent>
      </SlideContainer>
    </>
  );
}
