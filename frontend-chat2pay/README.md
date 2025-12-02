# chat2pay frontend

Next.js app dengan infrastruktur integrasi API bawaan.

## Environment

Salin `.env.example` lalu isi nilai yang sesuai.

```bash
cp .env.example .env.local
```

| Key | Deskripsi |
| --- | --- |
| `NEXT_PUBLIC_API_BASE_URL` | Base URL API (contoh `https://api.example.com`). Dipakai di server dan client. |

## Menjalankan

```bash
npm install
npm run dev
```

Buka http://localhost:3000

## Cara pakai API client

API helper tersedia di `lib/api`:
- `apiClient` (default): memakai `NEXT_PUBLIC_API_BASE_URL`.
- `endpoints`: registry path API.
- `ApiError`: class error dengan status/payload.

### GET example

```ts
import { apiClient, endpoints } from "@/lib/api";

type Payment = { id: string; amount: number; status: string };

export async function listPayments() {
  const { data } = await apiClient.get<Payment[]>(
    endpoints.payments.root(),
    { query: { page: 1, per_page: 10 } }
  );

  return data;
}
```

### POST example

```ts
import { apiClient, endpoints } from "@/lib/api";

type Payment = { id: string; amount: number; status: string };

export async function createPayment(payload: {
  amount: number;
  currency: string;
}) {
  const { data } = await apiClient.post<Payment>(
    endpoints.payments.root(),
    { json: payload }
  );

  return data;
}
```

### Error handling singkat

```ts
import { ApiError } from "@/lib/api";

try {
  await listPayments();
} catch (error) {
  if (error instanceof ApiError) {
    console.error("API failed", error.status, error.payload);
  }
}
```
