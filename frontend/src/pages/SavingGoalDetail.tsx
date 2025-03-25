import {
  BackgroundRefetchErrorSnackbar,
  Button,
  CircularProgress,
  Container,
  FontAwesomeIcon,
  Navbar,
  PageTitle,
  RecurringSaving,
} from "../components";
import { useGetSavingGoal } from "../queries";
import { useParams } from "@tanstack/react-router";
import { Error } from "./Error.tsx";
import { IconButton, Typography } from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import { SavingGoal } from "../types";
import { faTrash } from "@fortawesome/free-solid-svg-icons/faTrash";
import {
  faBullseye,
  faCalendar,
  faClock,
  faPencil,
} from "@fortawesome/free-solid-svg-icons";
import { currencyFormatter, monthYearFormatter } from "../utils";

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

      <Grid container spacing={2}>
        {/*Goal detail card*/}
        <Grid xs={6}>
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
                progress={savingGoal.progress}
                target={savingGoal.target}
                size={"8rem"}
                subtitle={"Progress"}
              />
            </div>

            {/*Breakdown in numbers*/}
            <div className={"flex w-full justify-center pt-8"}>
              <div className={"grid grid-cols-3 justify-center w-[90%]"}>
                {/*Goal*/}
                <div>
                  <div className={"flex items-center"}>
                    <span className={"text-amber-400"}>
                      <FontAwesomeIcon icon={faBullseye} />
                    </span>
                    <h4 className={"ml-1 text-xl font-bold leading-[0px]"}>Goal</h4>
                  </div>
                  <h4 className={"text-xl"}>
                    {currencyFormatter.format(savingGoal.target)}
                  </h4>
                </div>

                {/*Progress*/}
                <div>
                  <div className={"flex items-center"}>
                    <FontAwesomeIcon colorClassName={"text-sky-600"} icon={faClock} />
                    <h4 className={"ml-1 text-xl font-bold leading-[0px]"}>Progress</h4>
                  </div>
                  <h4 className={"text-xl"}>
                    {currencyFormatter.format(savingGoal.progress)}
                  </h4>
                </div>

                {/*Deadline*/}
                <div>
                  <div className={"flex items-center"}>
                    <FontAwesomeIcon colorClassName={"text-red-200"} icon={faCalendar} />
                    <h4 className={"ml-1 text-xl font-bold leading-[0px]"}>Deadline</h4>
                  </div>
                  <h4 className={"text-xl"}>
                    {monthYearFormatter.format(new Date(savingGoal.deadline))}
                  </h4>
                </div>
              </div>
            </div>
          </div>
        </Grid>

        {/*Automatic savings*/}
        <Grid xs={5}>
          <RecurringSaving savingGoal={savingGoal} />
        </Grid>
      </Grid>
    </Container>
  );
}
