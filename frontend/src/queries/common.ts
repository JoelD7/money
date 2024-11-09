import { AxiosError } from "axios";
import { QUERY_RETRIES } from "./queries.ts";

export function queryRetryFn(failureCount: number, e: AxiosError) {
  if (failureCount > QUERY_RETRIES) {
    return false;
  }

  return e.response ? e.response.status !== 404 : true;
}
