export const endpoints = {
  // Auth - Customer
  auth: {
    customer: {
      register: "/api/auth/customer/register",
      login: "/api/auth/customer/login",
    },
    merchant: {
      register: "/api/auth/merchant/register",
      login: "/api/auth/merchant/login",
    },
  },

  // Products
  products: {
    list: "/api/products",
    listByMerchant: (merchantId: string) => `/api/products?merchant_id=${merchantId}`,
    byId: (id: string) => `/api/products/${id}`,
    create: "/api/products",
    update: (id: string) => `/api/products/${id}`,
    delete: (id: string) => `/api/products/${id}`,
    ask: "/api/products/ask",
  },

  // Merchants
  merchants: {
    list: "/api/merchants",
    byId: (id: string) => `/api/merchants/${id}`,
  },

  // Customers
  customers: {
    list: "/api/customers",
    byId: (id: string) => `/api/customers/${id}`,
  },

  // Orders
  orders: {
    list: "/api/orders",
    customer: "/api/orders/customer",
    merchant: "/api/orders/merchant",
    byId: (id: string) => `/api/orders/${id}`,
    create: "/api/orders",
    updateStatus: (id: string) => `/api/orders/${id}/status`,
  },

  // Shipping
  shipping: {
    provinces: "/api/shipping/provinces",
    cities: (provinceId?: string) => provinceId ? `/api/shipping/cities?province_id=${provinceId}` : "/api/shipping/cities",
    cost: (origin: string, destination: string, weight: number, courier?: string) => {
      let url = `/api/shipping/cost?origin=${origin}&destination=${destination}&weight=${weight}`;
      if (courier) url += `&courier=${courier}`;
      return url;
    },
    track: (waybill: string, courier: string) => `/api/shipping/track?waybill=${waybill}&courier=${courier}`,
  },

  // Chat
  chat: {
    history: "/api/chat/history",
    saveMessage: "/api/chat/messages",
    clearHistory: "/api/chat/history",
  },
} as const;
