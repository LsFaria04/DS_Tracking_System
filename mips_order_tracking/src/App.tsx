import OrderTrackingPage from './pages/OrderTrackingPage';
import HomePage from './pages/HomePage';
import OrdersPage from './pages/OrdersPage';
import { Routes, Route } from 'react-router-dom';
import MyNavbar from './components/Navbar';

export default function App() {
  return (
    <div>
      <MyNavbar />
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/orders" element={<OrdersPage />} />
        <Route path="/order/:id" element={<OrderTrackingPage />} />
      </Routes>
    </div>
  );
}