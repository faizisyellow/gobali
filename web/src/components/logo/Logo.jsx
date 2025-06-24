import { styled } from "@mui/material/styles";
import Typography from "@mui/material/Typography";

const StyledLogo = styled(Typography)`
  @font-face {
    font-family: 'Mont';
    src: url('../../../assets/Mont.otf') format('truetype');
    font-weight: normal;
    font-style: normal;
  }

  font-family: 'Mont', sans-serif;
  font-style: italic;
  color: #5DC9E2;
`;

export default function Logo() {
  return <StyledLogo variant="h4" fontWeight={600}>Gobali</StyledLogo>
}
