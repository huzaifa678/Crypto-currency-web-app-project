import React, { useState, useEffect } from "react";
import { ArrowUpRight, ArrowDownRight } from "lucide-react";
import { useAuth, api } from "../contexts/AuthContext";

interface Transaction {
  transactionId: string;
  type: string;
  currency: string;
  amount: string;
  status: string;
  createdAt: string;
}

const Transactions: React.FC = () => {
  const { user } = useAuth();
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [formData, setFormData] = useState({
    type: "deposit",
    currency: "USD",
    amount: "",
    address: "",
  });

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const fetchTransactions = async (email?: string) => {
    try {
      setLoading(true);
      setError(null);

      const res = await api.get(`/v1/transactions/list/${email}`); 
      const mapped = (res.data.transactions || []).map((t: any) => ({
      transactionId: t.transaction_id,
      type: t.type.toLowerCase(),
      currency: t.currency,
      amount: t.amount,
      status: t.status.toLowerCase() || "pending",
      createdAt: t.created_at,
      }));

      console.log(res.data)
      setTransactions(mapped);
    } catch (err: any) {
      console.error("Error fetching transactions:", err);
      setError("Failed to fetch transactions");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (!user?.email) return;

    fetchTransactions(user.email);

    const interval = setInterval(() => {
      fetchTransactions(user.email);
    }, 5000);

    return () => clearInterval(interval);
  }, [user?.email]);

  const fetchTransactionById = async (id: string) => {
    try {
      setLoading(true);
      setError(null);
      const res = await api.get(`/v1/transactions/${id}`);
      const t = res.data.transaction;

      if (t) {
      const mapped = {
        transactionId: t.transaction_id,
        type: t.type,
        currency: t.currency,
        amount: t.amount,
        status: t.status || "pending",
        createdAt: t.created_at,
      };

      setTransactions((prev) => [...prev, mapped]);
    }

    } catch (err) {
      console.error("Error fetching transaction:", err);
      setError("Failed to fetch transaction");
    } finally {
      setLoading(false);
    }
  };

  const createTransaction = async () => {
    try {
      setLoading(true);
      setError(null);

      const res = await api.post("/v1/transactions", {
        username: user?.username,
        user_email: user?.email,
        type: formData.type,
        currency: formData.currency,
        amount: formData.amount,
        address: formData.address,
        tx_hash: crypto.randomUUID(),
      });

      const createdTx = res.data.transaction;

      console.log(createdTx.transaction_id)

      await api.patch(`/v1/transactions/${createdTx.transaction_id}`, {
        status: "COMPLETED",
      });

      setTransactions((prev) => [
        ...prev,
        {
          transactionId: createdTx.transaction_id,
          type: createdTx.type,
          currency: createdTx.currency,
          amount: createdTx.amount,
          status: "completed", 
          createdAt: createdTx.created_at,
        },
      ]);

    } catch (err) {
      console.error("Error creating transaction:", err);
      setError("Failed to create transaction");
    } finally {
      setLoading(false);
    }
  };

  const deleteTransaction = async (id: string) => {
    if (!id || id === "undefined") {
      console.error("Invalid transaction ID:", id);
      return;
    }

    try {
      await api.delete(`/v1/transactions/${id}`);
      setTransactions((t) => t.filter((tx) => tx.transactionId !== id));
    } catch (err) {
      console.error("Error deleting transaction:", err);
    }
  };

  const formatCurrency = (amount: string, currency: string) => {
    if (currency === "USD") {
      return new Intl.NumberFormat("en-US", {
        style: "currency",
        currency: "USD",
      }).format(parseFloat(amount));
    }
    return `${parseFloat(amount).toFixed(8)} ${currency}`;
  };

  const formatDate = (dateString: string) =>
    new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });

  return (
    <div className="fixed top-20 left-1/2 -translate-x-1/2 max-w-4xl px-6 py-10 ml-10">
      {/* New Transaction Form */}
      <div className="bg-white p-4 rounded-lg shadow mb-6 ml-10">
        <h2 className="text-lg font-bold mb-4">New Transaction</h2>
        <div className="flex space-x-4">
          <select
            name="type"
            value={formData.type}
            onChange={handleChange}
            className="border px-3 py-2 rounded"
          >
            <option value="deposit">Deposit</option>
            <option value="withdrawal">Withdrawal</option>
          </select>

          <input
            type="number"
            name="amount"
            placeholder="Amount"
            value={formData.amount}
            onChange={handleChange}
            className="border px-3 py-2 rounded"
          />

          <input
            type="text"
            name="currency"
            placeholder="Currency (e.g. USD)"
            value={formData.currency}
            onChange={handleChange}
            className="border px-3 py-2 rounded"
          />

          <input
            type="text"
            name="address"
            placeholder="Wallet Address"
            value={formData.address}
            onChange={handleChange}
            className="border px-3 py-2 rounded"
          />

          <button
            onClick={createTransaction}
            className="bg-blue-600 text-white px-4 py-2 rounded-lg shadow hover:bg-blue-700"
          >
            Submit
          </button>
        </div>
      </div>

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
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Type</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Amount</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Currency</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Date</th>
                <th className="px-6 py-3"></th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {transactions.map((transaction) => (
                <tr key={transaction.transactionId} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap flex items-center">
                    {transaction.type === "deposit" ? (
                      <ArrowDownRight className="h-5 w-5 text-green-500 mr-2" />
                    ) : (
                      <ArrowUpRight className="h-5 w-5 text-red-500 mr-2" />
                    )}
                    <span className="text-sm font-medium text-gray-900 capitalize">
                      {transaction.type}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    {transaction.type === "deposit" ? "+" : "-"}
                    {formatCurrency(transaction.amount, transaction.currency)}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {transaction.currency}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm">
                    <span
                      className={`px-2 py-1 rounded text-white ${
                        transaction.status === "completed"
                          ? "bg-green-500"
                          : transaction.status === "pending"
                          ? "bg-yellow-500"
                          : "bg-red-500"
                      }`}
                    >
                      {transaction.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {formatDate(transaction.createdAt)}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right text-sm">
                    <button
                      onClick={() => fetchTransactionById(transaction.transactionId)}
                      className="text-blue-600 hover:text-blue-800"
                    >
                        Fetch
                    </button>

                    <button
                      onClick={() => deleteTransaction(transaction.transactionId)}
                      className="text-red-600 hover:text-red-800"
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {transactions.length === 0 && !loading && (
          <div className="text-center py-12">
            <p className="text-gray-500">No transactions found.</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default Transactions;