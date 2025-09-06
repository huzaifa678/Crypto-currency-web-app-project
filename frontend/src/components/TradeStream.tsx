import { useEffect, useRef, useState } from "react";
import { Market } from "../pages/websocket";

interface Trade {
  symbol: string;
  price: number;
}

export function useTradeStream() {
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
  return markets;
}