import { createContext, useContext, useState } from "react";
import { Market } from "../pages/Markets";

const MarketsContext = createContext<{market: Market[], setMarket: React.Dispatch<React.SetStateAction<Market[]>>} | null>(null);

export const MarketsProvider: React.FC<{children: React.ReactNode}> = ({ children }) => {
  const [market, setMarket] = useState<Market[]>([]);
  return (
    <MarketsContext.Provider value={{ market, setMarket }}>
      {children}
    </MarketsContext.Provider>
  );
};

export const useMarkets = () => {
  const ctx = useContext(MarketsContext);
  if (!ctx) throw new Error("useMarkets must be used inside MarketsProvider");
  return ctx;
};