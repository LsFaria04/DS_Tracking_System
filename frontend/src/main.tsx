import React from 'react'
import ReactDOM from 'react-dom/client'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'

import AppLayout from './AppLayout'
import './index.css'
import HomePage from './pages/HomePage'
import OrderDetailPage from './pages/OrderTrackingPage'
import OrdersPage from './pages/OrdersPage'

const router = createBrowserRouter([
    {
        path: '/',
        element: <AppLayout />,
        children: [
            {
                path: '/',
                element: <HomePage />,
            },
            {
                path: '/orders',
                element: <OrdersPage />,
            },
            {
                path: '/order/:id',
                element: <OrderDetailPage />,
            },
        ],
    },
])

ReactDOM.createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
        <RouterProvider router={router} />
    </React.StrictMode>,
)