import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { Toaster } from 'react-hot-toast';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import { MarketsProvider } from './contexts/MarketContext';
import Layout from './components/Layout';
import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import Markets from './pages/Markets';
import Trading from './pages/Trading';
import Orders from './pages/Orders';
import Wallet from './pages/Wallet';
import Transactions from './pages/Transactions';
import Profile from './pages/Profile';
import MarketsTable from './pages/WebSocket';
import { OrderProvider } from './contexts/OrderContext';
import "./index.css";
import { GoogleOAuthProvider } from '@react-oauth/google';

const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated } = useAuth();
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />;
};

const PublicRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated } = useAuth();
  return isAuthenticated ? <Navigate to="/dashboard" replace /> : <>{children}</>;
};

const AdminRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated, user } = useAuth();

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (user?.role !== 'admin') {
    return <Navigate to="/dashboard" replace />;
  }

  return <>{children}</>;
};

function App() {
  return (
    <GoogleOAuthProvider clientId={import.meta.env.VITE_GOOGLE_CLIENT_ID || ''}>
      <AuthProvider>
          <div className="App">
            <Toaster 
              position="top-right"
              toastOptions={{
                duration: 4000,
                style: {
                  background: '#363636',
                  color: '#fff',
                },
              }}
            />
            <Routes>
              <Route path="/login" element={
                <PublicRoute>
                  <Login />
                </PublicRoute>
              } />
              <Route path="/register" element={
                <PublicRoute>
                  <Register />
                </PublicRoute>
              } />
              <Route path="/" element={
                <ProtectedRoute>
                  <Layout />
                </ProtectedRoute>
              }>
                <Route index element={<Navigate to="/dashboard" replace />} />
                <Route path="dashboard" element={<Dashboard />} />
                <Route path="websocket" element={<MarketsTable />} />
                <Route path="markets" element={
                  <AdminRoute>
                    <MarketsProvider> 
                      <Markets /> 
                    </MarketsProvider>
                  </AdminRoute>
                } />
                <Route path="trading" element={<OrderProvider> <Trading /> </OrderProvider>} />
                <Route path="orders" element={
                  <MarketsProvider>
                  <OrderProvider> 
                    <Orders /> 
                  </OrderProvider>
                  </MarketsProvider>} />
                <Route path="wallet" element={<Wallet />} />
                <Route path="transactions" element={<Transactions />} />
                <Route path="profile" element={<Profile />} />
              </Route>
            </Routes>
          </div>
      </AuthProvider>
    </GoogleOAuthProvider>
  );
}

export default App;