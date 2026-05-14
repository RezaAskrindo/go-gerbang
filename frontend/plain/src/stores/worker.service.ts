import { useWebWorkerFn } from '@vueuse/core';

// export const baseHost = import.meta.env.MODE === "development" ? "http://localhost:9000" : `${location.origin}/backend`;
export const baseHost = `${location.origin}/backend`;

export const { workerFn, workerStatus, workerTerminate } = useWebWorkerFn(async (url: string) => {
  try {
    const response = await fetch(`${url}`, {
      credentials: 'include',
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      // NOTE: WORKAROUND FOR SERVICE UNAVAILABLE
      // If the service is unavailable, we will try to start it using a POST request.
      if (response.status === 503) {
        const credentials = btoa(`admin:@dmin9192`);
        const startService = await fetch(`https://siasuransi.com/go-services/`, {
          method: 'POST',
          mode: 'no-cors',
          headers: {
            'Authorization': `Basic ${credentials}`,
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ "app-name": "apigateway-9000" }),
        });

        await startService.json();
      };

      if (response.status === 401) return { status: false };
      
      console.error(`HTTP error! Status: ${response.status}`);
      return response.json();
    }
    
    return await response.json();
  } catch (error: unknown) {
    if (error instanceof Error) {
      console.error("Fetch error:", error.message);
    } else {
      console.error("Fetch error:", error);
    }
  }
});

export const getCSRFToken = async() => {
  if (workerStatus) workerTerminate();

  const csrf = await workerFn(`${baseHost}/secure-gateway-c`);
  workerTerminate();
  return csrf?.data;
}

export const getGoogleClientId = async() => {
  if (workerStatus) workerTerminate();

  const client_id = await workerFn(`${baseHost}/api/v1/auth/get-google-client-id`);
  workerTerminate();
  return client_id;
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
