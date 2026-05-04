import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import Trades from "./Trading";
import { MemoryRouter } from "react-router-dom";

const localStorageMock: Record<string, string> = {};
vi.stubGlobal("localStorage", {
  getItem: (k: string) => localStorageMock[k] ?? null,
  setItem: (k: string, v: string) => (localStorageMock[k] = v),
  removeItem: (k: string) => delete localStorageMock[k],
  clear: () => Object.keys(localStorageMock).forEach(k => delete localStorageMock[k]),
});

const mockApiGet = vi.fn();
const mockApiPost = vi.fn();
const mockApiDelete = vi.fn();

vi.mock("../contexts/AuthContext", () => ({
  useAuth: () => ({
    user: { email: "test@example.com", username: "testuser" },
  }),
  api: {
    get: (...args: any[]) => mockApiGet(...args),
    post: (...args: any[]) => mockApiPost(...args),
    delete: (...args: any[]) => mockApiDelete(...args),
  },
}));

vi.mock("../contexts/OrderContext", () => {
  return {
    useOrder: () => ({
      orders: [
        {
          id: "orderBUY",
          market_id: "marketABC",
          type: "BUY",
          status: "OPEN",
          price: "100",
          amount: "2",
          filled_amount: "0",
          created_at: "2025-01-01T00:00:00Z",
        },
        {
          id: "orderSELL",
          market_id: "marketABC",
          type: "SELL",
          status: "OPEN",
          price: "100",
          amount: "2",
          filled_amount: "0",
          created_at: "2025-01-01T00:00:00Z",
        },
      ],
      setOrders: vi.fn(),
    }),
  };
});

beforeEach(() => {
  vi.clearAllMocks();
  localStorage.clear();
  localStorage.setItem("user", JSON.stringify({ username: "testuser" }));
  localStorage.setItem("marketId", "marketABC");
});

describe("Trades Component", () => {
  it("fetches and renders trades on mount", async () => {
    mockApiGet.mockResolvedValueOnce({
      data: {
        trades: [
          {
            trade_id: "T1",
            username: "testuser",
            buy_order_id: "orderBUY",
            sell_order_id: "orderSELL",
            market_id: "marketABC",
            price: "100",
            amount: "2",
            fee: "0.002",
            created_at: "2025-01-01T00:00:00Z",
          },
        ],
      },
    });

    render(
      <MemoryRouter>
        <Trades />
      </MemoryRouter>
    );

    expect(await screen.findByText("T1")).toBeInTheDocument();
  });

  it("creates a new trade successfully", async () => {
    mockApiGet.mockResolvedValueOnce({ data: { trades: [] } });

    mockApiPost.mockResolvedValueOnce({
      data: {
        trade: {
          trade_id: "new123",
          username: "huzaifa210",
          buy_order_id: "orderBUY",
          sell_order_id: "orderSELL",
          market_id: "marketABC",
          price: "100",
          amount: "2",
          fee: "0.002",
          created_at: "2025-01-01T00:00:00Z",
        },
      },
    });

    mockApiGet.mockResolvedValueOnce({
      data: {
        trade: {
          trade_id: "new123",
          market_id: "marketABC",
        },
      },
    });

    render(
      <MemoryRouter>
        <Trades />
      </MemoryRouter>
    );

    fireEvent.change(screen.getByPlaceholderText("Price"), {
      target: { value: "100" },
    });

    fireEvent.change(screen.getByPlaceholderText("Amount"), {
      target: { value: "2" },
    });

    fireEvent.click(screen.getByText("Create"));

    await waitFor(() =>
      expect(mockApiPost).toHaveBeenCalledWith("/v1/trades", {
        username: "testuser",
        buyer_user_email: "test@example.com",
        seller_user_email: "",
        buy_order_id: "orderBUY",
        sell_order_id: "orderSELL",
        market_id: "marketABC",
        price: "100",
        amount: "2",
        fee: "0.001", // default
      })
    );

    expect(await screen.findByText("new123")).toBeInTheDocument();
  });

  it("deletes a trade", async () => {
    mockApiGet.mockResolvedValueOnce({
      data: {
        trades: [
          {
            trade_id: "DEL1",
            username: "user",
            buy_order_id: "orderBUY",
            sell_order_id: "orderSELL",
            market_id: "marketABC",
            price: "100",
            amount: "2",
            fee: "0.002",
            created_at: "2025-01-01T00:00:00Z",
          },
        ],
      },
    });

    render(
      <MemoryRouter>
        <Trades />
      </MemoryRouter>
    );

    expect(await screen.findByText("DEL1")).toBeInTheDocument();

    mockApiDelete.mockResolvedValueOnce({});

    fireEvent.click(screen.getByText("Delete"));

    await waitFor(() =>
      expect(mockApiDelete).toHaveBeenCalledWith("/v1/trades/DEL1")
    );

    expect(screen.queryByText("DEL1")).not.toBeInTheDocument();
  });
});
