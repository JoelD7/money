import {
  BackgroundRefetchErrorSnackbar,
  Button,
  Container,
  GoalDetailCard,
  Navbar,
  NewSaving,
  PageTitle,
  RecurringSaving,
  SavingsTable,
  Snackbar,
} from "../components";
import { useGetSavingGoal } from "../queries";
import { useNavigate, useParams } from "@tanstack/react-router";
import Grid from "@mui/material/Unstable_Grid2";
import { SavingGoal, SnackAlert } from "../types";
import { Typography } from "@mui/material";
import { useState } from "react";
import { NotFound } from "./NotFound.tsx";
import { Error } from "./Error.tsx";

export function SavingGoalDetail() {
  // @ts-expect-error ...
  const { savingGoalId } = useParams({ strict: false });

  const [openNewSaving, setOpenNewSaving] = useState(false);
  const [key, setKey] = useState<number>(0);
  const [alert, setAlert] = useState<SnackAlert>({
    open: false,
    type: "success",
    title: "",
  });

  const getSavingGoalQuery = useGetSavingGoal(savingGoalId);
  const savingGoal: SavingGoal | undefined = getSavingGoalQuery.data;

  const navigate = useNavigate();

  if (getSavingGoalQuery.error?.response?.status === 404) {
    return (
      <Container>
        <Navbar />
        <NotFound
          title={"Whoops..."}
          subtitle={"Saving goal not found"}
          body={"This saving goal has likely been deleted or it never existed."}
        >
          <Button
            variant={"contained"}
            sx={{
              marginTop: "10px",
              fontSize: "18px",
            }}
            onClick={() => {
              navigate({
                to: "/savings",
              }).catch((e) => console.error("Error navigating to /savings", e));
            }}
          >
            Go back to savings
          </Button>
        </NotFound>
      </Container>
    );
  }

  if (getSavingGoalQuery.isError && !savingGoal) {
    return (
      <Container>
        <Navbar />
        <Error />
      </Container>
    );
  }

  function handleAlert(alert?: SnackAlert) {
    if (alert) {
      setAlert(alert);
    }
  }

  return (
    <Container>
      <Navbar />
      <BackgroundRefetchErrorSnackbar show={getSavingGoalQuery.isRefetching} />

      <Snackbar
        open={alert.open}
        title={alert.title}
        message={alert.message}
        severity={alert.type}
        onClose={() => setAlert({ ...alert, open: false })}
      />

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

        <Grid xs={12}>
          <div className={"pt-4"}>
            <Typography variant={"h4"}>Latest savings</Typography>
            <div className={"pt-4 pb-2"}>
              <Button variant={"contained"} onClick={() => setOpenNewSaving(true)}>
                New saving
              </Button>
            </div>
            <SavingsTable savingGoalID={savingGoalId} />
          </div>
        </Grid>
      </Grid>

      {savingGoal && (
        <NewSaving
          key={key}
          open={openNewSaving}
          savingGoalId={savingGoal.saving_goal_id}
          onClose={() => {
            setOpenNewSaving(false);
            setKey(key + 1);
          }}
          onAlert={handleAlert}
        />
      )}
    </Container>
  );
}
