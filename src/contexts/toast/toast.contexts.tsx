// src/components/ToastProvider.tsx

import React from "react";
import { Toaster } from "react-hot-toast";

interface ToastProviderProps {
    children: React.ReactNode;
}

const ToastProvider: React.FC<ToastProviderProps> = ({ children }) => {
    return (
        <>
            {children}
            {/* The Toaster component is where react-hot-toast renders all toasts. */}
            <Toaster
                position="top-right" // You can change this (e.g., 'bottom-center')
                reverseOrder={false} // New toasts appear at the bottom of the stack
                gutter={16} // Spacing between multiple toasts
                containerClassName="" // No default styling applied to the container itself
                containerStyle={{}} // No default inline styles
                toastOptions={{
                    duration: 3000, // Default duration for toasts if not overridden
                }}
            />
        </>
    );
};

export default ToastProvider;
