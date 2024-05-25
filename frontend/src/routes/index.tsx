import { createFileRoute, redirect } from "@tanstack/react-router";
import { Home } from "../pages";
import { store } from "../store";

function isAuth() {
  return store.getState().authReducer.isAuthenticated;
}

export const Route = createFileRoute("/")({
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
  component: Index,
});

function Index() {
  return <Home />;
}
