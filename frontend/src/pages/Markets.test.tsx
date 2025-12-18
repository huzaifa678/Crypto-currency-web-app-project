import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { vi } from 'vitest';
import Markets from './Markets';

const mockSetMarket = vi.fn();
const mockUseMarkets = {
  market: [],
  setMarket: mockSetMarket,
};

vi.mock('../contexts/MarketContext', () => ({
  useMarkets: () => mockUseMarkets,
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

vi.mock('react-router-dom', async (importOriginal) => {
  const actual = await importOriginal<any>();
  return {
    ...actual,
    Link: ({ children, to }: any) => <a href={to}>{children}</a>,
  };
});

describe('Markets Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockMarkets = [
    {
      market_id: '1',
      name: 'BTC/USDT',
      base_currency: 'BTC',
      quote_currency: 'USDT',
      min_order_amount: 0.001,
      price_precision: 2,
      created_at: '2025-01-01T10:00:00Z',
    },
    {
      market_id: '2',
      name: 'ETH/USDT',
      base_currency: 'ETH',
      quote_currency: 'USDT',
      min_order_amount: 0.01,
      price_precision: 3,
      created_at: '2025-01-02T10:00:00Z',
    },
  ];

  it('fetches markets and renders table rows', async () => {
    mockApiGet.mockResolvedValueOnce({
      data: { markets: mockMarkets },
    });

    render(
      <MemoryRouter>
        <Markets />
      </MemoryRouter>
    );

    expect(screen.getByRole('status')).toBeInTheDocument();

    await waitFor(() => {
      expect(mockApiGet).toHaveBeenCalledWith('/v1/markets', expect.any(Object));
    });

    expect(screen.getByText('BTC')).toBeInTheDocument();
    expect(screen.getByText('ETH')).toBeInTheDocument();
  });


  it('filters markets by search input', async () => {
    mockApiGet.mockResolvedValueOnce({
      data: { markets: mockMarkets },
    });

    render(
      <MemoryRouter>
        <Markets />
      </MemoryRouter>
    );

    await waitFor(() => screen.getByText('BTC'));

    fireEvent.change(screen.getByPlaceholderText(/search markets/i), {
      target: { value: 'btc' },
    });

    expect(screen.getByText('BTC')).toBeInTheDocument();
    expect(screen.queryByText('ETH')).toBeNull();
  });

  it('sorts markets by base currency when clicking header', async () => {
    mockApiGet.mockResolvedValueOnce({
      data: { markets: mockMarkets },
    });

    render(
      <MemoryRouter>
        <Markets />
      </MemoryRouter>
    );

    await waitFor(() => screen.getByText('BTC'));

    const baseHeader = screen.getByText('Base');

    fireEvent.click(baseHeader);

    const rows = screen.getAllByRole('row');

    expect(rows[1]).toHaveTextContent('BTC');

    fireEvent.click(baseHeader);

    const rowsDesc = screen.getAllByRole('row');
    expect(rowsDesc[1]).toHaveTextContent('ETH');
  });

  it('calls deleteMarket and removes row', async () => {
    mockApiGet.mockResolvedValueOnce({
      data: { markets: mockMarkets },
    });

    mockApiDelete.mockResolvedValueOnce({});

    render(
      <MemoryRouter>
        <Markets />
      </MemoryRouter>
    );

    await waitFor(() => screen.getByText('BTC'));

    const deleteButton = screen.getAllByText(/delete/i)[0];

    fireEvent.click(deleteButton);

    expect(mockApiDelete).toHaveBeenCalledWith('/v1/markets/1');

    await waitFor(() => {
      expect(screen.queryByText('BTC')).toBeNull();
    });
  });

  it('creates market and adds new row', async () => {
    mockApiGet.mockResolvedValueOnce({
      data: { markets: [] },
    });

    mockApiPost.mockResolvedValueOnce({
      data: { market_id: '99' },
    });

    mockApiGet.mockResolvedValueOnce({
      data: {
        market: {
          market_id: '99',
          name: 'SOL/USDT',
          base_currency: 'SOL',
          quote_currency: 'USDT',
          min_order_amount: 1,
          price_precision: 2,
          created_at: '2025-02-02T00:00:00Z',
        },
      },
    });

    render(
      <MemoryRouter>
        <Markets />
      </MemoryRouter>
    );

    await waitFor(() => screen.getByText('Create Market'));

    fireEvent.change(screen.getByPlaceholderText('Base Currency'), { target: { value: 'SOL' } });
    fireEvent.change(screen.getByPlaceholderText('Quote Currency'), { target: { value: 'USDT' } });
    fireEvent.change(screen.getByPlaceholderText('Min Order Amount'), { target: { value: '1' } });

    fireEvent.click(screen.getByText('Create Market'));

    expect(mockApiPost).toHaveBeenCalledWith('/v1/markets', {
      base_currency: 'SOL',
      quote_currency: 'USDT',
      min_order_amount: '1',
      price_precision: 2,
    });

    await waitFor(() => {
      expect(screen.getByText('SOL')).toBeInTheDocument();
    });
  });
});
