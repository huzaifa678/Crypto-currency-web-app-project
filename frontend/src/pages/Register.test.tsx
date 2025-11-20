vi.mock('../contexts/AuthContext', () => ({
  useAuth: vi.fn(),
}));

import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { vi } from 'vitest';
import Register from './Register';
import { setupMockAuth, mockRegister } from '../mock/mockAuth';


const mockNavigate = vi.fn();
vi.mock('react-router-dom', async (importOriginal) => {
  const actual = await vi.importActual<any>('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

describe('Register Component', () => {
  const mockRegister = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    setupMockAuth({
        register: mockRegister
    });
  });

  it('renders all form fields', () => {
    render(
      <MemoryRouter>
        <Register />
      </MemoryRouter>
    );

    expect(screen.getByLabelText(/username/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/email address/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/^password$/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/confirm password/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/role/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /create account/i })).toBeInTheDocument();
  });

  it('shows alert if passwords do not match', async () => {
    const alertMock = vi.spyOn(window, 'alert').mockImplementation(() => {});
    render(
      <MemoryRouter>
        <Register />
      </MemoryRouter>
    );

    fireEvent.change(screen.getByLabelText(/^Password$/i), { target: { value: '123456' } });
    fireEvent.change(screen.getByLabelText(/^Confirm Password$/i), { target: { value: 'abcdef' } });

    fireEvent.submit(screen.getByRole('button', { name: /create account/i }));

    await waitFor(() => {
      expect(alertMock).toHaveBeenCalledWith('Passwords do not match');
    });

    alertMock.mockRestore();
  });

  it('shows alert if password is too short', async () => {
    const alertMock = vi.spyOn(window, 'alert').mockImplementation(() => {});
    render(
      <MemoryRouter>
        <Register />
      </MemoryRouter>
    );

    fireEvent.change(screen.getByLabelText(/^Password$/i), { target: { value: '123' } });
    fireEvent.change(screen.getByLabelText(/^Confirm Password$/i), { target: { value: '123' } });

    fireEvent.submit(screen.getByRole('button', { name: /create account/i }));

    await waitFor(() => {
      expect(alertMock).toHaveBeenCalledWith('Password must be at least 6 characters long');
    });

    alertMock.mockRestore();
  });

  it('calls register and navigates on successful submit', async () => {
    mockRegister.mockResolvedValueOnce({});
    render(
      <MemoryRouter>
        <Register />
      </MemoryRouter>
    );

    fireEvent.change(screen.getByLabelText(/username/i), { target: { value: 'JohnDoe' } });
    fireEvent.change(screen.getByLabelText(/email address/i), { target: { value: 'john@example.com' } });
    fireEvent.change(screen.getByLabelText(/^password$/i), { target: { value: '123456' } });
    fireEvent.change(screen.getByLabelText(/confirm password/i), { target: { value: '123456' } });

    fireEvent.submit(screen.getByRole('button', { name: /create account/i }));

    await waitFor(() => {
      expect(mockRegister).toHaveBeenCalledWith('JohnDoe', 'john@example.com', '123456', 'user');
      expect(mockNavigate).toHaveBeenCalledWith('/login');
    });
  });

  it('displays loading state when submitting', async () => {
    mockRegister.mockImplementation(
      () => new Promise((resolve) => setTimeout(resolve, 100))
    );

    render(
      <MemoryRouter>
        <Register />
      </MemoryRouter>
    );

    fireEvent.change(screen.getByLabelText(/username/i), { target: { value: 'JohnDoe' } });
    fireEvent.change(screen.getByLabelText(/email address/i), { target: { value: 'john@example.com' } });
    fireEvent.change(screen.getByLabelText(/^password$/i), { target: { value: '123456' } });
    fireEvent.change(screen.getByLabelText(/confirm password/i), { target: { value: '123456' } });

    fireEvent.submit(screen.getByRole('button', { name: /create account/i }));

    expect(screen.getByText(/creating account/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /creating account/i })).toBeDisabled();

    await waitFor(() => expect(mockRegister).toHaveBeenCalled());
  });
});
