import React, { useState, useEffect } from 'react';
import { api } from '../contexts/AuthContext';
import { Clock, CheckCircle, XCircle, Trash2, PlusCircle } from 'lucide-react';
import toast from 'react-hot-toast';
import { useMarkets } from '../contexts/MarketContext';
import { useOrder } from '../contexts/OrderContext';

export interface Order {
  id: string;
  market_id: string;
  type: 'BUY' | 'SELL';
  status: 'OPEN' | 'PARTIALLY_FILLED' | 'FILLED' | 'CANCELLED';
  price: string;
  amount: string;
  filled_amount?: string;
  created_at: string;
}

const Orders: React.FC = () => {
  const { setOrder } = useOrder();
  const { market } = useMarkets();
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);

  const [creating, setCreating] = useState(false);
  const [type, setType] = useState<'BUY' | 'SELL'>('BUY');
  const [price, setPrice] = useState('');
  const [amount, setAmount] = useState('');

  const [fetchedOrder, setFetchedOrder] = useState<Order | null>(null);

  useEffect(() => {
    const fetchOrders = async () => {
      try {
        setLoading(true);
        const response = await api.get('/v1/orders');
        setOrders(response.data.orders ?? []);
        console.log('Fetched orders:', response.data.orders);
        console.log('data', response.data);
        console.log("Orders API response (stringified):", JSON.stringify(response.data, null, 2));

      } finally {
        setLoading(false);
      }
    }

    fetchOrders();
  }, []);

  const handleCreateOrder = async () => {
    console.log("markets", market);
    console.log("markets (stringified)", JSON.stringify(market));
    console.log("Marketid ", market[0]);
    try {
      setCreating(true);


      const response = await api.post<{ order_id: string }>('/v1/orders', {
        user_email: localStorage.getItem('user')
          ? JSON.parse(localStorage.getItem('user')!).email
          : '',
        market_id: market[0].market_id,
        type,
        price,
        amount,
      });

      const orderId = response.data.order_id;

      const orderResponse = await api.get<{ order: Order }>(`/v1/orders/${orderId}`);
      const createdOrder = orderResponse.data.order;

      setOrder(createdOrder);

      setFetchedOrder(createdOrder);
      setOrders((prev) => [...prev, createdOrder]);

      toast.success('Order created & fetched successfully!');
      setPrice('');
      setAmount('');
    } catch (error: any) {
      const message = error.response?.data?.error || 'Failed to create order';
      console.log(error);
      toast.error(message);
    } finally {
      setCreating(false);
    }
  };

  const handleCancelOrder = async (orderId: string) => {
    try {
      await api.delete(`/v1/orders/${orderId}`);
      setOrders(orders.filter(order => order.id !== orderId));
      if (fetchedOrder?.id === orderId) setFetchedOrder(null);
      toast.success('Order cancelled successfully!');
    } catch (error: any) {
      toast.error('Failed to cancel order');
    }
  };

  const formatCurrency = (amount: string) =>
    new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(parseFloat(amount));

  const formatDate = (dateString: string) =>
    new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric', month: 'short', day: 'numeric',
      hour: '2-digit', minute: '2-digit'
    });

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'FILLED': return <CheckCircle className="h-5 w-5 text-green-500" />;
      case 'OPEN': return <Clock className="h-5 w-5 text-yellow-500" />;
      case 'PARTIALLY_FILLED': return <Clock className="h-5 w-5 text-blue-500" />;
      case 'CANCELLED': return <XCircle className="h-5 w-5 text-red-500" />;
      default: return <Clock className="h-5 w-5 text-gray-500" />;
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
    <div className="fixed top-10 bottom-0 right-0 left-[16rem] z-30 bg-white shadow p-6 overflow-y-auto">
      {/* Header + Create Order */}
      <div className="bg-white rounded-lg shadow p-6 sticky z-20 mb-10">
        <h1 className="text-2xl font-bold text-gray-900">Orders</h1>
        <p className="text-gray-600 mt-1">Manage your trading orders</p>

        {/* Create Order Form */}
        <div className="mt-4 grid grid-cols-2 md:grid-cols-6 gap-3">
          <select
            value={type}
            onChange={(e) => setType(e.target.value as 'BUY' | 'SELL')}
            className="border rounded-md px-3 py-2 text-sm"
          >
            <option value="BUY">BUY</option>
            <option value="SELL">SELL</option>
          </select>
          <input
            type="number"
            value={price}
            onChange={(e) => setPrice(e.target.value)}
            placeholder="Price"
            className="border rounded-md px-3 py-2 text-sm"
          />
          <input
            type="number"
            value={amount}
            onChange={(e) => setAmount(e.target.value)}
            placeholder="Amount"
            className="border rounded-md px-3 py-2 text-sm"
          />
          <button
            onClick={handleCreateOrder}
            disabled={creating}
            className="px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-md hover:bg-blue-700 transition disabled:opacity-50 flex items-center justify-center"
          >
            <PlusCircle className="h-4 w-4 mr-2" />
            {creating ? 'Creating...' : 'Create Order'}
          </button>
        </div>
      </div>

      {/* Highlighted Fetched Order */}
      {fetchedOrder && (
        <div className="bg-yellow-50 border border-yellow-200 rounded-lg shadow p-4">
          <h2 className="text-lg font-semibold text-yellow-800 mb-2">Last Created Order</h2>
          <div className="overflow-x-auto">
            <table className="min-w-full border">
              <tbody>
                <tr className="bg-white">
                  <td className="px-6 py-4 font-medium">#{fetchedOrder.id}</td>
                  <td className="px-6 py-4">{fetchedOrder.market_id.toUpperCase()}</td>
                  <td className="px-6 py-4">
                    <span className={`px-2 py-1 text-xs rounded-full ${
                      fetchedOrder.type === 'BUY' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                    }`}>
                      {fetchedOrder.type}
                    </span>
                  </td>
                  <td className="px-6 py-4 flex items-center space-x-2">
                    {getStatusIcon(fetchedOrder.status)}
                    <span>{fetchedOrder.status}</span>
                  </td>
                  <td className="px-6 py-4">{formatCurrency(fetchedOrder.price)}</td>
                  <td className="px-6 py-4">{fetchedOrder.amount}</td>
                  <td className="px-6 py-4">{formatDate(fetchedOrder.created_at)}</td>
                  <td className="px-6 py-4">
                    {fetchedOrder.status === 'OPEN' && (
                      <button
                        onClick={() => handleCancelOrder(fetchedOrder.id)}
                        className="text-red-600 hover:text-red-900"
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    )}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* Orders Table */}
      <div className="bg-white rounded-lg shadow overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              {['Order ID','Market','Type','Status','Price','Amount','Date','Actions'].map(h => (
                <th key={h} className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{h}</th>
              ))}
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {orders.map((order) => (
              <tr key={order.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">#{order.id}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{order.market_id.toUpperCase()}</td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={`inline-flex px-2 py-1 text-xs font-medium rounded-full ${
                    order.type === 'BUY' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                  }`}>
                    {order.type}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="flex items-center">
                    {getStatusIcon(order.status)}
                    <span className="ml-2 text-sm text-gray-900">{order.status}</span>
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{formatCurrency(order.price)}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{order.amount}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{formatDate(order.created_at)}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  {order.status === 'OPEN' && (
                    <button
                      onClick={() => handleCancelOrder(order.id)}
                      className="text-red-600 hover:text-red-900"
                    >
                      <Trash2 className="h-4 w-4" />
                    </button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        {orders.length === 0 && (
          <div className="text-center py-12">
            <p className="text-gray-500">No orders found.</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default Orders;