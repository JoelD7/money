import { createFileRoute } from "@tanstack/react-router";
import { Savings } from "../pages";

export const Route = createFileRoute("/savings/")({
  component: Savings,
});
