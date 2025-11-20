import { render, screen, fireEvent } from '@testing-library/react';
import { vi } from 'vitest';
import Profile from '../pages/Profile';
import { setupMockAuth } from '../mock/mockAuth';

vi.mock('../contexts/AuthContext', () => ({
  useAuth: vi.fn(),
}));

describe('Profile Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupMockAuth({
      user: {
        id: 'u123',
        username: 'john_doe',
        email: 'john@example.com',
        role: 'user',
        is_verified: true,
        created_at: '2024-10-01T00:00:00Z',
        updated_at: '2025-01-01T00:00:00Z',
      },
    });
  });

  it('renders user information correctly', () => {
    render(<Profile />);

    expect(screen.getByText('Profile')).toBeInTheDocument();
    expect(screen.getByText('Account Information')).toBeInTheDocument();
    expect(screen.getByDisplayValue('john_doe')).toBeInTheDocument();
    expect(screen.getByDisplayValue('john@example.com')).toBeInTheDocument();
    expect(screen.getByText('user')).toBeInTheDocument();
    expect(screen.getByText('Verified')).toBeInTheDocument();
  });

  it('toggles edit mode when Edit is clicked', () => {
    render(<Profile />);

    fireEvent.click(screen.getByRole('button', { name: 'Edit' }));

    const usernameInput = screen.getByDisplayValue('john_doe') as HTMLInputElement;
    expect(usernameInput.disabled).toBe(false);

    expect(screen.getByRole('button', { name: /save changes/i })).toBeInTheDocument();
  });

  it('updates form fields in edit mode', () => {
    render(<Profile />);

    fireEvent.click(screen.getByRole('button', { name: 'Edit' }));

    const usernameInput = screen.getByDisplayValue('john_doe');
    const emailInput = screen.getByDisplayValue('john@example.com');

    fireEvent.change(usernameInput, { target: { value: 'jane_doe' } });
    fireEvent.change(emailInput, { target: { value: 'jane@example.com' } });

    expect(screen.getByDisplayValue('jane_doe')).toBeInTheDocument();
    expect(screen.getByDisplayValue('jane@example.com')).toBeInTheDocument();
  });

  it('closes edit mode on Cancel click', () => {
    render(<Profile />);

    fireEvent.click(screen.getByRole('button', { name: 'Edit' }));
    fireEvent.click(screen.getByTestId('cancel-edit-btn')); 

    expect(screen.getByRole('button', { name: 'Edit' })).toBeInTheDocument();
  });

  it('submits form and exits edit mode on Save', () => {
    render(<Profile />);

    fireEvent.click(screen.getByRole('button', { name: 'Edit' }));
    fireEvent.click(screen.getByRole('button', { name: /save changes/i }));

    expect(screen.getByRole('button', { name: 'Edit' })).toBeInTheDocument();
  });

  it('renders verified status correctly from client if user missing', () => {
    setupMockAuth({
      user: null,
      client: {
        username: 'client_user',
        email: 'client@example.com',
        role: 'client',
        created_at: '2023-05-12T00:00:00Z',
        is_verified: false,
        id: 'c456',
        provider: 'google',
      },
    });

    render(<Profile />);

    expect(screen.getByText('client')).toBeInTheDocument();
    expect(screen.getByText('Unverified')).toBeInTheDocument();
  });
});
