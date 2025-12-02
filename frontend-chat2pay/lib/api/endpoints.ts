export const endpoints = {
  health: () => "/health",
  payments: {
    root: () => "/payments",
    byId: (paymentId: string) => `/payments/${paymentId}`,
    capture: (paymentId: string) => `/payments/${paymentId}/capture`,
  },
  customers: {
    root: () => "/customers",
    byId: (customerId: string) => `/customers/${customerId}`,
  },
};
