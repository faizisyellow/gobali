import { Typography } from "@mui/material";

function Section({ title }) {
  return (
    <Typography variant="h4" noWrap align="left" sx={{ fontWeight: 500,marginBottom:5 }}>
      {title}
    </Typography>
  );
}

export default Section
