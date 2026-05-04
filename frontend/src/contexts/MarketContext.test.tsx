/// <reference types="vitest" />
import { beforeEach, describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import React from 'react'
import { MarketsProvider, useMarkets } from './MarketContext'

type Market = {
    market_id: string;
    name: string;
    base_currency: string;
    quote_currency: string;
    min_order_amount?: number;
    price_precision?: number;
    created_at?: string;
}

const TestComponent = () => {
const { market, setMarket } = useMarkets()

return ( <div> <div data-testid="market-count">{market.length}
    </div>
        <button
            onClick={() =>
            setMarket([{ market_id: '1', name: 'BTC', base_currency: 'BTC', quote_currency: 'USD', min_order_amount: 0.0001, price_precision: 2 }])}
        >
            Add Market 
        </button> 
    </div>
)
}

beforeEach(() => {
})

describe('MarketsProvider', () => {
it('initializes with an empty market list', () => {
render( <MarketsProvider> <TestComponent /> </MarketsProvider>
)


expect(screen.getByTestId('market-count').textContent).toBe('0')


})

it('updates market state when setMarket is called', () => {
render( <MarketsProvider> <TestComponent /> </MarketsProvider>)


fireEvent.click(screen.getByText('Add Market'))

expect(screen.getByTestId('market-count').textContent).toBe('1')


})

it('throws error when useMarkets is used outside MarketsProvider', () => {
// Suppress expected error logs from React
const spy = vi.spyOn(console, 'error').mockImplementation(() => {})


const OutsideComponent = () => {
  useMarkets()
  return <div>Outside</div>
}

expect(() => render(<OutsideComponent />)).toThrowError(
  'useMarkets must be used inside MarketsProvider'
)

spy.mockRestore()


})
})
