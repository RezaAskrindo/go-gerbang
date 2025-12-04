// import { useSWRLite } from "../composables/useSWRLite";
import useSWRV, { mutate } from 'swrv'

const env = import.meta.env;
// export const BackendUrlBase = env.DEV ? "http://localhost:9000" : "/backend";
export const BackendUrlBase = "http://localhost:9000";

const FrontendUrl = env.DEV ? "http://localhost:3000" : window.location.origin;
// export const BackendUrlBase = "http://localhost:9000";
// export const FrontendUrl = "http://localhost:5173";

export async function FetchCsrfToken(): Promise<string> {
  const urlCsrf = `${BackendUrlBase}/secure-gateway-c`;
  const res = await fetch(urlCsrf, { credentials: 'include' });
  if (!res.ok) return "";
  const data = await res.json();
  return data?.data;
}

export function GetAuthSession() {
  // const { data, isLoading, mutate } = useSWRLite(
  const { data, isLoading } = useSWRV(
    `${BackendUrlBase}/api/v1/auth/get-session`,
    (url: string) => fetch(url, { credentials: 'include' }).then(r => r.json())
  )

  function prefetch() {
    mutate(
      `${BackendUrlBase}/api/v1/auth/get-session`,
      (url: string) => fetch(url, { credentials: 'include' }).then(r => r.json())
    )
  }

  return { data, isLoading, prefetch };
}

export async function LoginUser(form:{ identity: string; password: string }, domain?: string) {
  const getCsrf = await FetchCsrfToken();
  const url = `${BackendUrlBase}/api/v1/auth/login?httponly=true&session=true&domain=${domain ?? 'localhost'}&url=${FrontendUrl}`;

  const res = await fetch(url, {
    method: "POST",
    credentials: "include",
    headers: { 
      "Content-Type": "application/json", 
      "X-SGCsrf-Token": getCsrf 
    },
    body: JSON.stringify(form),
  });

  const data = await res.json();

  if (!res.ok) {
    throw new Error(data.message || "Request failed");
  }
  if (!data.status) {
    throw new Error(data.message || "Failed to login");
  }

  console.log("here")

  // const { prefetch } = GetAuthSession();

  // prefetch();

  return data;
}

export function LogoutUser() {
  window.location.href = `${BackendUrlBase}/api/v1/auth/logout?redirectUrl=${FrontendUrl}`
}