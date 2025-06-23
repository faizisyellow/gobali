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
  Button,
} from "@mui/material";
import Section from "../../../components/section/Section";
import { useEffect, useState } from "react";
import { axiosQueryWithAuth } from "../../../services/axios/auth/auth";
import ErrorPage from "../../../components/error/Error";
import { useFormik } from "formik";
import * as Yup from "yup";
import useDebouncedFormikField from "../../../lib/hooks/DebounceFormikField";

const ITEM_HEIGHT = 48;
const ITEM_PADDING_TOP = 8;
const MenuProps = {
  PaperProps: {
    style: {
      maxHeight: ITEM_HEIGHT * 4.5 + ITEM_PADDING_TOP,
      width: 250,
    },
  },
};

const validationSchema = Yup.object({
  name: Yup.string().required("Name is required"),
  description: Yup.string().required("Description is required"),
  min_guest: Yup.number().required("Min guest is required"),
  bedrooms: Yup.number().required("Bedrooms is required"),
  baths: Yup.number().required("Baths is required"),
  price: Yup.number().required("Price is required"),
  category_id: Yup.number().required("Category is required"),
  location_id: Yup.number().required("Location is required"),
  amenity_id: Yup.array()
    .of(Yup.number())
    .required("Select at least one amenity"),
});

export default function VillaAdd() {
  const [datas, setDatas] = useState({
    locations: [],
    categories: [],
    amenities: [],
  });
  const [errorResponse, setErrorResponse] = useState({
    locations: null,
    categories: null,
    amenities: null,
  });

  const formik = useFormik({
    initialValues: {
      name: "",
      description: "",
      min_guest: "",
      bedrooms: "",
      baths: "",
      price: "",
      category_id: "",
      location_id: "",
      amenity_id: [],
    },
    validationSchema,
    onSubmit: (values) => {
      console.log("Submitted:", values);
    },
  });

  const nameFieldDebounce = useDebouncedFormikField(formik, "name");
  const descriptionFieldDebounce = useDebouncedFormikField(
    formik,
    "description"
  );

  useEffect(() => {
    async function fetchDatas() {
      const results = await Promise.allSettled([
        axiosQueryWithAuth.getAllLocations(),
        axiosQueryWithAuth.getAllCategories(),
        axiosQueryWithAuth.getAllAmenties(),
      ]);

      const [locationsResult, categoriesResult, amenitiesResult] = results;

      setDatas({
        locations:
          locationsResult.status === "fulfilled"
            ? locationsResult.value?.data?.data || []
            : [],
        categories:
          categoriesResult.status === "fulfilled"
            ? categoriesResult.value?.data?.data || []
            : [],
        amenities:
          amenitiesResult.status === "fulfilled"
            ? amenitiesResult.value?.data?.data || []
            : [],
      });

      setErrorResponse({
        locations:
          locationsResult.status === "rejected" ? locationsResult.reason : null,
        categories:
          categoriesResult.status === "rejected"
            ? categoriesResult.reason
            : null,
        amenities:
          amenitiesResult.status === "rejected" ? amenitiesResult.reason : null,
      });
    }

    fetchDatas();
  }, []);

  // If error checks
  if (
    errorResponse.amenities ||
    errorResponse.categories ||
    errorResponse.locations
  ) {
    // what's this ??
    const statuses = ["amenities", "categories", "locations"]
      .map((field) => errorResponse[field]?.status)
      .filter(Boolean);

    if (statuses.includes(500) || statuses.length === 0) {
      return (
        <ErrorPage
          status={500}
          message="Looks like our server there's an unexpected happening"
        />
      );
    } else {
      return (
        <ErrorPage
          status={400}
          message="Looks like our side there's an unexpected error"
        />
      );
    }
  }

  return (
    <>
      <Section title="Create New Villa" />
      <form onSubmit={formik.handleSubmit}>
        <Grid container spacing={2}>
          {/* Name */}
          <Grid size={12}>
            <TextField
              fullWidth
              label="Name"
              name="name"
              {...nameFieldDebounce}
              error={formik.touched.name && Boolean(formik.errors.name)}
              helperText={formik.touched.name && formik.errors.name}
            />
          </Grid>

          {/* Category */}
          <Grid size={{ xs: 12, lg: 6 }}>
            <FormControl fullWidth>
              <InputLabel id="category-label">Category</InputLabel>
              <Select
                name="category_id"
                labelId="category-label"
                value={formik.values.category_id}
                onChange={formik.handleChange}
                input={<OutlinedInput label="Category" />}
              >
                {datas.categories.map((cat) => (
                  <MenuItem value={cat.id} key={cat.id}>
                    {cat.name}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>

          {/* Location */}
          <Grid size={{ xs: 12, lg: 6 }}>
            <FormControl fullWidth>
              <InputLabel id="location-label">Location</InputLabel>
              <Select
                name="location_id"
                labelId="location-label"
                value={formik.values.location_id}
                onChange={formik.handleChange}
                input={<OutlinedInput label="Location" />}
              >
                {datas.locations.map((loc) => (
                  <MenuItem value={loc.id} key={loc.id}>
                    {loc.area}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>

          {/* Min Guest */}
          <Grid size={{ xs: 12, lg: 3 }}>
            <TextField
              fullWidth
              label="Min Guest"
              name="min_guest"
              type="number"
              value={formik.values.min_guest}
              onChange={formik.handleChange}
            />
          </Grid>

          {/* Bedrooms */}
          <Grid size={{ xs: 12, lg: 3 }}>
            <TextField
              fullWidth
              label="Bedrooms"
              name="bedrooms"
              type="number"
              value={formik.values.bedrooms}
              onChange={formik.handleChange}
            />
          </Grid>

          {/* Price */}
          <Grid size={{ xs: 12, lg: 3 }}>
            <TextField
              fullWidth
              label="Price"
              name="price"
              type="number"
              value={formik.values.price}
              onChange={formik.handleChange}
            />
          </Grid>

          {/* Baths */}
          <Grid size={{ xs: 12, lg: 3 }}>
            <TextField
              fullWidth
              label="Baths"
              name="baths"
              type="number"
              value={formik.values.baths}
              onChange={formik.handleChange}
            />
          </Grid>

          {/* Description */}
          <Grid size={{ xs: 12, lg: 6 }}>
            <TextField
              fullWidth
              multiline
              rows={4}
              label="Description"
              name="description"
              {...descriptionFieldDebounce}
            />
          </Grid>

          {/* Amenities */}
          <Grid size={{ xs: 12, lg: 6 }}>
            <FormControl fullWidth>
              <InputLabel id="amenities-label">Amenities</InputLabel>
              <Select
                multiple
                name="amenity_id"
                labelId="amenities-label"
                value={formik.values.amenity_id}
                onChange={(e) =>
                  formik.setFieldValue(
                    "amenity_id",
                    typeof e.target.value === "string"
                      ? e.target.value.split(",").map(Number)
                      : e.target.value.map(Number)
                  )
                }
                input={<OutlinedInput label="Amenities" />}
                renderValue={(selected) => (
                  <Box sx={{ display: "flex", flexWrap: "wrap", gap: 0.5 }}>
                    {selected.map((id) => {
                      const amenity = datas.amenities.find((a) => a.id === id);
                      return <Chip key={id} label={amenity?.name || id} />;
                    })}
                  </Box>
                )}
                MenuProps={MenuProps}
              >
                {datas.amenities.map((am) => (
                  <MenuItem key={am.id} value={am.id}>
                    {am.name}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>

          {/* Submit */}
          <Grid alignItems={"flex-end"}>
            <Button type="submit" variant="contained" color="primary">
              Submit Villa
            </Button>
          </Grid>
        </Grid>
      </form>
    </>
  );
}
