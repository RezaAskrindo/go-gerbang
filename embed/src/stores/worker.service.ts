import { useWebWorkerFn } from '@vueuse/core';

export const baseHost = import.meta.env.MODE === "development" ? "http://localhost:9000" : "https://apiv1.siskor.web.id";
// export const baseHost = "https://apiv1.siskor.web.id";

export const { workerFn, workerStatus, workerTerminate } = useWebWorkerFn(async (url: string) => {
  try {
    const response = await fetch(url, {
      credentials: 'include',
      headers: {
        "Content-Type": "application/json",
      },
    });
    if (!response.ok) {
      if (response.status === 401) return { status: false }; 
      
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error("Fetch error:", error);
    throw error;
  }
});

export const getCSRFToken = async() => {
  if (workerStatus) workerTerminate();

  const csrf = await workerFn(`${baseHost}/secure-gateway-c`);
  workerTerminate();
  return csrf;
}

const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

export const getUserSession = async() => {
  // const session = await workerFn(`${baseHost}/api/v1/auth/get-session`);
  // workerTerminate();

  // MOCKING SESSION
  let session = { status: false , token: ''};

  await delay(200); 

  // setTimeout(() => {
    session = { status: true, token: 'Testing' };
  // }, 3000);
  
  return session;
}
