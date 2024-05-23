import { createFileRoute, redirect } from "@tanstack/react-router";
import { Home } from "../pages";
import { NavbarLayout } from "../components";
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
  return (
    <NavbarLayout>
      <Home />
    </NavbarLayout>
  );
}
