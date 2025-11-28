import { render, screen, waitFor } from '@testing-library/react';
import HomePage from './HomePage';
import { MemoryRouter } from 'react-router-dom';
import { vi } from 'vitest';

describe('HomePage (microfrontend)', () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });

  it('renders a list of orders fetched from the API', async () => {
    const mockOrders = [
      { Id: 1, Tracking_Code: 'M-T-1', Delivery_Address: 'M-Address 1', Price: 15, Created_At: '2023-01-01', Products: [] },
    ];

    global.fetch = vi.fn(() => Promise.resolve({ ok: true, json: () => Promise.resolve({ orders: mockOrders }) })) as unknown as typeof fetch;

    render(
      <MemoryRouter>
        <HomePage />
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('M-T-1')).toBeTruthy();
      expect(screen.getByText('M-Address 1')).toBeTruthy();
    });
  });

  it('shows an error message when the API fails', async () => {
    global.fetch = vi.fn(() => Promise.reject(new Error('Network fail'))) as unknown as typeof fetch;

    render(
      <MemoryRouter>
        <HomePage />
      </MemoryRouter>
    );

    await waitFor(() => expect(screen.getByText('Unable to load orders')).toBeTruthy());
  });
});
