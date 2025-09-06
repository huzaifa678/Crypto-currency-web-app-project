import { createContext, useContext, useState } from "react";
import { Order } from "../pages/Orders";

const OrderContext = createContext<{order: Order | null, setOrder: React.Dispatch<React.SetStateAction<Order | null>>} | null>(null);

export const OrderProvider: React.FC<{children: React.ReactNode}> = ({ children }) => {
    const [order, setOrder] = useState<Order | null>(null);
    return (
        <OrderContext.Provider value={{ order, setOrder }}>
            {children}
        </OrderContext.Provider>
    );
};

export const useOrder = () => {
    const ctx = useContext(OrderContext);
    if (!ctx) throw new Error("useOrder must be used inside OrderProvider");
    return ctx;
}