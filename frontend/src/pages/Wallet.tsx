import React, { useState, useEffect } from 'react';
import { api } from '../contexts/AuthContext';
import { Wallet as WalletIcon, Plus, Minus, Trash2 } from 'lucide-react';
import toast from 'react-hot-toast';

interface WalletResponse {
  id: string;
  currency: string;
  balance: string;
  locked_balance: string;
}

interface Wallet {
  id: string;
  currency: string;
  balance: number;
  locked_balance: number;
}

const WalletPage: React.FC = () => {
  const [wallets, setWallets] = useState<Wallet[]>([]);
  const [loading, setLoading] = useState(true);
  const [creating, setCreating] = useState(false);
  const [newCurrency, setNewCurrency] = useState('');

  const [searchedWallet, setSearchedWallet] = useState<Wallet | null>(null);

  useEffect(() => {
    const fetchWallets = async () => {
      try {
        const res = await api.get<{ wallets: WalletResponse[] }>('/v1/wallets', {
          params: { user_email: localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user')!).email : '' }
        });
        const normalized = (res.data.wallets || []).filter(Boolean).map(w => ({
          id: w.id,
          currency: w.currency,
          balance: parseFloat(w.balance ?? "0"),
          locked_balance: parseFloat(w.locked_balance ?? "0")
        }));
        setWallets(normalized);
        console.log(res.data.wallets);
      } catch (err) {
        console.error(err);
      } finally {
      setLoading(false);
      }
    };

    fetchWallets();
  }, []);

  const handleCreateWallet = async () => {
    if (!newCurrency) {
      toast.error('Please enter a currency');
      return;
    }

    try {
      setCreating(true);

      const response = await api.post<{ wallet_id: string }>('/v1/wallets', {
        user_email: localStorage.getItem('user')
          ? JSON.parse(localStorage.getItem('user')!).email
          : '',
        currency: newCurrency,
      });

      const walletId = response.data.wallet_id;

      const walletResponse = await api.get<{ wallet: Wallet }>(`/v1/wallets/${walletId}`);
      const raw = walletResponse.data.wallet;
      const createdWallet: Wallet = {
        id: raw.id,
        currency: raw.currency,
        balance: Number(raw.balance),
        locked_balance: Number(raw.locked_balance),
      };
      setSearchedWallet(createdWallet);
      setWallets((prev) => [...prev, createdWallet]);

      toast.success('Wallet created & fetched successfully!');
      setNewCurrency('');
    } catch (error: any) {
      const message = error.response?.data?.error || 'Failed to create wallet';
      toast.error(message);
    } finally {
      setCreating(false);
    }
  };

  const handleDeleteWallet = async (walletId: string) => {
    try {
      await api.delete(`/v1/wallets/${walletId}`);
      toast.success('Wallet deleted successfully!');
      setWallets(wallets.filter((w) => w.id !== walletId));
      if (searchedWallet?.id === walletId) setSearchedWallet(null);
    } catch (error: any) {
      toast.error('Failed to delete wallet');
    }
  };

  const handleUpdateWallet = async (walletId: string, delta: number) => {
    try {
      const wallet = wallets.filter(Boolean).find((w) => w.id === walletId);
      if (!wallet) return;

      const updatedBalance = wallet.balance + delta;

      await api.patch<{ wallet: Wallet }>(`/v1/wallets/${walletId}`, {
        balance: updatedBalance,
        locked_balance: wallet.locked_balance,
        id: walletId,
      });

      toast.success(delta >= 0 ? 'Deposit successful!' : 'Withdrawal successful!');
    } catch (error: any) {
      console.error(error);
      toast.error('Failed to update wallet');
    }
  };

  const formatCurrency = (amount: number, currency: string) => {
    if (currency === 'USD') {
      return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'USD',
      }).format(amount);
    }
    return `${amount} ${currency}`;
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="fixed top-10 mt-10 left-1/2 transform -translate-x-1/3 w-full max-w-3xl px-6 py-6 h-[calc(100vh-6rem)] overflow-y-auto"> 
      {/* Header */}
      <div className="bg-white rounded-lg shadow p-6 mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Wallet</h1>
        <p className="text-gray-600 mt-1">Manage your cryptocurrency balances</p>

        {/* Create Wallet */}
        <div className="mt-4 flex space-x-2">
          <input
            type="text"
            value={newCurrency}
            onChange={(e) => setNewCurrency(e.target.value.toUpperCase())}
            placeholder="Currency (e.g. BTC, ETH, USD)"
            className="border rounded-md px-3 py-2 text-sm flex-1"
          />
          <button
            type="button"
            onClick={handleCreateWallet}
            disabled={creating}
            className="px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-md hover:bg-blue-700 transition-colors disabled:opacity-50"
          >
            {creating ? 'Creating...' : 'Create Wallet'}
          </button>
        </div>

        {/* Show created wallet immediately */}
        {searchedWallet && (
          <div className="mt-4 p-4 border rounded-md bg-gray-50">
            <h3 className="text-lg font-semibold">Latest Wallet:</h3>
            <p><strong>ID:</strong> {searchedWallet.id}</p>
            <p><strong>Currency:</strong> {searchedWallet.currency}</p>
            <p><strong>Balance:</strong> {formatCurrency(searchedWallet.balance, searchedWallet.currency)}</p>
            <p><strong>Locked:</strong> {formatCurrency(searchedWallet.locked_balance, searchedWallet.currency)}</p>
          </div>
        )}
      </div>

      {/* Wallet Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
        {Array.isArray(wallets) && wallets.filter(Boolean).map((wallet) => (
          <div
            key={wallet.id}
            className="bg-white rounded-lg shadow p-5 min-h-[320px] flex flex-col justify-between"
          >
            {/* Top content */}
            <div className="flex-1">
              <div className="flex items-center justify-between mb-5">
                <div className="flex items-center">
                  <div className="p-2 bg-blue-100 rounded-lg">
                    <WalletIcon className="h-6 w-6 text-blue-600" />
                  </div>
                  <div className="ml-3">
                    <h3 className="text-sm font-semibold text-gray-900">{wallet.currency}</h3>
                    <p className="text-sm text-gray-500">Available Balance</p>
                  </div>
                </div>
                <button
                  onClick={() => handleDeleteWallet(wallet.id)}
                  className="text-red-500 hover:text-red-700"
                >
                  <Trash2 className="h-5 w-5" />
                </button>
              </div>

              <div className="space-y-6">
                <div>
                  <p className="text-2xl font-bold text-gray-900">
                    {formatCurrency(wallet.balance, wallet.currency)}
                  </p>
                  <p className="text-sm text-gray-500">Available</p>
                </div>

                <div>
                  <p className="text-sm font-medium text-gray-700">
                    {formatCurrency(wallet.locked_balance, wallet.currency)}
                  </p>
                  <p className="text-sm text-gray-500">Locked in Orders</p>
                </div>
              </div>
            </div>

            {/* Action buttons */}
            <div className="flex space-x-2 mt-4">
              <button
                type="button"
                onClick={() => handleUpdateWallet(wallet.id, 10)}
                className="flex-1 flex items-center justify-center px-2 py-2 bg-green-600 text-white text-sm font-medium rounded-md hover:bg-green-700 transition-colors"
              >
                <Plus className="h-4 w-4 mr-2" />
                Deposit
              </button>
              <button
                type="button"
                onClick={() => handleUpdateWallet(wallet.id, -10)}
                className="flex-1 flex items-center justify-center px-2 py-2 bg-red-600 text-white text-sm font-medium rounded-md hover:bg-red-700 transition-colors"
              >
                <Minus className="h-4 w-4 mr-1" />
                Withdraw
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default WalletPage;