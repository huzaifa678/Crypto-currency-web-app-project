import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../contexts/AuthContext';
import { Search, TrendingUp, TrendingDown, ArrowUpRight } from 'lucide-react';

interface Market {
  id: string;
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
  const [markets, setMarkets] = useState<Market[]>([]);
  const [filteredMarkets, setFilteredMarkets] = useState<Market[]>([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [loading, setLoading] = useState(true);
  const [sortBy, setSortBy] = useState<'name' | 'price' | 'change' | 'volume'>('name');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc');

  useEffect(() => {
    const fetchMarkets = async () => {
      try {
        setLoading(true);
        const response = await api.get('/markets');
        setMarkets(response.data);
        setFilteredMarkets(response.data);
      } catch (error) {
        console.error('Error fetching markets:', error);
        const mockMarkets: Market[] = [
          {
            id: '1',
            name: 'Bitcoin',
            base_currency: 'BTC',
            quote_currency: 'USD',
            current_price: '48500.00',
            price_change_24h: '2.5',
            volume_24h: '2500000000',
            high_24h: '49000.00',
            low_24h: '48000.00'
          },
          {
            id: '2',
            name: 'Ethereum',
            base_currency: 'ETH',
            quote_currency: 'USD',
            current_price: '3600.00',
            price_change_24h: '-1.2',
            volume_24h: '1500000000',
            high_24h: '3650.00',
            low_24h: '3550.00'
          },
          {
            id: '3',
            name: 'Cardano',
            base_currency: 'ADA',
            quote_currency: 'USD',
            current_price: '1.25',
            price_change_24h: '5.8',
            volume_24h: '500000000',
            high_24h: '1.30',
            low_24h: '1.20'
          },
          {
            id: '4',
            name: 'Solana',
            base_currency: 'SOL',
            quote_currency: 'USD',
            current_price: '120.00',
            price_change_24h: '-0.8',
            volume_24h: '800000000',
            high_24h: '125.00',
            low_24h: '118.00'
          },
          {
            id: '5',
            name: 'Polkadot',
            base_currency: 'DOT',
            quote_currency: 'USD',
            current_price: '25.50',
            price_change_24h: '3.2',
            volume_24h: '300000000',
            high_24h: '26.00',
            low_24h: '24.80'
          }
        ];
        setMarkets(mockMarkets);
        setFilteredMarkets(mockMarkets);
      } finally {
        setLoading(false);
      }
    };

    fetchMarkets();
  }, []);

  useEffect(() => {
    const filtered = markets.filter(market =>
      market.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      market.base_currency.toLowerCase().includes(searchTerm.toLowerCase()) ||
      market.quote_currency.toLowerCase().includes(searchTerm.toLowerCase())
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
    <div className="space-y-6">
      {/* Header */}
      <div className="bg-white rounded-lg shadow p-6">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Markets</h1>
            <p className="text-gray-600 mt-1">Explore and trade cryptocurrency markets</p>
          </div>
          <div className="mt-4 sm:mt-0">
            <Link
              to="/trading/btc-usd"
              className="inline-flex items-center px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-md hover:bg-blue-700 transition-colors"
            >
              <ArrowUpRight className="h-4 w-4 mr-2" />
              Start Trading
            </Link>
          </div>
        </div>
      </div>

      {/* Search and Filters */}
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
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100" onClick={() => handleSort('name')}>
                  Market
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100" onClick={() => handleSort('price')}>
                  Price
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100" onClick={() => handleSort('change')}>
                  24h Change
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100" onClick={() => handleSort('volume')}>
                  24h Volume
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  24h High
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  24h Low
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Action
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {sortMarkets(filteredMarkets).map((market) => (
                <tr key={market.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center">
                      <div className="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
                        <span className="text-sm font-medium text-blue-600">
                          {market.base_currency.charAt(0)}
                        </span>
                      </div>
                      <div className="ml-4">
                        <div className="text-sm font-medium text-gray-900">{market.name}</div>
                        <div className="text-sm text-gray-500">{market.base_currency}/{market.quote_currency}</div>
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm font-medium text-gray-900">
                      {formatCurrency(market.current_price)}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className={`flex items-center text-sm ${
                      parseFloat(market.price_change_24h) >= 0 ? 'text-green-600' : 'text-red-600'
                    }`}>
                      {parseFloat(market.price_change_24h) >= 0 ? (
                        <TrendingUp className="h-4 w-4 mr-1" />
                      ) : (
                        <TrendingDown className="h-4 w-4 mr-1" />
                      )}
                      {parseFloat(market.price_change_24h) >= 0 ? '+' : ''}{market.price_change_24h}%
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {formatVolume(market.volume_24h)}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {formatCurrency(market.high_24h)}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {formatCurrency(market.low_24h)}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    <Link
                      to={`/trading/${market.base_currency.toLowerCase()}-${market.quote_currency.toLowerCase()}`}
                      className="text-blue-600 hover:text-blue-900 font-medium"
                    >
                      Trade
                    </Link>
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