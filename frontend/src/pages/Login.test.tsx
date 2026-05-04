vi.mock('../contexts/AuthContext', () => ({
  useAuth: vi.fn(),
}));

import { render, screen, fireEvent, waitFor, act } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { vi } from 'vitest';
import Login from './Login';
import { setupMockAuth, mockLogin, mockLoginWithGoogle } from '../mock/mockAuth';
import userEvent from '@testing-library/user-event';


const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<any>('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

vi.mock('@react-oauth/google', () => ({
  GoogleLogin: ({ onSuccess }: any) => (
    <button onClick={() => onSuccess({ credential: 'fake-google-token' })}>
      Mock Google Login
    </button>
  ),
}));

describe('Login Component', () => {

  beforeEach(() => {
    vi.clearAllMocks();
    setupMockAuth({
        login: mockLogin,
        loginWithGoogle: mockLoginWithGoogle,
    });
  });

  function renderLogin() {
    return render(
      <MemoryRouter>
        <Login />
      </MemoryRouter>
    );
  }

  test('renders login form elements', () => {
    renderLogin();

    expect(screen.getByText(/sign in to your account/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/email address/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /sign in/i })).toBeInTheDocument();
  });

  test('updates email and password inputs', () => {
    renderLogin();

    const emailInput = screen.getByLabelText(/email address/i) as HTMLInputElement;
    const passwordInput = screen.getByLabelText(/password/i) as HTMLInputElement;

    fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
    fireEvent.change(passwordInput, { target: { value: 'secret' } });

    expect(emailInput.value).toBe('test@example.com');
    expect(passwordInput.value).toBe('secret');
  });

  test('toggles password visibility', () => {
    renderLogin();

    const passwordInput = screen.getByLabelText(/password/i);
    const toggleButton = screen.getByRole('button', { name: '' }); // Eye icon button

    expect(passwordInput).toHaveAttribute('type', 'password');
    fireEvent.click(toggleButton);
    expect(passwordInput).toHaveAttribute('type', 'text');
    fireEvent.click(toggleButton);
    expect(passwordInput).toHaveAttribute('type', 'password');
  });

  test('calls login and navigates on successful form submit', async () => {
    mockLogin.mockResolvedValueOnce({});
    renderLogin();

    const emailInput = screen.getByLabelText(/email address/i);
    const passwordInput = screen.getByLabelText(/password/i);
    const submitButton = screen.getByRole('button', { name: /sign in/i });

    fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
    fireEvent.change(passwordInput, { target: { value: 'password123' } });
    fireEvent.click(submitButton);

    await waitFor(() => expect(mockLogin).toHaveBeenCalledWith('test@example.com', 'password123'));
    expect(mockNavigate).toHaveBeenCalledWith('/dashboard');
  });

  test('shows loading state during login', async () => {
    let resolveLogin: Function;
    mockLogin.mockReturnValue(
      new Promise((resolve) => {
        resolveLogin = resolve;
      })
    );

    renderLogin();

    const submitButton = screen.getByTestId('submit-button');

    const user = userEvent.setup();

    await act(async () => {
      await fireEvent.submit(submitButton.closest('form')!);
    });

    await waitFor(() => {
        expect(submitButton).toBeDisabled();
        expect(screen.getByText(/signing in/i)).toBeInTheDocument();
    });

    await act(() => resolveLogin!({}));

    await waitFor(() => expect(submitButton).not.toBeDisabled());
  });

  test('handles Google login success', async () => {
    mockLoginWithGoogle.mockResolvedValueOnce(true);
    renderLogin();

    fireEvent.click(screen.getByText(/mock google login/i));

    await waitFor(() => expect(mockLoginWithGoogle).toHaveBeenCalledWith('fake-google-token'));
    expect(mockNavigate).toHaveBeenCalledWith('/dashboard');
  });
});
