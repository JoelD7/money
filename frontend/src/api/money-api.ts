import {SignUpUser} from "../types";
import axios from "axios"
export const BASE_URL = "https://38qslpe8d9.execute-api.us-east-1.amazonaws.com/staging"

export function signUp(newUser: SignUpUser){
    return axios.post(BASE_URL, newUser)
}