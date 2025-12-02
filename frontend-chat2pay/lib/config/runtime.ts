export type RuntimeConfig = {
 
  apiBaseUrl: string;
};

const removeTrailingSlash = (value?: string | null) =>
  value ? value.replace(/\/+$/, "") : "";

let cachedConfig: RuntimeConfig | null = null;

export function getRuntimeConfig(): RuntimeConfig {
  if (cachedConfig) return cachedConfig;

  const apiBaseUrl = removeTrailingSlash(process.env.NEXT_PUBLIC_API_BASE_URL);

  if (!apiBaseUrl && process.env.NODE_ENV !== "production") {
    console.warn(
      "NEXT_PUBLIC_API_BASE_URL belum diset."
    );
  }

  cachedConfig = { apiBaseUrl: apiBaseUrl ?? "" };

  return cachedConfig;
}
