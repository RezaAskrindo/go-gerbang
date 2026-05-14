import useSWR, { type SWRConfiguration } from "swr"
import { BackendUrlBase, FetchCsrfToken } from "./baseService"

// export const fetchSWR = (url: string) => fetch(url, {credentials: "include"}).then(r => r.json())
export const fetchSWR = async (url: string, options = {}, retry = true): Promise<any> => {
  const response = await fetch(url, {
    ...options,
    credentials: "include",
  });

  if (response.status !== 401) {
    if (!response.ok) {
      throw new Error(`Error: ${response.status}`);
    }
    return response.json();
  }

  if (retry) {
    const refreshed = await getRefreshToken();

    if (refreshed) {
      return fetchSWR(url, options, false);
    }
  }

  throw new Error("Unauthorized");
}

async function getRefreshToken() {
  try {
    const getCsrf = await FetchCsrfToken();
    const domain = localStorage.getItem("domain");
    const refreshTokenUrl = `${BackendUrlBase}/api/v1/auth/less-secure/refresh-token?httponly=true&domain=${domain !== null ? domain : 'localhost'}`;
    return await (await fetch(refreshTokenUrl, {
      method: "POST",
      credentials: "include",
      headers: {
        'Content-Type': 'application/json',
        'X-SGCsrf-Token': getCsrf
      },
    })).json();
  } catch (error) {
    console.log(error)
    return null;
  }
}


export const SWRBasedConfig: SWRConfiguration = {
  revalidateOnFocus: false
}

export const SWRDashboardConfig: SWRConfiguration = {
  refreshInterval: 5000
}
export function GetAuthSession() {
  // const { data, isLoading } = useSWR(`${BackendUrlBase}/api/v1/auth/get-session`, fetchSWR); // using Session
  const { data, isLoading } = useSWR(`${BackendUrlBase}/api/v1/auth/get-jwt-info`, fetchSWR, {
    revalidateOnFocus: false
  });// using JWT

  return { data, isLoading };
}

export function CheckMigration() {
  const { data, isLoading } = useSWR(`${BackendUrlBase}/check-migration`, fetchSWR);
  return { data, isLoading };
}

export function useConfiguration(group: string, config_name?: string) {
  let url = `${BackendUrlBase}/Configuration/${group}`
  if (config_name) {
    url = `${BackendUrlBase}/Configuration/${group}?config_name=${config_name}`
  }
  
  const { data, error, isLoading, mutate } = useSWR(url, fetchSWR)
  return {
    data,
    isLoading,
    error,
    mutate
  }
}

export const useDeleteConfiguration = async (
  group: string,
  name: string
) => {
  const getCsrf = await FetchCsrfToken();
  
  const res = await fetch(`${BackendUrlBase}/Configuration/${group}?config_name=${name}`, {
    method: "DELETE",
    credentials: "include",
    headers: { "Content-Type": "application/json", "X-SGCsrf-Token": getCsrf },
  })

  if (!res.ok) throw new Error("Request failed");
  
  const response = await res.json();
  
  if (!response.status) throw new Error(response.message || "Failed to login");
  
  return response;
}