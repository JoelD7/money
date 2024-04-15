import { createFileRoute } from "@tanstack/react-router";
import { Home } from "../pages";
import { NavbarLayout } from "../components";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index() {
  return (
    <NavbarLayout>
      <Home />
    </NavbarLayout>
  );
}
