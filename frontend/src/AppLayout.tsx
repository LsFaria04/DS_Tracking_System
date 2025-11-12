import { Outlet } from 'react-router-dom';

export default function AppLayout() {
    return (
        <div className="app-container">
            {/* Example:
                <MyNavbar /> 
            */}

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