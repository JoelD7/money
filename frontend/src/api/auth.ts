import { Credentials, IdempotencyKVP, SignUpUser } from "../types";
import { API_BASE_URL } from "./money-api.ts";
import axios, { AxiosError } from "axios";
import { getIdempotencyKey, redirectToLogin, retryableRequest } from "./utils.ts";
import { keys } from "../utils"; // Promise<AxiosResponse<any, any>>

export function signUp(newUser: SignUpUser) {
  const idempotenceKVP: IdempotencyKVP = getIdempotencyKey(newUser, "", newUser.username);

  return axios
    .post(API_BASE_URL + "/auth/signup", newUser, {
      headers: {
        "Idempotency-Key": idempotenceKVP.idempotencyKey,
      },
    })
    .then((res) => {
      console.log("removing idempotency key item request success");
      localStorage.removeItem(idempotenceKVP.encodedRequestBody);
      return res;
    })
    .catch((err) => {
      console.log("removing idempotency key item request error", err);
      const axiosError = err as AxiosError;
      if (
        axiosError.response &&
        axiosError.response.status >= 400 &&
        axiosError.response.status < 500
      ) {
        console.log("removing idempotency key item non-500 error");
        localStorage.removeItem(idempotenceKVP.encodedRequestBody);
      }

      return Promise.reject(err);
    });
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
