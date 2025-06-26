import {
  Box,
  Chip,
  FormControl,
  Grid,
  InputLabel,
  MenuItem,
  OutlinedInput,
  Select,
  TextField,
  Typography,
} from "@mui/material";
import VillaCard from "../villa-card/VillaCard";
import { useEffect, useState } from "react";
import { axiosQueryWithAuth } from "../../services/axios/auth/auth";

const locations = ["Canggu", "Ubud", "Seminyak", "Nusa Penida"];
const categories = ["Deluxe", "Luxury", "Regular"];

export default function Villas() {
  const [category, setCategory] = useState("");
  const [location, setLocation] = useState("");
  const [minGuest, setMinGuest] = useState("");
  const [bedrooms, setBedrooms] = useState("");
  const [villas, setVillas] = useState([]);

  useEffect(() => {
    async function fetch() {
      try {
        const query = {
          location,
          category,
          minGuest,
          bedrooms
        };

        const result = await axiosQueryWithAuth.GetAllVillas(query);
        const data = result?.data?.data;
        if (Array.isArray(data)) {
          setVillas(data);
        }
      } catch (error) {
        console.log(error);
      }
    }

    fetch();
  }, [location, category,minGuest,bedrooms]);

  return (
    <>
      <Typography variant="h5" mb={2}>
        Filters
      </Typography>
      <Grid container spacing={2} mb={2}>
        <Grid size={{ xs: 6, sm: 3 }}>
          <FormControl fullWidth>
            <InputLabel id="location-label">Location</InputLabel>
            <Select
              labelId="location-label"
              value={location}
              onChange={(e) => setLocation(e.target.value)}
              input={<OutlinedInput label="Location" />}
              renderValue={() =>
                location ? (
                  <Box display="flex" alignItems="center">
                    <Box mr={1}>
                      <Chip
                        label={location}
                        onDelete={() => {
                          setLocation("");
                          const selectEl = document.activeElement;
                          if (selectEl) selectEl.blur();
                        }}
                        onMouseDown={(event) => {
                          event.stopPropagation();
                        }}
                      />
                    </Box>
                  </Box>
                ) : (
                  <Typography color="text.secondary">
                    Select Location
                  </Typography>
                )
              }
            >
              {locations.map((loc) => (
                <MenuItem key={loc} value={loc}>
                  {loc}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Grid>

        <Grid size={{ xs: 6, sm: 3 }}>
          <FormControl fullWidth>
            <InputLabel id="category-label">Category</InputLabel>
            <Select
              labelId="category-label"
              value={category}
              onChange={(e) => setCategory(e.target.value)}
              input={<OutlinedInput label="Category" />}
              renderValue={() =>
                category ? (
                  <Box display="flex" alignItems="center">
                    <Box mr={1}>
                      <Chip
                        label={category}
                        onDelete={() => {
                          setCategory("");
                          const selectEl = document.activeElement;
                          if (selectEl) selectEl.blur();
                        }}
                        onMouseDown={(event) => {
                          event.stopPropagation();
                        }}
                      />
                    </Box>
                  </Box>
                ) : (
                  <Typography color="text.secondary">
                    Select Category
                  </Typography>
                )
              }
            >
              {categories.map((cat) => (
                <MenuItem key={cat} value={cat}>
                  {cat}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Grid>

        <Grid size={{ xs: 6, sm: 3 }}>
          <TextField
            fullWidth
            name="min_guest"
            type="number"
            label="Min Guest"
            value={minGuest}
            onChange={(e) => {
              setMinGuest(e.target.value);
            }}
          />
        </Grid>

        <Grid size={{ xs: 6, sm: 3 }}>
          <TextField
            fullWidth
            name="bedrooms"
            type="number"
            label="Bedrooms"
            value={bedrooms}
            onChange={(e) => {
              setBedrooms(e.target.value);
            }}
          />
        </Grid>
      </Grid>

      <Grid container spacing={4}>
        {villas.map((villa, index) => (
          <Grid size={{ sm: 12, md: 4 }} key={index}>
            <VillaCard content={villa} />
          </Grid>
        ))}
      </Grid>
    </>
  );
}
