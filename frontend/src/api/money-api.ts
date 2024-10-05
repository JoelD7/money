import {AccessToken, User} from "../types";
import axios, {AxiosError, AxiosResponse} from "axios";
import {keys} from "../utils";
import createAuthRefreshInterceptor from "axios-auth-refresh";
import {redirectToLogin, retryableRequest} from "./utils.ts";
import {logout} from "./auth.ts";

export const API_BASE_URL = process.env.REACT_APP_API_BASE_URL;

type ErrorHandler = () => Promise<void>;
const refreshTokenErrorHandlers = new Map<number, ErrorHandler>([
  [400, logoutWoApiRequest],
  [401, logoutWithApiRequest],
]);

export const axiosClient = axios.create({
  baseURL: API_BASE_URL,
});

createAuthRefreshInterceptor(axiosClient, refreshAuthInterceptor, {
  statusCodes: [401],
  pauseInstanceWhileRefreshing: true,
});

export async function getUser(): Promise<User> {
  const res: AxiosResponse = await axiosClient.get<User>(API_BASE_URL + "/users", {
    withCredentials: true,
    headers: {
      Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
    },
  })

  return res.data;
}

async function refreshToken() {
  await retryableRequest(async () => {
    const response = await axios.post(API_BASE_URL + "/auth/token", null, {
      withCredentials: true,
    });

    localStorage.setItem(keys.ACCESS_TOKEN, response.data.accessToken);
  });
}

async function refreshAuthInterceptor(
  failedRequest: AxiosError,
): Promise<string> {
  try {
    await refreshToken();
    return updateBearerToken(failedRequest);
  } catch (err) {
    console.error("Error refreshing the token: ", err);
    const axiosError = err as AxiosError;

    await handleRefreshTokenError(axiosError);

    return Promise.reject(err);
  }
}

function updateBearerToken(failedRequest: AxiosError): Promise<string> {
  const newToken = localStorage.getItem(keys.ACCESS_TOKEN);

  if (newToken && failedRequest.response) {
    failedRequest.response.config.headers.Auth = "Bearer " + newToken;
    return Promise.resolve(newToken);
  }

  console.error("Couldn't update Bearer token. Logging out.");
  redirectToLogin();
  return Promise.reject();
}

async function handleRefreshTokenError(axiosError: AxiosError) {
  if (!axiosError.response) {
    console.error("Couldn't read error response. Logging out.");
    redirectToLogin();
    return;
  }

  const logout = refreshTokenErrorHandlers.get(axiosError.response.status);
  if (logout) {
    await logout();
    return;
  }

  console.error(
    `Couldn't get refresh token error handler for the status code '${axiosError.response.status}'. Logging out.`,
  );
  redirectToLogin();
}

async function logoutWoApiRequest() {
  localStorage.removeItem(keys.ACCESS_TOKEN);
  redirectToLogin();
}

async function logoutWithApiRequest() {
  const accessToken = parseJwt(localStorage.getItem(keys.ACCESS_TOKEN));

  try {
    await logout({ username: accessToken.sub, password: "" });
  } catch (error) {
    console.error("Error logging out", error);
    redirectToLogin();
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
