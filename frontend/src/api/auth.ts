import { Credentials, SignUpUser } from "../types";
import { API_BASE_URL } from "./money-api.ts";
import axios from "axios";
import { redirectToLogin, retryableRequest } from "./utils.ts";
import { keys } from "../utils";

export function signUp(newUser: SignUpUser) {
  return axios.post(API_BASE_URL + "/auth/signup", newUser);
}

export function login(credentials: Credentials) {
  return axios.post(API_BASE_URL + "/auth/login", credentials, {
    withCredentials: true, //required so that the browser will store the cookie. See more: https://developer.mozilla.org/en-US/docs/Web/API/Request/credentials
  });
}

export async function logout(credentials: Credentials) {
  await retryableRequest(async () => {
    await axios.post(API_BASE_URL + "/auth/logout", credentials, {
      withCredentials: true,
    });

    localStorage.removeItem(keys.ACCESS_TOKEN);
    localStorage.removeItem(keys.CURRENT_PERIOD);
    redirectToLogin();
  });
}
