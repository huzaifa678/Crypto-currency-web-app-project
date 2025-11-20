/// <reference types="vitest" />
import { vi, beforeEach, describe, it, expect } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'



vi.mock('axios', () => {
    const mockPost = vi.fn()

    const mockAxiosInstance = {
        post: mockPost,
        interceptors: {
            request: { use: vi.fn() },
            response: { use: vi.fn() },
        },
    }

    return {
        __esModule: true,
        default: {
        create: vi.fn(() => mockAxiosInstance),
        __mockPost: mockPost,
        },
    }
})

vi.mock('react-hot-toast', () => ({
  __esModule: true,
  default: {
    success: vi.fn(),
    error: vi.fn(),
  },
}))

import { AuthProvider, useAuth } from './AuthContext'

beforeEach(() => {
  localStorage.clear()
  vi.clearAllMocks()
  vi.stubGlobal('location', { href: '' })
})

const TestComponent = () => {
  const { user, client, isAuthenticated, login, logout, loading } = useAuth()
  return (
    <div>
      <div data-testid="user">{user?.username || 'no-user'}</div>
      <div data-testid="client">{client?.username || 'no-client'}</div>
      <div data-testid="auth">{isAuthenticated ? 'yes' : 'no'}</div>
      <div data-testid="loading">{loading ? 'loading' : 'idle'}</div>
      <button onClick={() => login('test@example.com', 'password')}>Login</button>
      <button onClick={logout}>Logout</button>
    </div>
  )
}

import axios from 'axios'

describe('AuthProvider', () => {
  it('initializes with no user or client', async () => {
    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    )

    expect(screen.getByTestId('user').textContent).toBe('no-user')
    expect(screen.getByTestId('client').textContent).toBe('no-client')
    expect(screen.getByTestId('auth').textContent).toBe('no')
  })

  it('login sets user and localStorage', async () => {

    const axiosMock = axios as unknown as { __mockPost: ReturnType<typeof vi.fn> }

    const mockUser = {
      id: '1',
      username: 'John',
      email: 'john@example.com',
      role: 'USER_ROLE_USER',
      is_verified: true,
      created_at: '2025-01-01',
      updated_at: '2025-01-01',
    }

    axiosMock.__mockPost.mockResolvedValueOnce({
        data: {
            access_token: 'fake-token',
            refresh_token: 'fake-refresh',
            user: mockUser,
        },
  })

    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    )

    fireEvent.click(screen.getByText('Login'))

    await waitFor(() => {
      expect(screen.getByTestId('user').textContent).toBe('John')
      expect(screen.getByTestId('auth').textContent).toBe('yes')
    })

    expect(localStorage.getItem('access_token')).toBe('fake-token')
    expect(localStorage.getItem('refresh_token')).toBe('fake-refresh')
    expect(localStorage.getItem('user')).toContain('John')
  })

  it('logout clears user, client, and localStorage', async () => {
    localStorage.setItem('user', JSON.stringify({ username: 'John' }))
    localStorage.setItem('access_token', 'token')
    localStorage.setItem('refresh_token', 'refresh')

    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    )

    fireEvent.click(screen.getByText('Logout'))

    await waitFor(() => {
      expect(screen.getByTestId('user').textContent).toBe('no-user')
      expect(screen.getByTestId('auth').textContent).toBe('no')
    })

    expect(localStorage.getItem('user')).toBeNull()
    expect(localStorage.getItem('access_token')).toBeNull()
    expect(localStorage.getItem('refresh_token')).toBeNull()
  })
})

