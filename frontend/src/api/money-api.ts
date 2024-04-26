import { LoginCredentials, SignUpUser, User } from "../types";
import axios from "axios";
import { keys } from "../utils";

export const BASE_URL =
  "https://38qslpe8d9.execute-api.us-east-1.amazonaws.com/staging";

export function signUp(newUser: SignUpUser) {
  return axios.post(BASE_URL + "/auth/signup", newUser);
}

export function login(credentials: LoginCredentials) {
  return axios.post(BASE_URL + "/auth/login", credentials, {
    withCredentials: true, //required so that the browser will store the cookie. See more: https://developer.mozilla.org/en-US/docs/Web/API/Request/credentials
  });
}

export function getUser() {
  return axios
    .get<User>(BASE_URL + "/users/", {
      withCredentials: true,
      headers: {
        Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
      },
    })
}

export async function refreshToken() {
  axios
    .post(BASE_URL + "/auth/token", {}, {
      withCredentials: true,
    })
    .then((res) => {
      localStorage.setItem(keys.ACCESS_TOKEN, res.data.accessToken);
    })
    .catch((err) => {
      console.error("Error refreshing token: ", err);
    });
}
