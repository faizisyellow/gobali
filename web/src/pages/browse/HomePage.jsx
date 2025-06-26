import { Box, Grid } from "@mui/material";
import { Header } from "../../components/header/Header";
import VillaCard from "../../components/villa-card/VillaCard";
import Villas from "../../components/villas/Villas";

export default function HomePage() {
  return (
    <>
      <Header />
      <Box m={4}>
        <Villas />
      </Box>
    </>
  );
}
