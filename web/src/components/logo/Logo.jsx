import { styled } from "@mui/material/styles";
import Typography from "@mui/material/Typography";

const StyledLogo = styled(Typography)`
  @font-face {
    font-weight: normal;
    font-style: normal;
  }

  font-style: italic;
  color: #5DC9E2;
`;

export default function Logo() {
  return <StyledLogo variant="h4" fontWeight={600}>Gobali</StyledLogo>
}
