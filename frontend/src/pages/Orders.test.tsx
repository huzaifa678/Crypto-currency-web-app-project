import { render, screen, waitFor, fireEvent, within } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { vi } from 'vitest';
import Orders from './Orders';

const mockUseMarkets = {
  market: [{ market_id: 'btc-usdt' }],
};
vi.mock('../contexts/MarketContext', () => ({
  useMarkets: () => mockUseMarkets,
}));

const mockSetOrder = vi.fn();
vi.mock('../contexts/OrderContext', () => ({
  useOrder: () => ({
    setOrder: mockSetOrder,
  }),
}));

const mockToastSuccess = vi.fn();
const mockToastError = vi.fn();
vi.mock('react-hot-toast', () => ({
  __esModule: true,
  default: {
    success: (...args: any[]) => mockToastSuccess(...args),
    error: (...args: any[]) => mockToastError(...args),
  },
}));

const mockApiGet = vi.fn();
const mockApiPost = vi.fn();
const mockApiDelete = vi.fn();

vi.mock('../contexts/AuthContext', () => ({
  api: {
    get: (...args: any[]) => mockApiGet(...args),
    post: (...args: any[]) => mockApiPost(...args),
    delete: (...args: any[]) => mockApiDelete(...args),
  },
}));


vi.stubGlobal('localStorage', {
  getItem: vi.fn(() =>
    JSON.stringify({ username: 'testuser', email: 'test@example.com' })
  ),
  setItem: vi.fn(),
  removeItem: vi.fn(),
});

describe('Orders Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockOrders = [
    {
      id: '1',
      market_id: 'btc-usdt',
      type: 'BUY',
      status: 'OPEN',
      price: '50000',
      amount: '0.1',
      created_at: '2025-01-01T00:00:00Z',
    },
  ];

  it('shows loading spinner then fetches & renders orders', async () => {
    mockApiGet.mockResolvedValueOnce({
      data: { orders: mockOrders },
    });

    render(
      <MemoryRouter>
        <Orders />
      </MemoryRouter>
    );

    expect(
      screen.getByRole('status', { hidden: true })
    ).toBeInTheDocument();

    await waitFor(() =>
      expect(mockApiGet).toHaveBeenCalledWith('/v1/orders', {
        params: { username: 'testuser' },
      })
    );

    const table = screen.getByRole('table', { name: /orders-table/i });

    expect(within(table).getByText('#1')).toBeInTheDocument();
    expect(within(table).getByText('BTC-USDT')).toBeInTheDocument();
    expect(within(table).getByText('BUY')).toBeInTheDocument();
  });

  it('creates a new order successfully', async () => {
    mockApiGet.mockResolvedValueOnce({ data: { orders: [] } });

    mockApiGet.mockResolvedValueOnce({
      data: { market: { market_id: 'btc-usdt' } },
    });

    mockApiPost.mockResolvedValueOnce({
      data: { order_id: '99' },
    });

    mockApiGet.mockResolvedValueOnce({
      data: {
        order: {
          id: '99',
          market_id: 'btc-usdt',
          type: 'BUY',
          status: 'OPEN',
          price: '100',
          amount: '1',
          created_at: '2025-01-01T00:00:00Z',
        },
      },
    });

    render(
      <MemoryRouter>
        <Orders />
      </MemoryRouter>
    );

    await waitFor(() => screen.getByText('Create Order'));

    fireEvent.change(screen.getByPlaceholderText('Base Currency'), {
      target: { value: 'BTC' },
    });
    fireEvent.change(screen.getByPlaceholderText('Quote Currency'), {
      target: { value: 'USDT' },
    });
    fireEvent.change(screen.getByPlaceholderText('Price'), {
      target: { value: '100' },
    });
    fireEvent.change(screen.getByPlaceholderText('Amount'), {
      target: { value: '1' },
    });

    fireEvent.click(screen.getByText('Create Order'));

    await waitFor(() =>
      expect(mockApiPost).toHaveBeenCalledWith('/v1/orders', {
        user_email: 'test@example.com',
        market_id: 'btc-usdt',
        base_currency: 'BTC',
        quote_currency: 'USDT',
        type: 'BUY',
        price: '100',
        amount: '1',
      })
    );

    await waitFor(() => {
      const elements = screen.getAllByText('#99');
      expect(elements.length).toBeGreaterThan(0);
    });

    expect(mockToastSuccess).toHaveBeenCalled();
  });

  it('cancels an order', async () => {
    mockApiGet.mockResolvedValueOnce({ data: { orders: mockOrders } });
    mockApiDelete.mockResolvedValueOnce({});

    render(
      <MemoryRouter>
        <Orders />
      </MemoryRouter>
    );

    await waitFor(() => screen.getByText('#1'));

    const cancelBtn = screen.getByRole('button', { name: '' });

    fireEvent.click(cancelBtn);

    expect(mockApiDelete).toHaveBeenCalledWith('/v1/orders/1');

    await waitFor(() =>
      expect(screen.queryByText('#1')).toBeNull()
    );

    expect(mockToastSuccess).toHaveBeenCalled();
  });
});
