import { BackgroundRefetchErrorSnackbar, Container, Navbar } from "../components";
import { useGetSavingGoal } from "../queries";
import { useParams } from "@tanstack/react-router";
import { Error } from "./Error.tsx";

export function SavingGoalDetail() {
  // @ts-expect-error ...
  const { savingGoalId } = useParams({ strict: false });

  const getSavingGoalQuery = useGetSavingGoal(savingGoalId);

  if (getSavingGoalQuery.isError) {
    return <Error />;
  }

  return (
    <Container>
      <Navbar />
      <BackgroundRefetchErrorSnackbar show={getSavingGoalQuery.isRefetching} />
    </Container>
  );
}
