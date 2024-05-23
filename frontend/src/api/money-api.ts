import { Credentials, SignUpUser, User } from "../types";
import axios, { AxiosError } from "axios";
import { keys } from "../utils";
import createAuthRefreshInterceptor from "axios-auth-refresh";
import { retryableRequest } from "./utils.ts";

export const API_BASE_URL = process.env.REACT_APP_API_BASE_URL;

const axiosClient = axios.create({
  baseURL: API_BASE_URL,
});

createAuthRefreshInterceptor(axiosClient, refreshAuthInterceptor, {
  statusCodes: [401],
  pauseInstanceWhileRefreshing: true,
});

export function signUp(newUser: SignUpUser) {
  return axiosClient.post(API_BASE_URL + "/auth/signup", newUser);
}

export function login(credentials: Credentials) {
  return axiosClient.post(API_BASE_URL + "/auth/login", credentials, {
    withCredentials: true, //required so that the browser will store the cookie. See more: https://developer.mozilla.org/en-US/docs/Web/API/Request/credentials
  });
}

export function getUser() {
  return axiosClient.get<User>(API_BASE_URL + "/users", {
    withCredentials: true,
    headers: {
      Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
    },
  });
}

async function refreshToken() {
  await retryableRequest(async () => {
    const response = await axios.post(API_BASE_URL + "/auth/token", null, {
      withCredentials: true,
    });

    localStorage.setItem(keys.ACCESS_TOKEN, response.data.accessToken);
  });
}

export async function logout(credentials: Credentials) {
  await retryableRequest(async () => {
    await axiosClient.post(API_BASE_URL + "/auth/logout", credentials, {
      withCredentials: true,
    });

    localStorage.removeItem(keys.ACCESS_TOKEN);
  });
}

async function refreshAuthInterceptor(failedRequest: AxiosError): Promise<string> {
  await refreshToken();

  const newToken = localStorage.getItem(keys.ACCESS_TOKEN);

  if (newToken && failedRequest.response) {
    failedRequest.response.config.headers.Auth = "Bearer " + newToken;
    return Promise.resolve(newToken);
  } else {
    return Promise.reject();
  }
}
