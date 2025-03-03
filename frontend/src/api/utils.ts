import { AxiosError } from "axios";

const BACKOFF_TIME_MS: number = 1000;
export const MAX_RETRIES: number = 3;

export function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export async function retryableRequest(request: () => Promise<void>) {
  let myErr;

  for (let i = 0; i < MAX_RETRIES; i++) {
    try {
      await request();
      myErr = undefined;
      return;
    } catch (err) {
      myErr = err as AxiosError;
      if (myErr.response && myErr.response.status >= 400 && myErr.response.status < 500) {
        throw myErr;
      }
    }

    await sleep(BACKOFF_TIME_MS);
  }

  if (myErr) {
    throw myErr;
  }
}

export function redirectToLogin() {
  window.location.replace("/login");
}

export function buildQueryParams(
  startKey: string = "",
  pageSize: number = 10,
  sortOrder: string = "",
  sortBy: string = "",
): string[] {
  const params = [];

  if (startKey) {
    params.push(`start_key=${startKey}`);
  }

  if (pageSize) {
    params.push(`page_size=${pageSize}`);
  }

  if (sortOrder !== "") {
    params.push(`sort_order=${sortOrder}`);
  }

  if (sortBy !== "") {
    params.push(`sort_by=${sortBy}`);
  }

  return params;
}
