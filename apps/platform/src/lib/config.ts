declare global {
  interface Window {
    __ENV__?: Record<string, string>
  }
}

const env = (key: string): string | undefined =>
  window.__ENV__?.[key] || import.meta.env[key]

export const config = {
  env: import.meta.env.MODE || 'development',
  landingUrl: 'https://www.blocknext.ai',
  platformApiUrl:
    env('VITE_PLATFORM_API_URL') || 'VITE_PLATFORM_API_URL is not set',
  mcpApiUrl: env('VITE_MCP_API_URL') || 'VITE_MCP_API_URL is not set',
}
