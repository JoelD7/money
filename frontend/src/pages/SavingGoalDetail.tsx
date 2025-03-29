import {
  BackgroundRefetchErrorSnackbar,
  Container,
  GoalDetailCard,
  Navbar,
  PageTitle,
  RecurringSaving,
} from "../components";
import { useGetSavingGoal } from "../queries";
import { useParams } from "@tanstack/react-router";
import { Error } from "./Error.tsx";
import Grid from "@mui/material/Unstable_Grid2";
import { SavingGoal } from "../types";

export function SavingGoalDetail() {
  // @ts-expect-error ...
  const { savingGoalId } = useParams({ strict: false });

  const getSavingGoalQuery = useGetSavingGoal(savingGoalId);
  const savingGoal: SavingGoal | undefined = getSavingGoalQuery.data;

  if (getSavingGoalQuery.isError && !savingGoal) {
    return (
      <Container>
        <Navbar />
        <Error />
      </Container>
    );
  }

  return (
    <Container>
      <Navbar />
      <BackgroundRefetchErrorSnackbar show={getSavingGoalQuery.isRefetching} />

      <PageTitle>Saving goal breakdown</PageTitle>

      <Grid container spacing={2} minHeight={"24rem"}>
        {/*Goal detail card*/}
        <Grid xs={6} height={"24rem"}>
          <GoalDetailCard />
        </Grid>

        {/*Automatic savings*/}
        <Grid xs={5} height={"24rem"}>
          <RecurringSaving savingGoalID={savingGoalId} />
        </Grid>
      </Grid>
    </Container>
  );
}
