import { User } from "../types";
import { keys } from "../utils";
import { axiosClient } from "./money-api";

export const API_BASE_URL = process.env.REACT_APP_API_BASE_URL;

export function patchUser(user: User) {
    return axiosClient.patch(API_BASE_URL + `/users/${user.username}`, user, {
        withCredentials: true,
        headers: {
            Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
        },
    });
}