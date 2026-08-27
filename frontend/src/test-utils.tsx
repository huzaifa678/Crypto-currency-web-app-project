import { ReactNode } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

export const createTestQueryClient = () =>
  new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
        refetchOnWindowFocus: false,
        gcTime: 0,
      },
    },
  });

export const withQueryClient = (ui: ReactNode) => (
  <QueryClientProvider client={createTestQueryClient()}>{ui}</QueryClientProvider>
);
