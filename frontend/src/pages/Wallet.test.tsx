import { render, screen, fireEvent, waitFor, cleanup, within } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { vi } from "vitest";
import WalletPage from "./Wallet";

afterEach(() => {
  cleanup();
  vi.clearAllMocks();
});

vi.mock("react-hot-toast", () => ({
  default: {
    success: vi.fn(),
    error: vi.fn(),
  },
}));

const mockGet = vi.fn();
const mockPost = vi.fn();
const mockPatch = vi.fn();
const mockDelete = vi.fn();

vi.mock("../contexts/AuthContext", () => ({
  api: {
    get: (...args: any[]) => mockGet(...args),
    post: (...args: any[]) => mockPost(...args),
    patch: (...args: any[]) => mockPatch(...args),
    delete: (...args: any[]) => mockDelete(...args),
  },
}));

beforeEach(() => {
  vi.clearAllMocks();
  Storage.prototype.getItem = vi.fn(() =>
    JSON.stringify({ email: "test@example.com" })
  );
});


const renderWalletPage = () =>
  render(
    <MemoryRouter>
      <WalletPage />
    </MemoryRouter>
  );


describe("WalletPage", () => {
  test("fetches and displays wallets on load", async () => {
    mockGet.mockResolvedValueOnce({
      data: {
        wallets: [
          { id: "1", currency: "USD", balance: "100", locked_balance: "20" },
        ],
      },
    });

    renderWalletPage();

    expect(screen.getByRole("status")).toBeInTheDocument(); 

    expect(mockGet).toHaveBeenCalledWith(
      "/v1/wallets",
      expect.any(Object)
    );

    expect(await screen.findByText("$100.00")).toBeInTheDocument();
  });

  test("creates a wallet successfully", async () => {
    mockGet.mockResolvedValueOnce({ data: { wallets: [] } });

    mockPost.mockResolvedValueOnce({ data: { wallet_id: "abc123" } });

    mockGet.mockResolvedValueOnce({
      data: {
        wallet: {
          id: "abc123",
          currency: "BTC",
          balance: 50,
          locked_balance: 0,
        },
      },
    });

    renderWalletPage();

    await waitFor(() =>
        expect(screen.queryByRole("status")).not.toBeInTheDocument()
    );

    fireEvent.change(screen.getByPlaceholderText(/currency/i), {
      target: { value: "btc" },
    });
    fireEvent.click(screen.getByText(/create wallet/i));

    await waitFor(() => {
      expect(mockPost).toHaveBeenCalledWith(
        "/v1/wallets",
        expect.objectContaining({
          currency: "BTC",
          user_email: "test@example.com",
        })
      );
    });

    expect(await screen.findByText("Latest Wallet:")).toBeInTheDocument();

    const latestWalletBox = screen.getByText("Latest Wallet:");
    const container = latestWalletBox.closest("div");


    expect(within(container!).getByText("Currency:")).toBeInTheDocument();
    expect(within(container!).getByText("BTC")).toBeInTheDocument();
  });

  test("deposits money into a wallet", async () => {
    mockGet.mockResolvedValueOnce({
      data: {
        wallets: [
          { id: "w1", currency: "USD", balance: "100", locked_balance: "0" },
        ],
      },
    });

    mockPatch.mockResolvedValueOnce({});

    renderWalletPage();

    fireEvent.click(await screen.findByText(/deposit/i));

    await waitFor(() =>
      expect(mockPatch).toHaveBeenCalledWith("/v1/wallets/w1", {
        id: "w1",
        balance: 110,
        locked_balance: 0,
      })
    );
  });

  test("withdraws money from a wallet", async () => {
    mockGet.mockResolvedValueOnce({
      data: {
        wallets: [
          { id: "w1", currency: "USD", balance: "100", locked_balance: "0" },
        ],
      },
    });

    mockPatch.mockResolvedValueOnce({});

    renderWalletPage();

    fireEvent.click(await screen.findByText(/withdraw/i));

    await waitFor(() =>
      expect(mockPatch).toHaveBeenCalledWith("/v1/wallets/w1", {
        id: "w1",
        balance: 90,
        locked_balance: 0,
      })
    );
  });

  test("deletes a wallet", async () => {
    mockGet.mockResolvedValueOnce({
      data: {
        wallets: [
          { id: "bye", currency: "USD", balance: "10", locked_balance: "0" },
        ],
      },
    });

    mockDelete.mockResolvedValueOnce({});

    renderWalletPage();

    fireEvent.click(await screen.findByRole("button", { name: "" })); 

    await waitFor(() =>
      expect(mockDelete).toHaveBeenCalledWith("/v1/wallets/bye")
    );
  });
});
