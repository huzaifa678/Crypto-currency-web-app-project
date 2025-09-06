import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { api } from '../contexts/AuthContext';
import { 
  TrendingUp, 
  TrendingDown, 
  ArrowUpRight, 
  ArrowDownRight,
  Clock,
  DollarSign
} from 'lucide-react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';

interface OrderBookEntry {
  price: string;
  amount: string;
  total: string;
}

interface Trade {
  id: string;
  price: string;
  amount: string;
  side: 'buy' | 'sell';
  timestamp: string;
}

const Trading: React.FC = () => {
  const { marketId } = useParams<{ marketId: string }>();
  const { user } = useAuth();
  const [activeTab, setActiveTab] = useState<'market' | 'limit'>('market');
  const [orderType, setOrderType] = useState<'buy' | 'sell'>('buy');
  const [amount, setAmount] = useState('');
  const [price, setPrice] = useState('');
  const [loading, setLoading] = useState(false);
  
  const [orderBook, setOrderBook] = useState<{
    bids: OrderBookEntry[];
    asks: OrderBookEntry[];
  }>({ bids: [], asks: [] });
  
  const [recentTrades, setRecentTrades] = useState<Trade[]>([]);
  const [currentPrice, setCurrentPrice] = useState('48500.00');
  const [priceChange, setPriceChange] = useState('2.5');

  const priceData = [
    { time: '00:00', price: 45000 },
    { time: '04:00', price: 46000 },
    { time: '08:00', price: 47000 },
    { time: '12:00', price: 46500 },
    { time: '16:00', price: 48000 },
    { time: '20:00', price: 47500 },
    { time: '24:00', price: 48500 },
  ];

  useEffect(() => {
    const fetchTradingData = async () => {
      try {
        setOrderBook({
          bids: [
            { price: '48450.00', amount: '0.5', total: '24225.00' },
            { price: '48400.00', amount: '1.2', total: '58080.00' },
            { price: '48350.00', amount: '0.8', total: '38680.00' },
            { price: '48300.00', amount: '2.1', total: '101430.00' },
            { price: '48250.00', amount: '1.5', total: '72375.00' },
          ],
          asks: [
            { price: '48550.00', amount: '0.3', total: '14565.00' },
            { price: '48600.00', amount: '0.7', total: '34020.00' },
            { price: '48650.00', amount: '1.1', total: '53515.00' },
            { price: '48700.00', amount: '0.9', total: '43830.00' },
            { price: '48750.00', amount: '1.4', total: '68250.00' },
          ]
        });

        setRecentTrades([
          { id: '1', price: '48500.00', amount: '0.1', side: 'buy', timestamp: '12:30:45' },
          { id: '2', price: '48495.00', amount: '0.2', side: 'sell', timestamp: '12:30:42' },
          { id: '3', price: '48500.00', amount: '0.05', side: 'buy', timestamp: '12:30:38' },
          { id: '4', price: '48490.00', amount: '0.3', side: 'sell', timestamp: '12:30:35' },
          { id: '5', price: '48495.00', amount: '0.15', side: 'buy', timestamp: '12:30:32' },
        ]);
      } catch (error) {
        console.error('Error fetching trading data:', error);
      }
    };

    fetchTradingData();
  }, [marketId]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      const orderData = {
        user_email: user?.email,
        market_id: marketId,
        type: orderType,
        status: 'pending',
        price: activeTab === 'limit' ? price : currentPrice,
        amount: amount
      };

      await api.post('/orders', orderData);
      alert('Order placed successfully!');
      
      // Reset form
      setAmount('');
      setPrice('');
    } catch (error) {
      console.error('Error placing order:', error);
      alert('Failed to place order. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const formatCurrency = (amount: string) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 8
    }).format(parseFloat(amount));
  };

  const formatAmount = (amount: string) => {
    return parseFloat(amount).toFixed(8);
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="bg-white rounded-lg shadow p-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">
              {marketId?.toUpperCase()} Trading
            </h1>
            <div className="flex items-center mt-2 space-x-4">
              <span className="text-3xl font-bold text-gray-900">
                {formatCurrency(currentPrice)}
              </span>
              <div className={`flex items-center text-sm ${
                parseFloat(priceChange) >= 0 ? 'text-green-600' : 'text-red-600'
              }`}>
                {parseFloat(priceChange) >= 0 ? (
                  <TrendingUp className="h-4 w-4 mr-1" />
                ) : (
                  <TrendingDown className="h-4 w-4 mr-1" />
                )}
                {parseFloat(priceChange) >= 0 ? '+' : ''}{priceChange}%
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Price Chart */}
        <div className="lg:col-span-2 bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Price Chart</h3>
          <ResponsiveContainer width="100%" height={400}>
            <LineChart data={priceData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="time" />
              <YAxis />
              <Tooltip formatter={(value) => formatCurrency(value as string)} />
              <Line type="monotone" dataKey="price" stroke="#3B82F6" strokeWidth={2} />
            </LineChart>
          </ResponsiveContainer>
        </div>

        {/* Trading Form */}
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Place Order</h3>
          
          {/* Order Type Tabs */}
          <div className="flex mb-4 bg-gray-100 rounded-lg p-1">
            <button
              onClick={() => setActiveTab('market')}
              className={`flex-1 py-2 px-4 text-sm font-medium rounded-md transition-colors ${
                activeTab === 'market'
                  ? 'bg-white text-blue-600 shadow-sm'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Market
            </button>
            <button
              onClick={() => setActiveTab('limit')}
              className={`flex-1 py-2 px-4 text-sm font-medium rounded-md transition-colors ${
                activeTab === 'limit'
                  ? 'bg-white text-blue-600 shadow-sm'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Limit
            </button>
          </div>

          {/* Buy/Sell Tabs */}
          <div className="flex mb-4">
            <button
              onClick={() => setOrderType('buy')}
              className={`flex-1 py-2 px-4 text-sm font-medium rounded-l-md border transition-colors ${
                orderType === 'buy'
                  ? 'bg-green-50 border-green-200 text-green-700'
                  : 'bg-gray-50 border-gray-200 text-gray-600 hover:bg-gray-100'
              }`}
            >
              Buy
            </button>
            <button
              onClick={() => setOrderType('sell')}
              className={`flex-1 py-2 px-4 text-sm font-medium rounded-r-md border border-l-0 transition-colors ${
                orderType === 'sell'
                  ? 'bg-red-50 border-red-200 text-red-700'
                  : 'bg-gray-50 border-gray-200 text-gray-600 hover:bg-gray-100'
              }`}
            >
              Sell
            </button>
          </div>

          <form onSubmit={handleSubmit} className="space-y-4">
            {activeTab === 'limit' && (
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Price
                </label>
                <input
                  type="number"
                  step="0.01"
                  value={price}
                  onChange={(e) => setPrice(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
                  placeholder="Enter price"
                  required
                />
              </div>
            )}

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Amount
              </label>
              <input
                type="number"
                step="0.00000001"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
                placeholder="Enter amount"
                required
              />
            </div>

            {amount && (
              <div className="bg-gray-50 p-3 rounded-md">
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Total:</span>
                  <span className="font-medium">
                    {formatCurrency(activeTab === 'limit' && price ? (parseFloat(price) * parseFloat(amount)).toString() : (parseFloat(currentPrice) * parseFloat(amount)).toString())}
                  </span>
                </div>
              </div>
            )}

            <button
              type="submit"
              disabled={loading || !amount || (activeTab === 'limit' && !price)}
              className={`w-full py-2 px-4 text-sm font-medium rounded-md transition-colors ${
                orderType === 'buy'
                  ? 'bg-green-600 hover:bg-green-700 text-white disabled:bg-green-300'
                  : 'bg-red-600 hover:bg-red-700 text-white disabled:bg-red-300'
              } disabled:cursor-not-allowed`}
            >
              {loading ? (
                <div className="flex items-center justify-center">
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                  Placing Order...
                </div>
              ) : (
                `${orderType === 'buy' ? 'Buy' : 'Sell'} ${marketId?.split('-')[0].toUpperCase()}`
              )}
            </button>
          </form>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Order Book */}
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Order Book</h3>
          <div className="grid grid-cols-3 gap-4 text-xs font-medium text-gray-500 mb-2">
            <span>Price</span>
            <span>Amount</span>
            <span>Total</span>
          </div>
          
          {/* Asks (Sell Orders) */}
          <div className="space-y-1 mb-4">
            {orderBook.asks.map((ask, index) => (
              <div key={`ask-${index}`} className="grid grid-cols-3 gap-4 text-sm">
                <span className="text-red-600">{formatCurrency(ask.price)}</span>
                <span>{formatAmount(ask.amount)}</span>
                <span>{formatCurrency(ask.total)}</span>
              </div>
            ))}
          </div>

          {/* Current Price */}
          <div className="border-t border-b border-gray-200 py-2 mb-4">
            <div className="text-center">
              <span className="text-lg font-bold text-gray-900">{formatCurrency(currentPrice)}</span>
            </div>
          </div>

          {/* Bids (Buy Orders) */}
          <div className="space-y-1">
            {orderBook.bids.map((bid, index) => (
              <div key={`bid-${index}`} className="grid grid-cols-3 gap-4 text-sm">
                <span className="text-green-600">{formatCurrency(bid.price)}</span>
                <span>{formatAmount(bid.amount)}</span>
                <span>{formatCurrency(bid.total)}</span>
              </div>
            ))}
          </div>
        </div>

        {/* Recent Trades */}
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Recent Trades</h3>
          <div className="grid grid-cols-4 gap-4 text-xs font-medium text-gray-500 mb-2">
            <span>Price</span>
            <span>Amount</span>
            <span>Side</span>
            <span>Time</span>
          </div>
          
          <div className="space-y-1">
            {recentTrades.map((trade) => (
              <div key={trade.id} className="grid grid-cols-4 gap-4 text-sm">
                <span className={trade.side === 'buy' ? 'text-green-600' : 'text-red-600'}>
                  {formatCurrency(trade.price)}
                </span>
                <span>{formatAmount(trade.amount)}</span>
                <span className={trade.side === 'buy' ? 'text-green-600' : 'text-red-600'}>
                  {trade.side.toUpperCase()}
                </span>
                <span className="text-gray-500">{trade.timestamp}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Trading; 