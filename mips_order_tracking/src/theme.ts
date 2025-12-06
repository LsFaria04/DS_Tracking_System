import { createTheme } from '@mui/material/styles';

/**
 * This theme serves as a fallback and base configuration.
 * The host application's theme will be used if it's provided via ThemeProvider.
 * The theme automatically adapts to light/dark mode based on the host's palette.mode setting.
 * 
 * Since @emotion/react is shared, Material UI will use the host's Emotion instance,
 * ensuring consistent styling across all micro-frontends.
 */
export const theme = createTheme({
  palette: {
    mode: 'dark', // Host can override this
    primary: {
      main: '#1976d2',
      light: '#42a5f5',
      dark: '#1565c0',
    },
    secondary: {
      main: '#dc004e',
      light: '#f05545',
      dark: '#9a0036',
    },
    success: {
      main: '#4caf50',
    },
    warning: {
      main: '#ff9800',
    },
    error: {
      main: '#f44336',
    },
    info: {
      main: '#2196f3',
    },
  },
  typography: {
    fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
  },
});

