import React from 'react';
import { Outlet, NavLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { 
  Home, 
  TrendingUp, 
  BarChart3, 
  Wallet, 
  FileText, 
  User, 
  LogOut,
  Menu,
  X
} from 'lucide-react';

const Layout: React.FC = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const [sidebarOpen, setSidebarOpen] = React.useState(false);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const navigation = [
    { name: 'Dashboard', href: '/dashboard', icon: Home },
    { name: 'Markets', href: '/markets', icon: TrendingUp },
    { name: 'Orders', href: '/orders', icon: BarChart3 },
    { name: 'Wallet', href: '/wallet', icon: Wallet },
    { name: 'Transactions', href: '/transactions', icon: FileText },
    { name: 'Profile', href: '/profile', icon: User },
    { name: 'WebSocket', href: '/websocket', icon: TrendingUp },
  ];

  return (
  <div className="min-h-screen bg-gray-50">
    {/* Mobile sidebar overlay */}
    {sidebarOpen && (
      <div 
        className="inset-0 z-40 bg-gray-600 bg-opacity-75 lg:hidden"
        onClick={() => setSidebarOpen(false)}
      />
    )}

    {/* Sidebar */}
    <div
      className={`fixed top-12 inset-y-0 left-0 z-50 py-10 w-64 bg-white shadow-lg transform transition-transform duration-300 ease-in-out lg:translate-x-0 lg:static lg:inset-0 flex flex-col ${
        sidebarOpen ? 'translate-x-0' : '-translate-x-full'
      }`}
    >
      {/* Sidebar header */}
      <div className="flex items-center justify-between h-16 px-6 border-b border-gray-200">
        <button
          onClick={() => setSidebarOpen(false)}
          className="lg:hidden p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100"
        >
          <X className="h-6 w-6" />
        </button>
      </div>

      {/* Navigation */}
      <nav className="px-10 overflow-hidden flex-1 mt-10">
        <div className="space-y-1 pb-6">
          {navigation.map((item) => {
            const Icon = item.icon;
            return (
              <NavLink
                key={item.name}
                to={item.href}
                className={({ isActive }) =>
                  `group flex items-center w-full px-3 py-2 text-sm font-medium rounded-md transition-colors ${
                    isActive
                      ? 'bg-blue-100 text-blue-700 border-r-2 border-blue-700'
                      : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
                  }`
                }
                onClick={() => setSidebarOpen(false)}
              >
                <Icon className="mr-3 h-5 w-5" />
                {item.name}
              </NavLink>
            );
          })}
        </div>
      </nav>

      {/* User info and logout */}
      <div className="mt-auto p-4 border-t border-gray-200">
        <div className="flex items-center justify-between">
          <div className="flex items-center">
            <div className="h-8 w-8 bg-blue-500 rounded-full flex items-center justify-center">
              <span className="text-white text-sm font-medium">
                {user?.username?.charAt(0).toUpperCase()}
              </span>
            </div>
            <div className="ml-3">
              <p className="text-sm font-medium text-gray-700">{user?.username}</p>
              <p className="text-xs text-gray-500">{user?.email}</p>
            </div>
          </div>
          <button
            onClick={handleLogout}
            className="p-2 text-gray-400 hover:text-gray-500 hover:bg-gray-100 rounded-md"
            title="Logout"
          >
            <LogOut className="h-5 w-5" />
          </button>
        </div>
      </div>
    </div>

    {/* Fixed header */}
    <header className="fixed top-0 right-0 left-0 lg:left-10 bg-white shadow-sm border-b border-gray-200 z-40">
      <div className="flex items-center h-12">
        <button
          onClick={() => setSidebarOpen(true)}
          className="lg:hidden p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100"
        >
          <Menu className="h-6 w-6" />
        </button>

        <div className="flex-1 lg:hidden" />

        <div className="flex items-center space-x-2">
          <span className="text-sm text-gray-500 whitespace-nowrap">Welcome back,</span>
          <span className="text-sm font-medium text-gray-900 whitespace-nowrap">{user?.username}</span>
        </div>
      </div>
    </header>

    <main className="lg:ml-64 px-5 pt-14 overflow-y-auto">
      <Outlet />
    </main>
  </div>
)};

export default Layout;