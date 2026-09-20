function resolveWsBaseUrl(): string {
  const explicit = import.meta.env.VITE_WS_BASE_URL as string | undefined;

  if (explicit) {
    return explicit.replace(/\/$/, "");
  }

  return "ws://localhost:8080";
}

export const WS_BASE_URL = resolveWsBaseUrl();

export function wsUrl(path: string): string {
  return `${WS_BASE_URL}${path.startsWith("/") ? path : `/${path}`}`;
}