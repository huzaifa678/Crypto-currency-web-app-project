import React, { useEffect, useState } from "react";
import { ArrowDownRight, ArrowUpRight } from "lucide-react";
import { useAuth, api } from "../contexts/AuthContext";
import { useOrder } from '../contexts/OrderContext';

interface Trade {
  tradeId: string;
  username: string;
  buyOrderId: string;
  sellOrderId: string;
  marketId: string;
  price: number;
  amount: number;
  fee: number;
  createdAt: string;
}

const Trades: React.FC = () => {
  const { user } = useAuth();
  const { order, setOrder } = useOrder();
  const [trades, setTrades] = useState<Trade[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [formData, setFormData] = useState({
    username: "",
    buy_order_id: "",
    sell_order_id: "",
    market_id: "BTC-USDT",
    price: "",
    amount: "",
    fee: "0.001",
    side: "buy" as "buy" | "sell",
  });

  useEffect(() => {
    const defaultUsername =
      (user as any)?.username ||
      (user?.email ? user.email.split("@")[0] : "");
    setFormData((f) => ({ ...f, username: defaultUsername }));
  }, [user?.email]);

  useEffect(() => {
    const fetchTrades = async () => {
      try {
        setLoading(true);
        setError(null);
        const res = await api.get("/v1/trades/{formData.market_id}");
        const trades = res.data.trades.map((t: any) => ({
          tradeId: t.trade_id,
          username: t.username,
          buyOrderId: t.buy_order_id,
          sellOrderId: t.sell_order_id,
          marketId: t.market_id,
          price: parseFloat(t.price),
          amount: parseFloat(t.amount),
          fee: parseFloat(t.fee),
          createdAt: t.created_at,
        }));
        setTrades(trades);
      } catch (err) {
        console.error("Error fetching trades:", err);
        setError("Failed to fetch trades");
      } finally {
        setLoading(false);
      }
    };

    fetchTrades();
  }, []);

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    const { name, value } = e.target;
    setFormData((f) => ({ ...f, [name]: value }));
  };

  const formatDate = (dateString: string) =>
    new Date(dateString).toLocaleString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });

  const formatCurrency = (value: string) =>
    new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
      minimumFractionDigits: 2,
      maximumFractionDigits: 8,
    }).format(parseFloat(value || "0"));

  const formatAmount = (value: string) =>
    isNaN(parseFloat(value)) ? "0.00000000" : parseFloat(value).toFixed(8);

  const sideOfTrade = (t: Trade) =>
    t.buyOrderId && !t.sellOrderId
      ? "buy"
      : t.sellOrderId && !t.buyOrderId
      ? "sell"
      : "—";

  const createTrade = async () => {
    try {
      setLoading(true);
      setError(null);

      console.log("Creating trade with formData:", formData, "and order:", order);

      let activeOrder = order;

      if (!activeOrder) {
        const res = await api.get('/v1/orders', {
          params: { username: formData.username }
        });

        const orders = res.data.orders;
        if (!orders || orders.length === 0) {
          setError("No order available to create trade");
          return;
        }
        activeOrder = orders[orders.length - 1]; 
        setOrder(activeOrder); 
      }

      console.log("Using activeOrder for trade creation:", activeOrder);
      console.log("Active order type", activeOrder?.type);

      const payload = {
        username: formData.username,
        buy_order_id: activeOrder?.type === "BUY" ? activeOrder.id : '',
        sell_order_id: activeOrder?.type === "SELL" ? activeOrder.id : '',
        market_id: activeOrder?.market_id,
        price: formData.price,
        amount: formData.amount,
        fee: formData.fee,
      };

      console.log("Payload for creating trade:", payload);

      const res = await api.post("/v1/trades", payload);
      const t = res.data.trade;

      const mapped: Trade = {
        tradeId: t.trade_id,
        username: t.username,
        buyOrderId: t.buy_order_id,
        sellOrderId: t.sell_order_id,
        marketId: t.market_id,
        price: t.price,
        amount: t.amount,
        fee: t.fee,
        createdAt: t.created_at,
      };

      setTrades((prev) => [mapped, ...prev]);

      setFormData((f) => ({ ...f, price: "", amount: "" }));
    } catch (err) {
      console.error("Error creating trade:", err);
      setError("Failed to create trade");
    } finally {
      setLoading(false);
    }
  };

  const deleteTrade = async (id: string) => {
    if (!id) return;
    try {
      await api.delete(`/v1/trades/${id}`);
      setTrades((prev) => prev.filter((t) => t.tradeId !== id));
    } catch (err) {
      console.error("Error deleting trade:", err);
    }
  };

  return (
    <div className="fixed top-20 left-1/2 -translate-x-1/2 max-w-5xl px-6 py-10 ml-10">
      {/* Create Trade */}
      <div className="bg-white p-4 rounded-lg shadow mb-6 ml-10">
        <h2 className="text-lg font-bold mb-4">New Trade</h2>
        <div className="grid grid-cols-1 md:grid-cols-7 gap-3">
          <input
            name="username"
            placeholder="Username"
            value={formData.username}
            onChange={handleChange}
            className="border px-3 py-2 rounded"
          />
          <select
            name="side"
            value={formData.side}
            onChange={handleChange}
            className="border px-3 py-2 rounded"
          >
            <option value="buy">Buy</option>
            <option value="sell">Sell</option>
          </select>
          <input
            name="market_id"
            placeholder="Market ID (e.g. BTC-USDT)"
            value={formData.market_id}
            onChange={handleChange}
            className="border px-3 py-2 rounded"
          />
          <input
            type="number"
            name="price"
            placeholder="Price"
            value={formData.price}
            onChange={handleChange}
            className="border px-3 py-2 rounded"
          />
          <input
            type="number"
            name="amount"
            placeholder="Amount"
            value={formData.amount}
            onChange={handleChange}
            className="border px-3 py-2 rounded"
          />
          <input
            type="number"
            name="fee"
            placeholder="Fee"
            value={formData.fee}
            onChange={handleChange}
            className="border px-3 py-2 rounded"
          />
          <button
            onClick={createTrade}
            className="bg-blue-600 text-white px-4 py-2 rounded-lg shadow hover:bg-blue-700"
          >
            Create
          </button>
        </div>
      </div>

      {/* Trades Table */}
      <div className="bg-white rounded-lg shadow overflow-hidden ml-10">
        {error && <p className="text-red-500 p-4">{error}</p>}
        {loading && (
          <div className="flex items-center justify-center h-32">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          </div>
        )}

        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Side</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Market</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Price</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Amount</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Fee</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">User</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Trade ID</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Date</th>
                <th className="px-6 py-3"></th>
              </tr>
            </thead>

            <tbody className="bg-white divide-y divide-gray-200">
              {trades.map((t) => {
                const side = sideOfTrade(t);
                return (
                  <tr key={t.tradeId} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap flex items-center">
                      {side === "buy" ? (
                        <ArrowDownRight className="h-5 w-5 text-green-500 mr-2" />
                      ) : side === "sell" ? (
                        <ArrowUpRight className="h-5 w-5 text-red-500 mr-2" />
                      ) : (
                        <span className="w-5 h-5 mr-2" />
                      )}
                      <span className="text-sm font-medium text-gray-900 capitalize">{side}</span>
                    </td>

                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                      {t.marketId}
                    </td>

                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      {formatCurrency(String(t.price))}
                    </td>

                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      {formatAmount(String(t.amount))}
                    </td>

                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      {formatAmount(String(t.fee))}
                    </td>

                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {t.username}
                    </td>

                    <td className="px-6 py-4 whitespace-nowrap text-xs text-gray-700">
                      {t.tradeId}
                    </td>

                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {formatDate(t.createdAt)}
                    </td>

                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm">
                      <button
                        onClick={() => deleteTrade(t.tradeId)}
                        className="text-red-600 hover:text-red-800"
                      >
                        Delete
                      </button>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>

        {trades.length === 0 && !loading && (
          <div className="text-center py-12">
            <p className="text-gray-500">No trades yet. Create one above.</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default Trades;