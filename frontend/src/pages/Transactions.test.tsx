import { render, screen, waitFor, fireEvent, within } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { vi } from "vitest";
import Transactions from "./Transactions";

vi.mock("../contexts/AuthContext", () => {
  return {
    useAuth: () => ({
      user: {
        username: "huzaifa210",
        email: "huzaifa210@example.com",
      },
    }),
    api: {
      get: (...args: any[]) => mockApiGet(...args),
      post: (...args: any[]) => mockApiPost(...args),
      patch: (...args: any[]) => mockApiPatch(...args),
      delete: (...args: any[]) => mockApiDelete(...args),
    },
  };
});

const mockApiGet = vi.fn();
const mockApiPost = vi.fn();
const mockApiPatch = vi.fn();
const mockApiDelete = vi.fn();

vi.stubGlobal("crypto", {
  randomUUID: () => "mock-uuid",
});

describe("Transactions Component", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockTransactions = [
    {
      transaction_id: "tx1",
      type: "deposit",
      currency: "USD",
      amount: "100",
      status: "completed",
      created_at: "2025-01-01T00:00:00Z",
    },
  ];

  it("fetches and renders transactions", async () => {
    mockApiGet.mockResolvedValueOnce({
      data: { transactions: mockTransactions },
    });

    render(
      <MemoryRouter>
        <Transactions />
      </MemoryRouter>
    );

    expect(
      screen.getByRole("status", { hidden: true })
    ).toBeInTheDocument();

    await waitFor(() =>
      expect(mockApiGet).toHaveBeenCalledWith(
        "/v1/transactions/list/huzaifa210@example.com"
      )
    );

    const table = screen.getByRole("table");

    expect(within(table).getByText("deposit")).toBeInTheDocument();
    expect(within(table).getByText("USD")).toBeInTheDocument();
  });

  it("creates a new transaction", async () => {
    mockApiGet.mockResolvedValueOnce({ data: { transactions: [] } });

    mockApiPost.mockResolvedValueOnce({
      data: {
        transaction: {
          transaction_id: "newTx",
          type: "deposit",
          currency: "USD",
          amount: "50",
          status: "pending",
          created_at: "2025-01-01T00:00:00Z",
        },
      },
    });

    mockApiPatch.mockResolvedValueOnce({});

    render(
      <MemoryRouter>
        <Transactions />
      </MemoryRouter>
    );

    fireEvent.change(screen.getByPlaceholderText("Amount"), {
      target: { name: "amount", value: "50" },
    });

    fireEvent.change(screen.getByPlaceholderText("Wallet Address"), {
      target: { name: "address", value: "wallet123" },
    });

    fireEvent.click(screen.getByText("Submit"));

    await waitFor(() =>
      expect(mockApiPost).toHaveBeenCalledWith("/v1/transactions", {
        username: "huzaifa210",
        user_email: "huzaifa210@example.com",
        type: "deposit",
        currency: "USD",
        amount: "50",
        address: "wallet123",
        tx_hash: "mock-uuid",
      })
    );

    expect(mockApiPatch).toHaveBeenCalledWith("/v1/transactions/newTx", {
      status: "COMPLETED",
    });

    const table = screen.getByRole("table");
    expect(within(table).getByText("deposit")).toBeInTheDocument();
  });

  it("deletes a transaction", async () => {
    mockApiGet.mockResolvedValueOnce({ data: { transactions: mockTransactions } });
    mockApiDelete.mockResolvedValueOnce({});

    render(
      <MemoryRouter>
        <Transactions />
      </MemoryRouter>
    );

    await waitFor(() => screen.getByText("deposit"));

    fireEvent.click(screen.getByText("Delete"));

    expect(mockApiDelete).toHaveBeenCalledWith("/v1/transactions/tx1");

    await waitFor(() =>
      expect(screen.queryByText("deposit")).not.toBeInTheDocument()
    );
  });

  it("handles fetch error", async () => {
    mockApiGet.mockRejectedValueOnce(new Error("API error"));

    render(
      <MemoryRouter>
        <Transactions />
      </MemoryRouter>
    );

    await waitFor(() =>
      expect(screen.getByText("Failed to fetch transactions")).toBeInTheDocument()
    );
  });
});
