import {
  CircularProgress as MuiCircularProgress,
  IconButton,
  Typography,
} from "@mui/material";
import { Button, CircularProgress, FontAwesomeIcon } from "../atoms";
import { faTrash } from "@fortawesome/free-solid-svg-icons/faTrash";
import {
  faBullseye,
  faCalendar,
  faClock,
  faPencil,
} from "@fortawesome/free-solid-svg-icons";
import { currencyFormatter, monthYearFormatter } from "../../utils";
import { useParams } from "@tanstack/react-router";
import { useGetSavingGoal } from "../../queries";
import { SavingGoal } from "../../types";
import Grid from "@mui/material/Unstable_Grid2";

export function GoalDetailCard() {
  // @ts-expect-error ...
  const { savingGoalId } = useParams({ strict: false });
  const containerClasses = "paper p-4 h-full";

  const getSavingGoalQuery = useGetSavingGoal(savingGoalId);
  const savingGoal: SavingGoal | undefined = getSavingGoalQuery.data;

  if (getSavingGoalQuery.isPending || savingGoal === undefined) {
    return (
      <div className={containerClasses}>
        <MuiCircularProgress size={"7rem"} />
      </div>
    );
  }

  if (getSavingGoalQuery.isError && !savingGoal) {
    return (
      <div className={containerClasses}>
        <Typography variant={"h5"}>Error</Typography>
      </div>
    );
  }

  return (
    <div className={containerClasses}>
      <Grid container height={"100%"}>
        {/*Title*/}
        <Grid xs={12}>
          <div className={"flex items-center w-full h-fit"}>
            <Typography variant={"h5"} sx={{ fontWeight: "bold" }}>
              {savingGoal.name}
            </Typography>

            <IconButton sx={{ marginLeft: "auto", marginRight: "5px" }} title={"Delete"}>
              <FontAwesomeIcon icon={faTrash} />
            </IconButton>

            <Button variant={"outlined"} startIcon={<FontAwesomeIcon icon={faPencil} />}>
              Edit
            </Button>
          </div>
        </Grid>

        {/*Percentage graphic*/}
        <Grid xs={12}>
          <div className={"flex w-full items-center justify-center"}>
            <CircularProgress
              progress={savingGoal.progress}
              target={savingGoal.target}
              size={"8rem"}
              subtitle={"Progress"}
            />
          </div>
        </Grid>

        {/*Breakdown in numbers*/}
        <Grid xs={12} alignSelf={"end"}>
          <div className={"flex w-full justify-center self-end"}>
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
        </Grid>
      </Grid>
    </div>
  );
}
