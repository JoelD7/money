import {
  BackgroundRefetchErrorSnackbar,
  Button,
  CircularProgress,
  Container,
  FontAwesomeIcon,
  Navbar,
  PageTitle,
} from "../components";
import { useGetSavingGoal } from "../queries";
import { useParams } from "@tanstack/react-router";
import { Error } from "./Error.tsx";
import { IconButton, TextField, Typography } from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import { SavingGoal } from "../types";
import { faTrash } from "@fortawesome/free-solid-svg-icons/faTrash";
import {
  faBullseye,
  faCalendar,
  faCircleInfo,
  faClock,
  faPencil,
} from "@fortawesome/free-solid-svg-icons";
import { currencyFormatter, tableDateFormatter } from "../utils";
import { useState } from "react";

export function SavingGoalDetail() {
  // @ts-expect-error ...
  const { savingGoalId } = useParams({ strict: false });

  const [recurringAmount, setRecurringAmount] = useState<number>(100);

  const monthYearFormatter = new Intl.DateTimeFormat("en-US", {
    year: "numeric",
    month: "long",
  });

  const getSavingGoalQuery = useGetSavingGoal(savingGoalId);
  const savingGoal: SavingGoal | undefined = getSavingGoalQuery.data;

  if (getSavingGoalQuery.isError || !savingGoal) {
    return <Error />;
  }

  function getReestimatedDeadline(): string {
    if (!savingGoal) return "";

    const periodsToReachGoal = Math.ceil(
      (savingGoal.target - savingGoal.progress) / recurringAmount,
    );

    const newDeadline = new Date(Date.now());
    newDeadline.setMonth(newDeadline.getMonth() + periodsToReachGoal);

    return monthYearFormatter.format(newDeadline);
  }

  function getReestimatedSavingAmount(): number {
    if (!savingGoal) return 0;

    const dateDiffInMs = new Date(savingGoal.deadline).getTime() - Date.now();
    const monthsUntilDeadline = Math.floor(dateDiffInMs / (1000 * 60 * 60 * 24 * 30));

    return (savingGoal.target - savingGoal.progress) / (monthsUntilDeadline - 1);
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
                    {tableDateFormatter.format(new Date(savingGoal.deadline))}
                  </h4>
                </div>
              </div>
            </div>
          </div>
        </Grid>

        {/*Automatic savings*/}
        <Grid xs={5}>
          <div className={"paper p-4"}>
            <Typography variant={"h5"} sx={{ fontWeight: "bold" }}>
              Automatic savings
            </Typography>

            <TextField
              margin={"normal"}
              name={"amount"}
              value={recurringAmount || ""}
              type={"number"}
              label={"Amount"}
              variant={"outlined"}
              required
              sx={{
                width: "50%",
              }}
              onChange={(e) => setRecurringAmount(Number(e.target.value))}
            />

            {/* Estimation box */}
            <div className={"rounded-xl bg-gray-200 mt-2 p-4"}>
              <Grid container spacing={1} height={"100%"}>
                <Grid xs={1}>
                  <div className={"flex h-full items-center"}>
                    <FontAwesomeIcon
                      colorClassName={"text-sky-600"}
                      icon={faCircleInfo}
                      size={"2xl"}
                    />
                  </div>
                </Grid>

                <Grid xs={11}>
                  <Typography variant={"body1"}>
                    By saving{" "}
                    <span className={"text-green-300"}>
                      {currencyFormatter.format(recurringAmount)}
                    </span>{" "}
                    each month, the deadline will move to {getReestimatedDeadline()}
                  </Typography>

                  <Typography variant={"body1"} sx={{ fontWeight: "bold" }}>
                    Or...
                  </Typography>

                  <Typography variant={"body1"}>
                    You would now have to save{" "}
                    <span className={"text-green-300"}>
                      {currencyFormatter.format(getReestimatedSavingAmount())}
                    </span>{" "}
                    each month to meet the same deadline of{" "}
                    {monthYearFormatter.format(new Date(savingGoal.deadline))}
                  </Typography>
                </Grid>
              </Grid>
            </div>
          </div>
        </Grid>
      </Grid>
    </Container>
  );
}
