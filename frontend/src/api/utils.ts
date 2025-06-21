import { AxiosError } from "axios";
import { v4 as uuidv4 } from "uuid";
import { IdempotencyKVP } from "../types";

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
  savingGoalID: string = "",
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

  if (savingGoalID !== "") {
    params.push(`saving_goal_id=${savingGoalID}`);
  }

  return params;
}

// getIdempotencyKey returns a previously saved idempotency key associated with the a request body, or generates a new one if it
// doesn't exist.
// Returns an IdempotencyKVP object that besides the idempotency key, also holds the associated encoded request body so
// that a caller can delete the local storage item on a succesful request.
export function getIdempotencyKey(data: unknown, accessToken: string, username: string): IdempotencyKVP {
  const encodedBody: string = window.btoa(encodeURIComponent(JSON.stringify(data)));
  let idempotencyKey: string | null = localStorage.getItem(encodedBody);

  if (username === ""){
    username = getUsernameFromAccessToken(accessToken);
  }

  if (!idempotencyKey) {
    idempotencyKey = `${uuidv4()}:${username}`;
    localStorage.setItem(encodedBody, idempotencyKey);
    return {
      encodedRequestBody: encodedBody,
      idempotencyKey,
    };
  }

  return {
    encodedRequestBody: encodedBody,
    idempotencyKey,
  };
}

function getUsernameFromAccessToken(accessToken: string): string {
  if (accessToken === "") {
    return "";
  }

  const base64Url = accessToken.split(".")[1];
  const base64 = base64Url.replace(/_/g, "/").replace(/-/g, "+");
  const jsonPayload = decodeURIComponent(atob(base64).split("=")[0]);

  return JSON.parse(jsonPayload).sub;
}