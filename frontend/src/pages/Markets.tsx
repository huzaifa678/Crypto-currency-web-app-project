import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../contexts/AuthContext';
import { useMarkets } from '../contexts/MarketContext';
import { Search, TrendingUp, TrendingDown, ArrowUpRight, Plus } from 'lucide-react';

export interface Market {
  market_id: string;
  name: string;
  base_currency: string;
  quote_currency: string;
  current_price: string;
  price_change_24h: string;
  volume_24h: string;
  high_24h: string;
  low_24h: string;
}

const Markets: React.FC = () => {
  const { market, setMarket } = useMarkets();
  const [markets, setMarkets] = useState<Market[]>([]);
  const [filteredMarkets, setFilteredMarkets] = useState<Market[]>([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [loading, setLoading] = useState(true);
  const [sortBy, setSortBy] = useState<'name' | 'price' | 'change' | 'volume'>('name');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc');
  const [marketId, setMarketId] = useState('');

  const [baseCurrency, setBaseCurrency] = useState('');
  const [quoteCurrency, setQuoteCurrency] = useState('');
  const [minOrderAmount, setMinOrderAmount] = useState('');
  const [pricePrecision, setPricePrecision] = useState(2);

  useEffect(() => {
    const fetchMarkets = async () => {
      try {
        setLoading(true);
        const response = await api.get('/v1/markets');
        console.log("Markets response:", response.data);
        setMarket(response.data.markets || []);
        setMarkets(response.data.markets || []);
        setFilteredMarkets(response.data.markets || []);

        console.log("market", market);
      } catch (error) {
        console.error('Error fetching markets:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchMarkets();
  }, []);

  const fetchMarketById = async (marketId: string) => {
    try {
      const response = await api.get(`/v1/markets/${marketId}`); 
      return response.data.market;
    } catch (error) {
      console.error(`Error fetching market ${marketId}:`, error);
      return null;
    }
  };

  const deleteMarket = async (marketId: string) => {
    try {
      await api.delete(`/v1/markets/${marketId}`);
      setMarkets((prev) => prev.filter((m) => m.market_id !== marketId));
      setFilteredMarkets((prev) => prev.filter((m) => m.market_id !== marketId));
    } catch (error) {
      console.error('Error deleting market:', error);
    }
  }

  const createMarket = async () => {
    try {
      const payload = {
        base_currency: baseCurrency,
        quote_currency: quoteCurrency,
        min_order_amount: minOrderAmount,
        price_precision: pricePrecision,
      };

      const response = await api.post('/v1/markets', payload);
      const newMarketId = response.data.market_id;

      setMarketId(newMarketId);

      const newMarket = await fetchMarketById(newMarketId);
      if (newMarket) {
        setMarkets((prev) => [...prev, newMarket]);
        setFilteredMarkets((prev) => [...prev, newMarket]);
      }

      setBaseCurrency('');
      setQuoteCurrency('');
      setMinOrderAmount('');
      setPricePrecision(2);

    } catch (error) {
      console.error('Error creating market:', error);
    }
  };

  useEffect(() => {
    const query = (searchTerm ?? "").toLowerCase();

    const filtered = markets.filter(market =>
      (market.name ?? "").toLowerCase().includes(query) ||
      (market.base_currency ?? "").toLowerCase().includes(query) ||
      (market.quote_currency ?? "").toLowerCase().includes(query)
    );

    setFilteredMarkets(filtered);
  }, [searchTerm, markets]);

  const sortMarkets = (markets: Market[]) => {
    return [...markets].sort((a, b) => {
      let aValue: string | number;
      let bValue: string | number;

      switch (sortBy) {
        case 'name':
          aValue = a.name;
          bValue = b.name;
          break;
        case 'price':
          aValue = parseFloat(a.current_price);
          bValue = parseFloat(b.current_price);
          break;
        case 'change':
          aValue = parseFloat(a.price_change_24h);
          bValue = parseFloat(b.price_change_24h);
          break;
        case 'volume':
          aValue = parseFloat(a.volume_24h);
          bValue = parseFloat(b.volume_24h);
          break;
        default:
          aValue = a.name;
          bValue = b.name;
      }

      if (sortOrder === 'asc') {
        return aValue > bValue ? 1 : -1;
      } else {
        return aValue < bValue ? 1 : -1;
      }
    });
  };

  const formatCurrency = (amount: string, currency: string = 'USD') => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency,
      minimumFractionDigits: 2,
      maximumFractionDigits: 8
    }).format(parseFloat(amount));
  };

  const formatVolume = (volume: string) => {
    const num = parseFloat(volume);
    if (num >= 1e9) {
      return `$${(num / 1e9).toFixed(2)}B`;
    } else if (num >= 1e6) {
      return `$${(num / 1e6).toFixed(2)}M`;
    } else if (num >= 1e3) {
      return `$${(num / 1e3).toFixed(2)}K`;
    }
    return formatCurrency(volume);
  };

  const handleSort = (column: 'name' | 'price' | 'change' | 'volume') => {
    if (sortBy === column) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      setSortBy(column);
      setSortOrder('asc');
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="fixed top-20 bottom-0 right-0 left-64 space-y-6 px-5 overflow-y-auto">
      {/* Header */}
      <div className="bg-white rounded-lg shadow p-6 sticky top-0 z-10">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Markets</h1>
            <p className="text-gray-600 mt-1">Explore and trade cryptocurrency markets</p>
          </div>
          <div className="mt-4 sm:mt-0 flex space-x-2">
            <Link
              to="/trading/btc-usd"
              className="inline-flex items-center px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-md hover:bg-blue-700 transition-colors"
            >
              <ArrowUpRight className="h-4 w-4 mr-2" />
              Start Trading
            </Link>
            <button
              onClick={createMarket}
              className="inline-flex items-center px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-md hover:bg-green-700 transition-colors"
            >
              <Plus className="h-4 w-4 mr-2" />
              Create Market
            </button>
          </div>
        </div>

        {/* Create Market Form */}
        <div className="mt-6 grid grid-cols-1 sm:grid-cols-4 gap-4">
          <input
            type="text"
            placeholder="Base Currency"
            value={baseCurrency}
            onChange={(e) => setBaseCurrency(e.target.value)}
            className="border rounded p-2"
          />
          <input
            type="text"
            placeholder="Quote Currency"
            value={quoteCurrency}
            onChange={(e) => setQuoteCurrency(e.target.value)}
            className="border rounded p-2"
          />
          <input
            type="text"
            placeholder="Min Order Amount"
            value={minOrderAmount}
            onChange={(e) => setMinOrderAmount(e.target.value)}
            className="border rounded p-2"
          />
          <input
            type="number"
            placeholder="Price Precision"
            value={pricePrecision}
            onChange={(e) => setPricePrecision(Number(e.target.value))}
            className="border rounded p-2"
          />
        </div>
      </div>

      {/* Search + Sort */}
      <div className="bg-white rounded-lg shadow p-6">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between space-y-4 sm:space-y-0">
          <div className="relative flex-1 max-w-md">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <Search className="h-5 w-5 text-gray-400" />
            </div>
            <input
              type="text"
              placeholder="Search markets..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md leading-5 bg-white placeholder-gray-500 focus:outline-none focus:placeholder-gray-400 focus:ring-1 focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
            />
          </div>
          <div className="flex items-center space-x-4">
            <select
              value={`${sortBy}-${sortOrder}`}
              onChange={(e) => {
                const [column, order] = e.target.value.split('-');
                setSortBy(column as any);
                setSortOrder(order as any);
              }}
              className="block w-full sm:w-auto px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
            >
              <option value="name-asc">Name (A-Z)</option>
              <option value="name-desc">Name (Z-A)</option>
              <option value="price-asc">Price (Low-High)</option>
              <option value="price-desc">Price (High-Low)</option>
              <option value="change-asc">Change (Low-High)</option>
              <option value="change-desc">Change (High-Low)</option>
              <option value="volume-asc">Volume (Low-High)</option>
              <option value="volume-desc">Volume (High-Low)</option>
            </select>
          </div>
        </div>
      </div>

      {/* Markets Table */}
      <div className="bg-white rounded-lg shadow top-5 overflow-hidden mt-4">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th onClick={() => handleSort('name')} className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100">Market</th>
                <th onClick={() => handleSort('price')} className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100">Price</th>
                <th onClick={() => handleSort('change')} className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100">24h Change</th>
                <th onClick={() => handleSort('volume')} className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100">24h Volume</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">24h High</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">24h Low</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Action</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {sortMarkets(filteredMarkets).map((market) => (
                <tr key={market.market_id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center">
                      <div className="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
                        <span className="text-sm font-medium text-blue-600">{market.base_currency.charAt(0)}</span>
                      </div>
                      <div className="ml-4">
                        <div className="text-sm font-medium text-gray-900">{market.name}</div>
                        <div className="text-sm text-gray-500">{market.base_currency}/{market.quote_currency}</div>
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">{formatCurrency(market.current_price)}</td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className={`flex items-center text-sm ${parseFloat(market.price_change_24h) >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                      {parseFloat(market.price_change_24h) >= 0 ? <TrendingUp className="h-4 w-4 mr-1" /> : <TrendingDown className="h-4 w-4 mr-1" />}
                      {parseFloat(market.price_change_24h) >= 0 ? '+' : ''}{market.price_change_24h}%
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{formatVolume(market.volume_24h)}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{formatCurrency(market.high_24h)}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{formatCurrency(market.low_24h)}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    <button
                      className="text-red-600 hover:text-red-900 font-medium"
                      onClick={async (e) => {
                        e.preventDefault();
                        await deleteMarket(market.market_id);
                      }}
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        {filteredMarkets.length === 0 && (
          <div className="text-center py-12">
            <p className="text-gray-500">No markets found matching your search.</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default Markets;