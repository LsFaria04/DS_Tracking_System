import { render, screen, waitFor } from '@testing-library/react';
import OrdersPage from './OrdersPage';
import { MemoryRouter } from 'react-router-dom';
import { vi } from 'vitest';

describe('OrdersPage', () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });

  it('renders a list of orders fetched from the API', async () => {
    const mockOrders = [
      { Id: 1, Tracking_Code: 'T-1', Delivery_Address: 'Address 1', Price: 10, Created_At: '2023-01-01', Products: [] },
      { Id: 2, Tracking_Code: 'T-2', Delivery_Address: 'Address 2', Price: 20, Created_At: '2023-01-02', Products: [] },
    ];

    global.fetch = vi.fn(() => Promise.resolve({ ok: true, json: () => Promise.resolve({ orders: mockOrders }) })) as unknown as typeof fetch;

    render(
      <MemoryRouter>
        <OrdersPage />
      </MemoryRouter>
    );

    // Wait for orders to be rendered
    await waitFor(() => {
      expect(screen.getByText('T-1')).toBeTruthy();
      expect(screen.getByText('Address 1')).toBeTruthy();
      expect(screen.getByText('T-2')).toBeTruthy();
      expect(screen.getByText('Address 2')).toBeTruthy();
    });
  });

  it('shows an error message when the API fails', async () => {
    global.fetch = vi.fn(() => Promise.reject(new Error('Network fail'))) as unknown as typeof fetch;

    render(
      <MemoryRouter>
        <OrdersPage />
      </MemoryRouter>
    );

    await waitFor(() => expect(screen.getByText('Unable to load orders')).toBeTruthy());
  });
});
