import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  reactCompiler: true,
  async rewrites() {
    return [
      {
        source: "/api/backend/:path*",
        destination: `${process.env.BACKEND_URL || "http://127.0.0.1:8083"}/api/:path*`,
      },
    ];
  },
};

export default nextConfig;
