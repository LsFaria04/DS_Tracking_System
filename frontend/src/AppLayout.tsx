import { Outlet } from 'react-router-dom';
import MyNavbar from './components/Navbar';

export default function AppLayout() {
    return (
        <div className="app-container">
            <MyNavbar />

            <main>
                {/* "page.tsx" content renders inside this <Outlet /> */}
                <Outlet />
            </main>

            {/* Example:
                <MyFooter /> 
            */}
        </div>
    );
}