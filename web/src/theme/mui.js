import { createTheme } from "@mui/material/styles";

const theme = createTheme({
  typography: {
    fontFamily: '"Poppins", sans-serif',
  },
  components: {
    MuiOutlinedInput: {
      styleOverrides: {
        root: () => ({
          borderRadius: 0,
          "& fieldset": {
            borderColor: "#dee2e6",
            borderWidth: 1.5,
          },
          "&:hover fieldset": {
            borderColor: "#212529",
          },
          "&.Mui-focused fieldset": {
            borderColor: "#212529",
            boxShadow: "none",
          },
        }),
      },
    },
    MuiFormControl: {
      styleOverrides: {
        root: ({ theme }) => ({
          marginBottom: theme.spacing(2),
        }),
      },
    },
  },
});

export default theme;
