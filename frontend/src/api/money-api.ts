import { LoginCredentials, SignUpUser, User } from "../types";
import axios from "axios";
import { keys } from "../utils";

export const BASE_URL =
  "https://38qslpe8d9.execute-api.us-east-1.amazonaws.com/staging";

export function signUp(newUser: SignUpUser){
    return axios.post(BASE_URL+"/auth/signup", newUser)
}

export function login(credentials: LoginCredentials){
    return axios.post(BASE_URL+"/auth/login", credentials)
}

export function getUser(){
    return axios.get<User>(BASE_URL+"/users/", {
        headers: {
            "Auth": `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`
        }
    })
}