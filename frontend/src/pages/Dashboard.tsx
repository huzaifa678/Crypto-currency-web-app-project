import React, { useState, useEffect, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { api } from '../contexts/AuthContext';
import { useTradeStream } from '../components/TradeStream';

import { 
  TrendingUp, 
  TrendingDown, 
  DollarSign, 
  Wallet, 
  BarChart3,
  ArrowUpRight,
  ArrowDownRight,
  Activity
} from 'lucide-react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { TypeAnimation } from 'react-type-animation';

interface Market {
  id: string;
  name: string;
  base_currency: string;
  quote_currency: string;
  current_price: string;
  price_change_24h: string;
  volume_24h: string;
}

interface Order {
  id: string;
  market_id: string;
  type: string;
  status: string;
  price: string;
  amount: string;
  created_at: string;
}

interface Transaction {
  id: string;
  type: string;
  currency: string;
  amount: string;
  created_at: string;
}

const Dashboard: React.FC = () => {
  const { user } = useAuth();
  const [markets, setMarkets] = useState<Market[]>([]);
  const [recentOrders, setRecentOrders] = useState<Order[]>([]);
  const [recentTransactions, setRecentTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);
  const trades = useTradeStream();
  const [priceData, setPriceData] = useState<any[]>([]);

  useEffect(() => {
    setLoading(false)

    if (!trades || trades.length === 0) return;

    const newPoint: any = {
      time: new Date().toLocaleTimeString(),
    };


    trades.forEach((m: any ) => {
      if (m.base_currency === "BTC") {
        newPoint.BTC = parseFloat(m.current_price);
      } else if (m.base_currency === "ETH") {
        newPoint.ETH = parseFloat(m.current_price);
      }
    });

    setPriceData((prev) => {
      const updated = [...prev, newPoint];
      return updated.slice(-50); // keep last 50 points
    });

  }, [trades]);

  const formatCurrency = (amount: string, currency: string = 'USD') => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency,
    }).format(parseFloat(amount));
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="fixed top-20 left-30 right-30 bottom-0 px-5 py-5 bg-animated-gradient min-h-2 overflow-y-auto h-[calc(100vh-5rem)]">
      {/* Welcome Section */}
      <div className="sticky bg-white rounded-lg shadow p-6">
        <TypeAnimation sequence={[
          `Welcome back, ${user?.username}! 👋`,
          1500,
          `Your portfolio is up +12.5% today 📈`,
          1500,
          `24h profit: $1,250 💰`,
          1500,
          () => console.log('Typing loop completed!')
        ]}
        wrapper="h1"
        cursor={true}
        repeat={Infinity}
        style={{
          fontSize: '1.5rem',
          fontWeight: 'bold',
          color: '#111827'
        }} />
        <p className="text-gray-600 mt-2">
          Here are the recent crypto updates for today
        </p>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 py-5 px-5">
        <div className="bg-white rounded-lg shadow p-6 transform transition-all duration-300 hover:scale-105 hover:shadow-2xl hover:-translate-y-1">
          <div className="flex items-center">
            <div className="p-2 bg-green-100 rounded-lg">
              <TrendingUp className="h-6 w-6 text-green-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Portfolio Value</p>
              <p className="text-2xl font-bold text-gray-900">$24,500</p>
            </div>
          </div>
          <div className="flex items-center mt-4">
            <ArrowUpRight className="h-4 w-4 text-green-500" />
            <span className="text-sm text-green-500 ml-1">+12.5%</span>
            <span className="text-sm text-gray-500 ml-2">from last month</span>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center">
            <div className="p-2 bg-blue-100 rounded-lg">
              <DollarSign className="h-6 w-6 text-blue-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">24h Profit</p>
              <p className="text-2xl font-bold text-gray-900">+$1,250</p>
            </div>
          </div>
          <div className="flex items-center mt-4">
            <ArrowUpRight className="h-4 w-4 text-green-500" />
            <span className="text-sm text-green-500 ml-1">+5.2%</span>
            <span className="text-sm text-gray-500 ml-2">from yesterday</span>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center">
            <div className="p-2 bg-purple-100 rounded-lg">
              <Wallet className="h-6 w-6 text-purple-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Active Orders</p>
              <p className="text-2xl font-bold text-gray-900">8</p>
            </div>
          </div>
          <div className="flex items-center mt-4">
            <ArrowDownRight className="h-4 w-4 text-red-500" />
            <span className="text-sm text-red-500 ml-1">-2</span>
            <span className="text-sm text-gray-500 ml-2">from yesterday</span>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center">
            <div className="p-2 bg-orange-100 rounded-lg">
              <Activity className="h-6 w-6 text-orange-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Total Trades</p>
              <p className="text-2xl font-bold text-gray-900">156</p>
            </div>
          </div>
          <div className="flex items-center mt-4">
            <ArrowUpRight className="h-4 w-4 text-green-500" />
            <span className="text-sm text-green-500 ml-1">+23</span>
            <span className="text-sm text-gray-500 ml-2">this week</span>
          </div>
        </div>
      </div>

      {/* Charts and Tables */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Price Chart */}
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold text-gray-900">Price Chart</h3>
            <div className="flex space-x-2">
              <button className="px-3 py-1 text-sm bg-blue-100 text-blue-700 rounded-md">24H</button>
              <button className="px-3 py-1 text-sm text-gray-500 hover:bg-gray-100 rounded-md">7D</button>
              <button className="px-3 py-1 text-sm text-gray-500 hover:bg-gray-100 rounded-md">1M</button>
            </div>
          </div>
          <ResponsiveContainer width="100%" height={300}>
            <LineChart data={priceData} margin={{ top: 20, right: 20, left: 0, bottom: 0 }}>
              {/* Background grid */}
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="time" />
              <YAxis domain={['auto', 'auto']} />
              <Tooltip />
              <Line type="monotone" dataKey="BTC" stroke="#3B82F6" strokeWidth={2} dot={false} />
              <Line type="monotone" dataKey="ETH" stroke="#10B981" strokeWidth={2} dot={false} />
            </LineChart>
          </ResponsiveContainer>
        </div>

        {/* Top Markets */}
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold text-gray-900">Top Markets</h3>
            <Link to="/markets" className="text-blue-600 hover:text-blue-500 text-sm">
              View all
            </Link>
          </div>
          <div className="space-y-4">
            {markets.map((market) => (
              <div key={market.id} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <div className="flex items-center">
                  <div className="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
                    <span className="text-sm font-medium text-blue-600">
                      {market.base_currency.charAt(0)}
                    </span>
                  </div>
                  <div className="ml-3">
                    <p className="text-sm font-medium text-gray-900">
                      {market.base_currency}/{market.quote_currency}
                    </p>
                    <p className="text-xs text-gray-500">{market.name}</p>
                  </div>
                </div>
                <div className="text-right">
                  <p className="text-sm font-medium text-gray-900">
                    {formatCurrency(market.current_price)}
                  </p>
                  <p className={`text-xs ${
                    parseFloat(market.price_change_24h) >= 0 ? 'text-green-500' : 'text-red-500'
                  }`}>
                    {parseFloat(market.price_change_24h) >= 0 ? '+' : ''}{market.price_change_24h}%
                  </p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Recent Activity */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mt-10">
        {/* Recent Orders */}
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold text-gray-900">Recent Orders</h3>
            <Link to="/orders" className="text-blue-600 hover:text-blue-500 text-sm">
              View all
            </Link>
          </div>
          <div className="space-y-3">
            {recentOrders.map((order) => (
              <div key={order.id} className="flex items-center justify-between p-3 border border-gray-200 rounded-lg">
                <div className="flex items-center">
                  <div className={`w-2 h-2 rounded-full ${
                    order.type === 'buy' ? 'bg-green-500' : 'bg-red-500'
                  }`} />
                  <div className="ml-3">
                    <p className="text-sm font-medium text-gray-900">
                      {order.type.toUpperCase()} {order.market_id}
                    </p>
                    <p className="text-xs text-gray-500">{formatDate(order.created_at)}</p>
                  </div>
                </div>
                <div className="text-right">
                  <p className="text-sm font-medium text-gray-900">
                    {formatCurrency(order.price)} × {order.amount}
                  </p>
                  <span className={`inline-flex px-2 py-1 text-xs font-medium rounded-full ${
                    order.status === 'completed' 
                      ? 'bg-green-100 text-green-800'
                      : 'bg-yellow-100 text-yellow-800'
                  }`}>
                    {order.status}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Recent Transactions */}
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold text-gray-900">Recent Transactions</h3>
            <Link to="/transactions" className="text-blue-600 hover:text-blue-500 text-sm">
              View all
            </Link>
          </div>
          <div className="space-y-3">
            {recentTransactions.map((transaction) => (
              <div key={transaction.id} className="flex items-center justify-between p-3 border border-gray-200 rounded-lg">
                <div className="flex items-center">
                  <div className={`w-2 h-2 rounded-full ${
                    transaction.type === 'deposit' ? 'bg-green-500' : 'bg-red-500'
                  }`} />
                  <div className="ml-3">
                    <p className="text-sm font-medium text-gray-900">
                      {transaction.type.charAt(0).toUpperCase() + transaction.type.slice(1)}
                    </p>
                    <p className="text-xs text-gray-500">{formatDate(transaction.created_at)}</p>
                  </div>
                </div>
                <div className="text-right">
                  <p className={`text-sm font-medium ${
                    transaction.type === 'deposit' ? 'text-green-600' : 'text-red-600'
                  }`}>
                    {transaction.type === 'deposit' ? '+' : '-'}{transaction.amount} {transaction.currency}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard; 