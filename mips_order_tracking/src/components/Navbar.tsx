import { Link } from 'react-router-dom';

export default function MyNavbar() {
    return (
        <nav className="w-full bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-800">
            <div className="max-w-6xl mx-auto px-6 py-3 flex items-center justify-between">
                <div className="flex items-center gap-4">
                    <Link to="/" className="text-lg font-semibold text-gray-900 dark:text-white">Tracking Status</Link>
                    <Link to="/orders" className="text-sm text-gray-500 hover:text-gray-700 dark:hover:text-gray-300">Orders</Link>
                </div>
            </div>
        </nav>
    );
}
