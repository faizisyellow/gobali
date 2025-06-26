import {
  Card,
  CardMedia,
  CardContent,
  Typography,
  Grid,
  Button,
  Box,
  Badge,
  Tooltip,
  Chip,
} from "@mui/material";
import LocationOnIcon from "@mui/icons-material/LocationOn";
import GroupIcon from "@mui/icons-material/Group";
import BathtubIcon from "@mui/icons-material/Bathtub";
import HotelIcon from "@mui/icons-material/Hotel";

export default function VillaCard({ content }) {
  return (
    <Card
      sx={{ position: "relative", borderRadius: 0 }}
      elevation={0}
      variant="outlined"
    >
      {/* Image with Location Chip */}
      <Box position="relative">
        <CardMedia
          component="img"
          height="200"
          image={
            import.meta.env.VITE_BASE_URL_DEV +
            "/files/villas/" +
            content?.image_urls[0]
          }
          alt="Villa Image"
        />
        <Box position="absolute" bottom={8} left={16}>
          <Chip
            icon={<LocationOnIcon fontSize="small" />}
            label={content?.location?.area}
            size="small"
            color="success"
          />
        </Box>
      </Box>

      <CardContent>
        {/* Title */}
        <Typography variant="h6" gutterBottom>
          {content?.name}
        </Typography>

        {/* Category Badge */}
        <Box mb={2} mx={4}>
          <Tooltip title="Villa Category">
            <Badge badgeContent={content?.category.name} color="primary" />
          </Tooltip>
        </Box>

        {/* Stats as text with icons */}
        <Box display="flex" gap={2} alignItems="center" mb={2} flexWrap="wrap">
          <Box display="flex" alignItems="center" gap={0.5}>
            <GroupIcon fontSize="small" />
            <Typography variant="body2" color="text.secondary">
              {content?.min_guest} Guests
            </Typography>
          </Box>
          <Box display="flex" alignItems="center" gap={0.5}>
            <BathtubIcon fontSize="small" />
            <Typography variant="body2" color="text.secondary">
              {content?.baths} Baths
            </Typography>
          </Box>
          <Box display="flex" alignItems="center" gap={0.5}>
            <HotelIcon fontSize="small" />
            <Typography variant="body2" color="text.secondary">
              {content?.bedrooms} Bedrooms
            </Typography>
          </Box>
        </Box>

        {/* Price and Button */}
        <Grid container justifyContent="space-between" alignItems="center">
          <Typography variant="subtitle1">${content?.price} / Night</Typography>
          <Button variant="outlined">Show Details</Button>
        </Grid>
      </CardContent>
    </Card>
  );
}
