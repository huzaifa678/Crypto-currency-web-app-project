import { render, screen } from '@testing-library/react';
import { vi } from 'vitest';
import Login from './pages/Login';
import { MemoryRouter } from 'react-router-dom';
import App from './App';

vi.mock('@react-oauth/google', async () => {
  const actual: any = await vi.importActual('@react-oauth/google');
  return {
    ...actual, // preserve other exports
    GoogleOAuthProvider: ({ children }: any) => <>{children}</>, // mock provider
    GoogleLogin: ({ onSuccess, onError }: any) => (
      <button onClick={() => onSuccess?.({ credential: 'fake-token' })}>
        Sign in with Google
      </button>
    ), // mock GoogleLogin
  };
});

vi.mock('./contexts/AuthContext', () => {
  return {
    AuthProvider: ({ children }: any) => <>{children}</>,
    useAuth: () => ({ isAuthenticated: false }),
  };
});

test('renders login page', () => {
  render(
    <MemoryRouter initialEntries={['/login']}>
      <App />
    </MemoryRouter>
  );

  expect(
    screen.getByText(/sign in to your account/i)
  ).toBeInTheDocument();

  expect(
    screen.getByRole('button', { name: 'Sign in' }) 
  ).toBeInTheDocument();

  expect(
    screen.getByText(/create a new account/i)
  ).toBeInTheDocument();
});