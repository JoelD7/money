import {Period, PeriodSchema, PeriodsSchema} from "../types";
import {keys} from "../utils/index.ts";
import {API_BASE_URL, axiosClient} from "./money-api.ts";
import {AxiosResponse} from "axios";

export async function getPeriod(period: string) {
    const res: AxiosResponse = await axiosClient.get<Period>(
        API_BASE_URL + `/periods/${period}`,
        {
            withCredentials: true,
            headers: {
                Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
            },
        },
    );

    try {
        return PeriodSchema.parse(res.data);
    } catch (e) {
        console.error("[money] - Error parsing GET period response", e)
        return Promise.reject(e)
    }
}

export async function getPeriods() {
    const res: AxiosResponse = await axiosClient.get<Period[]>(
        API_BASE_URL + "/periods",
        {
            withCredentials: true,
            headers: {
                Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
            },
        },
    );

    try {
        return PeriodsSchema.parse(res.data);
    } catch (e) {
        console.error("[money] - Error parsing GET periods response", e)
        return Promise.reject(e)
    }
}

// TODO: Update the current period in localstorage after creating a new one
export function createPeriod() {
}
