import { createFileRoute, redirect } from "@tanstack/react-router";
import { Home } from "../pages";
import { store } from "../store";
import { z } from "zod";
import { zodSearchValidator } from "@tanstack/router-zod-adapter";

function isAuth() {
  return store.getState().authReducer.isAuthenticated;
}

const expensesSearchSchema = z.object({
  categories: z.string().default(""),
  pageSize: z.number().default(10),
  startKey: z.string().default(""),
  period: z.string().default("current"),
});

export const Route = createFileRoute("/expenses")({
  beforeLoad: async ({ location }) => {
    if (!isAuth()) {
      throw redirect({
        to: "/login",
        search: {
          redirect: location.href,
        },
      });
    }
  },
  validateSearch: zodSearchValidator(expensesSearchSchema),
  component: Index,
});

function Index() {
  return <Home />;
}
