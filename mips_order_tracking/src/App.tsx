import OrderTrackingPage from './pages/OrderTrackingPage';
import HomePage from './pages/HomePage';
import { Routes, Route } from 'react-router-dom';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { theme } from './theme';

export default function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <div className="mips-tracking-root">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/order/:id" element={<OrderTrackingPage />} />
        </Routes>
      </div>
    </ThemeProvider>
  );
}