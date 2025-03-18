import {
  BackgroundRefetchErrorSnackbar,
  Button,
  CircularProgress,
  Container,
  Navbar,
  PageTitle,
} from "../components";
import { useGetSavingGoal } from "../queries";
import { useParams } from "@tanstack/react-router";
import { Error } from "./Error.tsx";
import { IconButton, Typography } from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import { SavingGoal } from "../types";
import { faTrash } from "@fortawesome/free-solid-svg-icons/faTrash";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faPencil } from "@fortawesome/free-solid-svg-icons";

export function SavingGoalDetail() {
  // @ts-expect-error ...
  const { savingGoalId } = useParams({ strict: false });

  const getSavingGoalQuery = useGetSavingGoal(savingGoalId);
  const savingGoal: SavingGoal | undefined = getSavingGoalQuery.data;

  if (getSavingGoalQuery.isError || !savingGoal) {
    return <Error />;
  }

  return (
    <Container>
      <Navbar />
      <BackgroundRefetchErrorSnackbar show={getSavingGoalQuery.isRefetching} />

      <PageTitle>Saving goal breakdown</PageTitle>

      <Grid container>
        {/*Goal detail card*/}
        <Grid xs={7}>
          <div className={"paper p-4"}>
            {/*Title*/}
            <div className={"flex items-center w-full"}>
              <Typography variant={"h5"} sx={{ fontWeight: "bold" }}>
                {savingGoal.name}
              </Typography>

              <IconButton
                sx={{ marginLeft: "auto", marginRight: "5px" }}
                title={"Delete"}
              >
                <FontAwesomeIcon icon={faTrash} />
              </IconButton>

              <Button
                variant={"outlined"}
                startIcon={<FontAwesomeIcon icon={faPencil} />}
              >
                Edit
              </Button>
            </div>

            {/*Percentage graphic*/}
            <div className={"flex w-full items-center justify-center"}>
              <CircularProgress
                // progress={savingGoal.progress}
                // target={savingGoal.target}
                progress={87}
                target={108}
                size={"8rem"}
                subtitle={"So far"}
              />
            </div>
          </div>
        </Grid>

        {/*Automatic savings*/}
        <Grid xs={5}></Grid>
      </Grid>
    </Container>
  );
}
