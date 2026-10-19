import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone", // self-contained server.js for the slim runtime image
};

export default nextConfig;
