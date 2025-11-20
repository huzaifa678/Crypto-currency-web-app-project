import { createContext, useContext, useState } from "react";
import { Order } from "../pages/Orders";

type OrderContextType = {
  order: Order | null;
  setOrder: React.Dispatch<React.SetStateAction<Order | null>>;
  orders: Order[] | null;
  setOrders: React.Dispatch<React.SetStateAction<Order[] | null>>;
};

const OrderContext = createContext<OrderContextType | null>(null);

export const OrderProvider: React.FC<{children: React.ReactNode}> = ({ children }) => {
    const [order, setOrder] = useState<Order | null>(null);
    const [orders, setOrders] = useState<Order[] | null>(null);
    return (
        <OrderContext.Provider value={{ order, setOrder, orders, setOrders }}>
            {children}
        </OrderContext.Provider>
    );
};

export const useOrder = () => {
    const ctx = useContext(OrderContext);
    if (!ctx) throw new Error("useOrder must be used inside OrderProvider");
    return ctx;
}