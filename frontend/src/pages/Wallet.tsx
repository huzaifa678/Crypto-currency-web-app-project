import React, { useState, useEffect } from 'react';
import { api } from '../contexts/AuthContext';
import { Wallet, Plus, Minus, ArrowUpRight } from 'lucide-react';

interface WalletBalance {
  id: string;
  currency: string;
  balance: string;
  locked_balance: string;
}

const WalletPage: React.FC = () => {
  const [balances, setBalances] = useState<WalletBalance[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchWallet = async () => {
      try {
        setLoading(true);
        
        setBalances([
          { id: '1', currency: 'USD', balance: '10000.00', locked_balance: '500.00' },
          { id: '2', currency: 'BTC', balance: '0.5', locked_balance: '0.05' },
          { id: '3', currency: 'ETH', balance: '5.0', locked_balance: '0.5' },
        ]);
      } catch (error) {
        console.error('Error fetching wallet:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchWallet();
  }, []);

  const formatCurrency = (amount: string, currency: string) => {
    if (currency === 'USD') {
      return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'USD',
      }).format(parseFloat(amount));
    }
    return `${parseFloat(amount).toFixed(8)} ${currency}`;
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
  <div className="fixed top-10 right-20 px-6 py-6 mt-10 mr-10 lg:ml-64">
    <div className="bg-white rounded-lg shadow p-6 mb-6">
      <h1 className="text-2xl font-bold text-gray-900">Wallet</h1>
      <p className="text-gray-600 mt-1">
        Manage your cryptocurrency balances
      </p>
    </div>

    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
      {balances.map((balance) => (
        <div
          key={balance.id}
          className="bg-white rounded-lg shadow p-6 min-h-[320px] flex flex-col"
        >
          {/* Top content */}
          <div className="flex-1">
            <div className="flex items-center justify-between mb-5">
              <div className="flex items-center">
                <div className="p-2 bg-blue-100 rounded-lg">
                  <Wallet className="h-6 w-6 text-blue-600" />
                </div>
                <div className="ml-3">
                  <h3 className="text-sm font-semibold text-gray-900">
                    {balance.currency}
                  </h3>
                  <p className="text-sm text-gray-500">Available Balance</p>
                </div>
              </div>
            </div>

            <div className="space-y-6">
              <div>
                <p className="text-2xl font-bold text-gray-900">
                  {formatCurrency(balance.balance, balance.currency)}
                </p>
                <p className="text-sm text-gray-500">Available</p>
              </div>

              <div>
                <p className="text-sm font-medium text-gray-700">
                  {formatCurrency(balance.locked_balance, balance.currency)}
                </p>
                <p className="text-sm text-gray-500">Locked in Orders</p>
              </div>
            </div>
          </div>

          {/* Buttons */}
          <div className="flex space-x-2 mt-6">
            <button
              type="button"
              className="flex-1 flex items-center justify-center px-2 py-2 bg-green-600 text-white text-sm font-medium rounded-md hover:bg-green-700 transition-colors"
            >
              <Plus className="h-4 w-4 mr-2" />
              Deposit
            </button>
            <button
              type="button"
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