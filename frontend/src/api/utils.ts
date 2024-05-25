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
      myErr = err;
    }
    await sleep(BACKOFF_TIME_MS);
  }

  if (myErr) {
    throw myErr;
  }
}
