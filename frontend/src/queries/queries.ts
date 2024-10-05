import {useQuery} from "@tanstack/react-query";
import api from "../api";
import {keys, utils} from "../utils";
import {AxiosError} from "axios";
import {User} from "../types";
import {INCOME, PERIOD, PERIOD_STATS, PERIODS, USER} from "./keys";

export const QUERY_RETRIES = 2;

export const incomeKeys = {
    all: [{scope: INCOME}] as const,
    list: (pageSize?: number, startKey?: string, period?: string) =>
        [
            {
                ...incomeKeys.all[0],
                pageSize,
                startKey,
                period,
            },
        ] as const,
};

export function useGetUser() {
    return useQuery({
        queryKey: [USER],
        queryFn: () => {
            const result: Promise<User> = api.getUser();
            result.then((res) => {
                localStorage.setItem(keys.CURRENT_PERIOD, res.current_period);
            });

            return result;
        },
    });
}

export function useGetPeriod(user?: User) {
    const periodID =
        user?.current_period || localStorage.getItem(keys.CURRENT_PERIOD) || "";

    return useQuery({
        queryKey: [PERIOD],
        queryFn: () => api.getPeriod(periodID),
        enabled: periodID !== "",
        retry: (failureCount: number, e: AxiosError) => {
            if (failureCount > QUERY_RETRIES) {
                return false;
            }

            return e.response ? e.response.status !== 404 : true;
        },
    });
}

export function useGetPeriods() {
    return useQuery({
        queryKey: [PERIODS],
        queryFn: () => api.getPeriods(),
        retry: (failureCount: number, e: AxiosError) => {
            if (failureCount > QUERY_RETRIES) {
                return false;
            }

            return e.response ? e.response.status !== 404 : true;
        },
    });
}

export function useGetPeriodStats(user?: User) {
    const periodID =
        user?.current_period || localStorage.getItem(keys.CURRENT_PERIOD) || "";

    return useQuery({
        queryKey: [PERIOD_STATS, periodID],
        queryFn: () => api.getPeriodStats(periodID),
        enabled: periodID !== "",
        retry: (failureCount: number, e: AxiosError) => {
            if (failureCount > QUERY_RETRIES) {
                return false;
            }

            return e.response ? e.response.status !== 404 : true;
        },
    });
}

export function useGetIncome() {
    // eslint-disable-next-line prefer-const
    let {pageSize, startKey, period} = utils.useTransactionsParams();

    if (!period) {
        period = localStorage.getItem(keys.CURRENT_PERIOD) || "";
    }

    return useQuery({
        queryKey: incomeKeys.list(pageSize, startKey, period),
        queryFn: api.getIncomeList,
        retry: (failureCount: number, e: AxiosError) => {
            if (failureCount > QUERY_RETRIES) {
                return false;
            }

            return e.response ? e.response.status !== 404 : true;
        },
    });
}
