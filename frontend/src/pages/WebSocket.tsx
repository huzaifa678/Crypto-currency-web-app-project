import React, { useEffect, useState, useRef } from "react";

export interface Market {
  market_id: string;
  base_currency: string;
  quote_currency: string;
  current_price: number;
}

interface MarketRowProps {
  market: Market;
  lastPrice?: number;
}

const MarketRow: React.FC<MarketRowProps> = ({ market, lastPrice }) => {
  const isUp = lastPrice !== undefined && market.current_price > lastPrice;
  const isDown = lastPrice !== undefined && market.current_price < lastPrice;

  return (
    <tr>
      <td className="px-4 py-2 align-middle">{market.base_currency}/{market.quote_currency}</td>
      <td
        className="px-4 py-2 align-middle font-bold transition-colors"
        style={{
          color: isUp ? "green" : isDown ? "red" : "black",
          fontWeight: "bold",
          transition: "color 0.3s ease",
        }}
      >
        {market.current_price.toFixed(2)}
      </td>
    </tr>
  );
};


const MarketsTable: React.FC = () => {
  const [markets, setMarkets] = useState<Market[]>([]);
  const lastPricesRef = useRef<Record<string, number>>({});

  useEffect(() => {
    const controller = new AbortController();

    async function streamTrades() {
      const symbols = ["BTCUSDT", "ETHUSDT"];
      const queryParams = symbols.map(s => `symbols=${encodeURIComponent(s)}`).join("&");
      const url = `http://localhost:8081/v1/stream_trades?${queryParams}`;

      try {
        const res = await fetch(url, {
          method: "GET",
          signal: controller.signal,
        });

        if (!res.body) {
          console.error("ReadableStream not supported by server or browser");
          return;
        }

        const reader = res.body.getReader();
        const decoder = new TextDecoder("utf-8");
        let buffer = "";

        while (true) {
          const { value, done } = await reader.read();
          if (done) break;

          buffer += decoder.decode(value, { stream: true });

          let boundary = buffer.indexOf("\n");

          while (boundary !== -1) {
            const jsonStr = buffer.slice(0, boundary).trim();
            buffer = buffer.slice(boundary + 1);

            if (jsonStr.length > 0) {
              try {
                const wrapped = JSON.parse(jsonStr);
                const trade = wrapped.result;

                const symbol = trade.symbol;
                const price = parseFloat(trade.price);

                if (!symbol || typeof symbol !== "string") {
                  console.warn("Missing or invalid symbol:", symbol);
                  continue;
                }

                if (isNaN(price)) {
                  console.warn("Invalid price:", trade.price);
                  continue;
                }

                setMarkets((prev) => {
                  const updated = [...prev];
                  const idx = updated.findIndex((m) => m.market_id === symbol);
                  if (idx >= 0) {
                    updated[idx] = { ...updated[idx], current_price: price };
                  } else {
                    if (symbol.length >= 6) {
                      updated.push({
                        market_id: symbol,
                        base_currency: symbol.slice(0, 3),
                        quote_currency: symbol.slice(3),
                        current_price: price,
                      });
                    } else {
                      console.warn("Symbol too short:", symbol);
                    }
                  }
                  return updated;
                });
              } catch (e) {
                console.error("JSON parse error:", e);
              }
            }

            boundary = buffer.indexOf("\n");
          }
        }
      } catch (err: any) {
        if (err.name !== "AbortError") {
          console.error("Stream error:", err);
        }
      }
    }

    streamTrades();

    return () => {
      controller.abort();
    };
  }, []);

  return (
  <div className="fixed top-20 left-1/2 -translate-x-1/4 max-w-4xl px-6 py-5">
    <div className="mb-4 px-10">
      <h1 className="text-xl font-bold">Live Trading Updates</h1>
      <p className="text-gray-500">Streaming prices directly from server...</p>
    </div>
      <table className="min-w-full">
        <thead className="bg-gray-100 sticky top-[3.5rem] z-10">
          <tr>
            <th className="px-4 py-2 border-b text-left">Pair</th>
            <th className="px-4 py-2 border-b text-left">Current Price</th>
          </tr>
        </thead>
      <tbody>
        {markets.map((market) => {
          const lastPrice = lastPricesRef.current[market.market_id];
          lastPricesRef.current[market.market_id] = market.current_price;

          return (
            <MarketRow
              key={market.market_id}
              market={market}
              lastPrice={lastPrice}
            />
          );
        })}
      </tbody>
    </table>
  </div>
);
};

export default MarketsTable;