import { vi } from 'vitest';
import { useAuth } from '../contexts/AuthContext';
import type { AuthContextType } from '../contexts/AuthContext';


export const mockLogin = vi.fn();
export const mockLoginWithGoogle = vi.fn();
export const mockRegister = vi.fn();
export const mockLogout = vi.fn();
export const user = {}

export function setupMockAuth(overrides: Partial<AuthContextType> = {}) {
  const defaultMock: AuthContextType = {
    user: null,
    client: null,
    isAuthenticated: false,
    login: mockLogin,
    loginWithGoogle: mockLoginWithGoogle,
    register: mockRegister,
    logout: mockLogout,
    loading: false,
    ...overrides,
  };

  vi.mocked(useAuth).mockReturnValue(defaultMock);
}
