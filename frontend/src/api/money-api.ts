import { AccessToken, APIError, Credentials, SignUpUser, User } from "../types";
import axios, { AxiosError } from "axios";
import { keys } from "../utils";
import createAuthRefreshInterceptor from "axios-auth-refresh";
import { retryableRequest } from "./utils.ts";

export const API_BASE_URL = process.env.REACT_APP_API_BASE_URL;

type ErrorHandler = () => Promise<void>;
const refreshTokenErrorHandlers = new Map<number, ErrorHandler>([
  [400, logoutWoApiRequest],
  [401, logoutWithApiRequest],
]);

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
    window.location.pathname = "/login";
  });
}

async function refreshAuthInterceptor(
  failedRequest: AxiosError,
): Promise<string> {
  try {
    await refreshToken();

    const newToken = localStorage.getItem(keys.ACCESS_TOKEN);

    if (newToken && failedRequest.response) {
      failedRequest.response.config.headers.Auth = "Bearer " + newToken;
      return Promise.resolve(newToken);
    } else {
      console.error("Couldn't update Bearer token. Logging out.")
      window.location.pathname = "/login";
      return Promise.reject();
    }
  } catch (err) {
    console.error("Error refreshing the token: ", err);
    const axiosError = err as AxiosError;

    if (axiosError.response) {
      const apiErr: APIError = axiosError.response.data as APIError;
      const logout = refreshTokenErrorHandlers.get(apiErr.http_code);
      if (logout) {
        await logout();
      }
    }

    return Promise.reject(err);
  }
}

async function logoutWoApiRequest() {
  localStorage.removeItem(keys.ACCESS_TOKEN);
  window.location.pathname = "/login";
}

async function logoutWithApiRequest() {
  const accessToken = parseJwt(localStorage.getItem(keys.ACCESS_TOKEN));

  try {
    await logout({ username: accessToken.sub, password: "" });
  } catch (error) {
    console.error("Error logging out", error);
    window.location.pathname = "/login";
  }
}

function parseJwt(token: string | null): AccessToken {
  if (!token) {
    return { sub: "", exp: 0, iat: 0 };
  }

  const base64Url = token.split(".")[1];
  const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
  const jsonPayload = decodeURIComponent(
    window
      .atob(base64)
      .split("")
      .map(function (c) {
        return "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2);
      })
      .join(""),
  );

  return JSON.parse(jsonPayload);
}
