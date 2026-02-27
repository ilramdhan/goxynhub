"use client";

import { createContext, useContext, useState, useCallback, ReactNode } from "react";

type ToastType = "success" | "error" | "warning" | "info";

interface Toast {
  id: string;
  type: ToastType;
  title: string;
  message?: string;
  duration?: number;
}

interface ToastContextValue {
  toasts: Toast[];
  toast: (options: Omit<Toast, "id">) => void;
  success: (title: string, message?: string) => void;
  error: (title: string, message?: string) => void;
  warning: (title: string, message?: string) => void;
  info: (title: string, message?: string) => void;
  dismiss: (id: string) => void;
}

const ToastContext = createContext<ToastContextValue | null>(null);

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const dismiss = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const toast = useCallback(
    (options: Omit<Toast, "id">) => {
      const id = Math.random().toString(36).slice(2);
      const newToast: Toast = { id, duration: 4000, ...options };

      setToasts((prev) => [...prev, newToast]);

      // Auto-dismiss
      setTimeout(() => {
        dismiss(id);
      }, newToast.duration);
    },
    [dismiss]
  );

  const success = useCallback(
    (title: string, message?: string) => toast({ type: "success", title, message }),
    [toast]
  );
  const error = useCallback(
    (title: string, message?: string) => toast({ type: "error", title, message, duration: 6000 }),
    [toast]
  );
  const warning = useCallback(
    (title: string, message?: string) => toast({ type: "warning", title, message }),
    [toast]
  );
  const info = useCallback(
    (title: string, message?: string) => toast({ type: "info", title, message }),
    [toast]
  );

  return (
    <ToastContext.Provider value={{ toasts, toast, success, error, warning, info, dismiss }}>
      {children}
      <ToastContainer toasts={toasts} onDismiss={dismiss} />
    </ToastContext.Provider>
  );
}

export function useToast(): ToastContextValue {
  const ctx = useContext(ToastContext);
  if (!ctx) throw new Error("useToast must be used within ToastProvider");
  return ctx;
}

const TOAST_STYLES: Record<ToastType, { bg: string; icon: string; border: string }> = {
  success: { bg: "bg-green-50", icon: "✅", border: "border-green-200" },
  error: { bg: "bg-red-50", icon: "❌", border: "border-red-200" },
  warning: { bg: "bg-yellow-50", icon: "⚠️", border: "border-yellow-200" },
  info: { bg: "bg-blue-50", icon: "ℹ️", border: "border-blue-200" },
};

function ToastContainer({
  toasts,
  onDismiss,
}: {
  toasts: Toast[];
  onDismiss: (id: string) => void;
}) {
  if (toasts.length === 0) return null;

  return (
    <div className="fixed bottom-4 right-4 z-[100] flex flex-col gap-2 max-w-sm w-full">
      {toasts.map((toast) => {
        const styles = TOAST_STYLES[toast.type];
        return (
          <div
            key={toast.id}
            className={`${styles.bg} ${styles.border} border rounded-xl p-4 shadow-lg flex items-start gap-3 animate-fade-in`}
          >
            <span className="text-lg flex-shrink-0">{styles.icon}</span>
            <div className="flex-1 min-w-0">
              <p className="font-semibold text-gray-900 text-sm">{toast.title}</p>
              {toast.message && (
                <p className="text-gray-600 text-xs mt-0.5">{toast.message}</p>
              )}
            </div>
            <button
              onClick={() => onDismiss(toast.id)}
              className="text-gray-400 hover:text-gray-600 flex-shrink-0 text-lg leading-none"
            >
              ×
            </button>
          </div>
        );
      })}
    </div>
  );
}
