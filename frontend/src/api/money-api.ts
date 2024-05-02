import { LoginCredentials, SignUpUser, User } from "../types";
import axios, { AxiosError } from "axios";
import { keys } from "../utils";

export const MAX_RETRIES: number = 3;
const BACKOFF_TIME_MS: number = 1000;

export const API_BASE_URL = process.env.REACT_APP_API_BASE_URL;

export function signUp(newUser: SignUpUser) {
  return axios.post(API_BASE_URL + "/auth/signup", newUser);
}

export function login(credentials: LoginCredentials) {
  return axios.post(API_BASE_URL + "/auth/login", credentials, {
    withCredentials: true, //required so that the browser will store the cookie. See more: https://developer.mozilla.org/en-US/docs/Web/API/Request/credentials
  });
}

export function getUser() {
  return axios.get<User>(API_BASE_URL + "/users", {
    withCredentials: true,
    headers: {
      Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
    },
  });
}

export async function refreshToken() {
  await retryableRequest(async () => {
    const response = await axios.post(API_BASE_URL + "/auth/token", null, {
      withCredentials: true,
    });

    localStorage.setItem(keys.ACCESS_TOKEN, response.data.accessToken);
  });
}

function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

async function retryableRequest(request: () => Promise<void>) {
  let myErr: AxiosError | undefined;

  for (let i = 0; i < MAX_RETRIES; i++) {
    try {
      await request();
      return;
    } catch (err) {
      myErr = err as AxiosError;
      await sleep(BACKOFF_TIME_MS);
    }
  }

  if (myErr) {
    throw myErr;
  }
}

export async function logout() {
  await retryableRequest(async () => {
    await axios.post(API_BASE_URL + "/auth/logout", null, {
      withCredentials: true,
    });

    localStorage.removeItem(keys.ACCESS_TOKEN);
  });
}
