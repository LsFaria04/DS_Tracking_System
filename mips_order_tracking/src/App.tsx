import OrderTrackingPage from './pages/OrderTrackingPage';
import HomePage from './pages/HomePage';
import { Routes, Route } from 'react-router-dom';

export default function App() {
  return (
    <div>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/order/:id" element={<OrderTrackingPage />} />
      </Routes>
    </div>
  );
}