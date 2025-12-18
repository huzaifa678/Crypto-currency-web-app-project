import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../contexts/AuthContext';
import { useMarkets } from '../contexts/MarketContext';
import { Search, ArrowUpRight, Plus } from 'lucide-react';

export interface Market {
  market_id: string;
  name: string;
  base_currency: string;
  quote_currency: string;
  min_order_amount?: number;
  price_precision?: number;
  created_at?: string;
}

const Markets: React.FC = () => {
  const { market, setMarket } = useMarkets();
  const [markets, setMarkets] = useState<Market[]>([]);
  const [filteredMarkets, setFilteredMarkets] = useState<Market[]>([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [loading, setLoading] = useState(true);
  const [sortBy, setSortBy] = useState<'name' | 'base_currency' | 'quote_currency' | 'created_at'>('name');
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
        const user = localStorage.getItem('user');
        const username = user ? JSON.parse(user).username : '';
        const response = await api.get('/v1/markets', {
          params: { username: username }
        });

        const normalizedMarkets = (response.data.markets || []).filter(Boolean).map((m: any) => ({
          market_id: m.market_id,
          name: m.name,
          base_currency: m.base_currency,
          quote_currency: m.quote_currency,
          min_order_amount: parseFloat(m.min_order_amount ?? '0'),
          price_precision: parseFloat(m.price_precision ?? '2'),
          created_at: m.created_at,
        }));

        setMarket(normalizedMarkets);
        setMarkets(normalizedMarkets);
        setFilteredMarkets(normalizedMarkets);
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
  };

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
    const query = (searchTerm ?? '').toLowerCase();

    const filtered = markets.filter((market) =>
      (market.name ?? '').toLowerCase().includes(query) ||
      (market.base_currency ?? '').toLowerCase().includes(query) ||
      (market.quote_currency ?? '').toLowerCase().includes(query)
    );

    setFilteredMarkets(filtered);
  }, [searchTerm, markets]);

  const sortMarkets = (markets: Market[]) => {
    const sorted = [...markets].sort((a, b) => {
      let aValue: string | number = '';
      let bValue: string | number = '';

      switch (sortBy) {
        case 'name':
          aValue = a.name || '';
          bValue = b.name || '';
          break;
        case 'base_currency':
          aValue = a.base_currency || '';
          bValue = b.base_currency || '';
          break;
        case 'quote_currency':
          aValue = a.quote_currency || '';
          bValue = b.quote_currency || '';
          break;
        case 'created_at':
          aValue = new Date(a.created_at ?? 0).getTime();
          bValue = new Date(b.created_at ?? 0).getTime();
          break;
      }

      if (typeof aValue === 'string' && typeof bValue === 'string') {
        return sortOrder === 'asc'
          ? aValue.localeCompare(bValue)
          : bValue.localeCompare(aValue);
      } else {
        return sortOrder === 'asc'
          ? (aValue as number) - (bValue as number)
          : (bValue as number) - (aValue as number);
      }
    });

    return sorted;
  };

  const handleSort = (column: 'name' | 'base_currency' | 'quote_currency' | 'created_at') => {
    if (sortBy === column) {
      setSortOrder((prev) => (prev === 'asc' ? 'desc' : 'asc'));
    } else {
      setSortBy(column);
      setSortOrder('asc');
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64" role="status">
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
              to="/trading"
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
        </div>
      </div>

      {/* Markets Table */}
      <div className="bg-white rounded-lg shadow top-5 overflow-hidden mt-4">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th onClick={() => handleSort('base_currency')} className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100">
                  Base
                </th>
                <th onClick={() => handleSort('quote_currency')} className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100">
                  Quote
                </th>
                <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                  Min Order Amount
                </th>
                <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                  Price Precision
                </th>
                <th onClick={() => handleSort('created_at')} className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100">
                  Created At
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Action
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {sortMarkets(filteredMarkets).map((market) => (
                <tr key={market.market_id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{market.base_currency}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{market.quote_currency}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{market.min_order_amount ?? '-'}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{market.price_precision ?? '-'}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{market.created_at ? new Date(market.created_at).toLocaleString() : '-'}</td>
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
