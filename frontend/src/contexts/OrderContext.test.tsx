/// <reference types="vitest" />
import { describe, it, expect } from 'vitest'
import { renderHook, act, render, screen } from '@testing-library/react'
import { OrderProvider, useOrder } from './OrderContext'

const TestComponent = () => {
  const { order, orders, setOrder, setOrders } = useOrder()

  return (
    <div>
      <div data-testid="order">{order ? order.id : 'no-order'}</div>
      <div data-testid="orders-count">{orders ? orders.length : 0}</div>
      <button
        onClick={() =>
          setOrder({ id: '1', market_id: '5', type: 'BUY', status: 'OPEN', price: '50000', amount: '0.1', created_at: '2024-01-01T00:00:00Z' })
        }
      >
        Set Order
      </button>
      <button
        onClick={() =>
          setOrders([{ id: '2', market_id: '10', type: 'SELL', status: 'OPEN', price: '60000', amount: '0.2', created_at: '2024-01-02T00:00:00Z' }])
        }
      >
        Set Orders
      </button>
    </div>
  )
}

describe('OrderProvider', () => {
  it('initializes with null values', () => {
    render(
      <OrderProvider>
        <TestComponent />
      </OrderProvider>
    )

    expect(screen.getByTestId('order').textContent).toBe('no-order')
    expect(screen.getByTestId('orders-count').textContent).toBe('0')
  })

  it('updates order and orders state correctly', async () => {
    render(
      <OrderProvider>
        <TestComponent />
      </OrderProvider>
    )

    // Set single order
    act(() => {
      screen.getByText('Set Order').click()
    })
    expect(screen.getByTestId('order').textContent).toBe('1')

    // Set multiple orders
    act(() => {
      screen.getByText('Set Orders').click()
    })
    expect(screen.getByTestId('orders-count').textContent).toBe('1')
  })

  it('throws error when used outside of OrderProvider', () => {
    expect(() => renderHook(() => useOrder())).toThrowError(
        'useOrder must be used inside OrderProvider'
    )
})
})
