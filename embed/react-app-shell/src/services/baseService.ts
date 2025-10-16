import useSWR, { type SWRConfiguration } from "swr";

const env = import.meta.env;
export const BackendUrlBase = env.DEV ? "http://localhost:9000" : "/backend";
export const FrontendUrl = env.DEV ? "http://localhost:5173" : window.location.origin;
// export const BackendUrlBase = "http://localhost:9000";
// export const FrontendUrl = "http://localhost:5173";

export const fetchSWR = (url: string) => fetch(url, {credentials: "include"}).then(r => r.json())

export const SWRBasedConfig: SWRConfiguration = {
  revalidateOnFocus: false
}

export const SWRDashboardConfig: SWRConfiguration = {
  refreshInterval: 5000
}

export async function FetchCsrfToken(): Promise<string> {
  const urlCsrf = `${BackendUrlBase}/secure-gateway-c`;
  const res = await fetch(urlCsrf, { credentials: 'include' });
  if (!res.ok) return "";
  const data = await res.json();
  return data?.data;
}

export function GetAuthSession() {
  const { data, isLoading } = useSWR(`${BackendUrlBase}/api/v1/auth/get-session`, (url: string) => fetch(url, {credentials: 'include'}).then(r => r.json()));
  return { data, isLoading };
}