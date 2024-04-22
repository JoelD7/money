import { LoginCredentials, SignUpUser, User } from "../types";
import axios, { AxiosError, AxiosResponse } from "axios";
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
  // const req = authRequest(
  //   axios.get<User>(BASE_URL + "/users/", {
  //     headers: {
  //       Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
  //     },
  //   }), 1
  // );
  //
  // req.catch((err) => {
  //   console.error("Error getting user: ", err);
  // });
  //
  // return req;

  return axios
    .get<User>(BASE_URL + "/users/", {
      withCredentials: true,
      headers: {
        Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
      },
    })
}

async function authRequest(request: Promise<AxiosResponse>, retryCount: number) {
  const req = request.then((res) => res).catch(async (err: AxiosError) => {
    if (isUnauthorized(err) && retryCount < 2) {
      await refreshToken();
      await authRequest(request, retryCount++)
      return
    } else {
      console.error("Error in auth request: ", err);
      throw err;
    }
  });

  return req;
}

export async function refreshToken() {
  axios
    .post(BASE_URL + "/auth/token")
    .then((res) => {
      localStorage.setItem(keys.ACCESS_TOKEN, res.data.accessToken);
    })
    .catch((err) => {
      console.error("Error refreshing token: ", err);
    });
}

function isUnauthorized(err: AxiosError): boolean {
  return err.response?.status === 401;
}
