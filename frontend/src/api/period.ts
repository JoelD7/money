import {Period, PeriodSchema, PeriodsSchema} from "../types";
import {keys} from "../utils/index.ts";
import {API_BASE_URL, axiosClient} from "./money-api.ts";
import {AxiosResponse} from "axios";

export async function getPeriod(period: string): Promise<Period> {
    const res: AxiosResponse = await axiosClient.get<Period>(
        API_BASE_URL + `/periods/${period}`,
        {
            withCredentials: true,
            headers: {
                Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
            },
        },
    );

    return PeriodSchema.parse(res.data);
}

export async function getPeriods(): Promise<Period[]> {
    const res: AxiosResponse = await axiosClient.get<Period[]>(
        API_BASE_URL + "/periods",
        {
            withCredentials: true,
            headers: {
                Auth: `Bearer ${localStorage.getItem(keys.ACCESS_TOKEN)}`,
            },
        },
    );

    return PeriodsSchema.parse(res.data);
}

// TODO: Update the current period in localstorage after creating a new one
export function createPeriod() {
}
