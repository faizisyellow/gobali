import { Box, Typography, Button, Paper } from '@mui/material';
import ErrorOutlineIcon from '@mui/icons-material/ErrorOutline';

const ErrorPage = ({ status = 500, message = 'Unexpected server error.', onRetry }) => {
  return (
    <Box
      sx={{
        minHeight: '60vh',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        p: 2,
        backgroundColor: '#f8f9fa',
      }}
    >
      <Paper
        elevation={3}
        sx={{
          maxWidth: 480,
          width: '100%',
          textAlign: 'center',
          p: 4,
          borderRadius: 4,
          backgroundColor: '#fff',
        }}
      >
        <ErrorOutlineIcon sx={{ fontSize: 80, color: 'error.main', mb: 1 }} />

        <Typography variant="h2" fontWeight={700} color="error.main" gutterBottom>
          {status}
        </Typography>

        <Typography variant="h6" color="text.secondary" gutterBottom>
          {message}
        </Typography>

        {onRetry && (
          <Button
            variant="contained"
            color="error"
            onClick={onRetry}
            sx={{ mt: 3, borderRadius: 2, px: 4 }}
          >
            Try Again
          </Button>
        )}
      </Paper>
    </Box>
  );
};

export default ErrorPage;
