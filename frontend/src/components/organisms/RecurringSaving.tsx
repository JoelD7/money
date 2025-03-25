import { TextField, Tooltip, Typography } from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import { Button, FontAwesomeIcon } from "../atoms";
import {
  faCircleCheck,
  faCircleInfo,
  faPencil,
  faTriangleExclamation,
} from "@fortawesome/free-solid-svg-icons";
import { currencyFormatter, monthYearFormatter } from "../../utils";
import { useState } from "react";
import { SavingGoal } from "../../types";

type RecurringSavingProps = {
  savingGoal: SavingGoal;
};

const infoBoxContainerClass = "rounded-xl bg-gray-100 mt-2 p-4";

export function RecurringSaving(props: RecurringSavingProps) {
  const { savingGoal } = props;

  // Default value is 1 to avoid posible division by zero
  const [recurringAmount, setRecurringAmount] = useState<number>(
    savingGoal.recurring_amount ? savingGoal.recurring_amount : 1,
  );

  const reestimatedDeadline: Date = (() => {
    const periodsToReachGoal = Math.ceil(
      (savingGoal.target - savingGoal.progress) / recurringAmount,
    );

    const newDeadline = new Date(Date.now());
    newDeadline.setMonth(newDeadline.getMonth() + periodsToReachGoal);

    return newDeadline;
  })();

  const reestimatedSavingAmount: number = (() => {
    if (!savingGoal) return 0;

    const dateDiffInMs = new Date(savingGoal.deadline).getTime() - Date.now();
    const monthsUntilDeadline = Math.floor(dateDiffInMs / (1000 * 60 * 60 * 24 * 30));

    return (savingGoal.target - savingGoal.progress) / (monthsUntilDeadline - 1);
  })();

  function renderInfoBox() {
    if (!savingGoal.is_recurring) return <InfoBoxNoRecurringSaving />;

    const reestimatedDeadlineString = monthYearFormatter.format(reestimatedDeadline);
    const deadlineString = monthYearFormatter.format(new Date(savingGoal.deadline));

    const props: InfoBoxProps = {
      reestimatedDeadline: reestimatedDeadline,
      recurringAmount: recurringAmount,
      reestimatedSavingAmount: reestimatedSavingAmount,
      deadline: new Date(savingGoal.deadline),
    };

    if (reestimatedDeadlineString === deadlineString) {
      return <InfoBoxKeepItUp {...props} />;
    }

    return <InfoBoxDefault {...props} />;
  }

  return (
    <div className={"paper p-4"}>
      <div className={"flex justify-between"}>
        <Typography variant={"h5"} sx={{ fontWeight: "bold" }}>
          Automatic savings
        </Typography>

        {savingGoal.is_recurring && (
          <Button variant={"outlined"} startIcon={<FontAwesomeIcon icon={faPencil} />}>
            Edit
          </Button>
        )}
      </div>

      {savingGoal.is_recurring && (
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
      )}

      {/* Estimation box */}
      {renderInfoBox()}
    </div>
  );
}

type InfoBoxProps = {
  reestimatedDeadline: Date;
  recurringAmount: number;
  reestimatedSavingAmount: number;
  deadline: Date;
};

function InfoBoxKeepItUp(props: InfoBoxProps) {
  const { deadline } = props;
  return (
    <div className={infoBoxContainerClass}>
      <Grid container spacing={1} height={"100%"}>
        {/*Checkmark icon*/}
        <Grid xs={1}>
          <div className={"flex h-full items-center"}>
            <FontAwesomeIcon
              colorClassName={"text-green-300"}
              icon={faCircleCheck}
              size={"2xl"}
            />
          </div>
        </Grid>

        {/*Explanatory text*/}
        <Grid xs={11}>
          <Typography variant={"h5"} sx={{ fontWeight: "bold" }}>
            Keep it up!
          </Typography>

          <Typography variant={"body1"}>
            You are on track to reach you goal by{" "}
            {monthYearFormatter.format(new Date(deadline))}
          </Typography>
        </Grid>
      </Grid>
    </div>
  );
}

function InfoBoxNoRecurringSaving() {
  return (
    <div className={infoBoxContainerClass}>
      <Grid container spacing={1} height={"100%"}>
        {/*Warning icon*/}
        <Grid xs={1}>
          <div className={"flex h-full items-center"}>
            <FontAwesomeIcon
              colorClassName={"text-blue-200"}
              icon={faCircleInfo}
              size={"2xl"}
            />
          </div>
        </Grid>

        {/*Explanatory text*/}
        <Grid xs={11}>
          <Typography variant={"body1"}>
            When you set up recurring savings, at the start of a new period the app will
            automatically create a new savings entry for this goal with a fixed amount.
            This way, you won’t have to manually add one every month.
          </Typography>
        </Grid>

        {/*Buttons*/}
        <Grid xs={12}>
          <div className={"flex justify-end gap-1"}>
            <Button variant={"contained"}>Set up recurring savings</Button>
          </div>
        </Grid>
      </Grid>
    </div>
  );
}

function InfoBoxBehind(props: InfoBoxProps) {
  const { reestimatedDeadline, reestimatedSavingAmount } = props;
  return (
    <div className={infoBoxContainerClass}>
      <Grid container spacing={1} height={"100%"}>
        {/*Information icon*/}
        <Grid xs={1}>
          <div className={"flex h-full items-center"}>
            <FontAwesomeIcon
              colorClassName={"text-yellow-200"}
              icon={faTriangleExclamation}
              size={"2xl"}
            />
          </div>
        </Grid>

        {/*Explanatory text*/}
        <Grid xs={11}>
          <Typography variant={"h5"} sx={{ fontWeight: "bold" }}>
            You're behind...
          </Typography>

          <Typography variant={"body1"}>
            Change the running amount to{" "}
            <span className={"text-green-300"}>
              {currencyFormatter.format(reestimatedSavingAmount)}
            </span>{" "}
            to meet the deadline.
          </Typography>

          <Typography variant={"body1"} sx={{ fontWeight: "bold" }}>
            Or...
          </Typography>

          <Typography variant={"body1"}>
            Change the deadline to
            {monthYearFormatter.format(reestimatedDeadline)}
          </Typography>
        </Grid>

        {/*Buttons*/}
        <Grid xs={12}>
          <div className={"flex justify-end gap-1"}>
            <Tooltip title={`Changes the recurring amount to ${reestimatedSavingAmount}`}>
              <Button
                variant={"contained"}
              >{`Accept ${currencyFormatter.format(reestimatedSavingAmount)}`}</Button>
            </Tooltip>

            <Tooltip
              title={`Changes the deadline to ${monthYearFormatter.format(reestimatedDeadline)}`}
            >
              <Button variant={"outlined"}>{`Change deadline`}</Button>
            </Tooltip>
          </div>
        </Grid>
      </Grid>
    </div>
  );
}

function InfoBoxDefault(props: InfoBoxProps) {
  const { reestimatedDeadline, recurringAmount, reestimatedSavingAmount, deadline } =
    props;
  return (
    <div className={infoBoxContainerClass}>
      <Grid container spacing={1} height={"100%"}>
        {/*Information icon*/}
        <Grid xs={1}>
          <div className={"flex h-full items-center"}>
            <FontAwesomeIcon
              colorClassName={"text-sky-600"}
              icon={faCircleInfo}
              size={"2xl"}
            />
          </div>
        </Grid>

        {/*Explanatory text*/}
        <Grid xs={11}>
          <Typography variant={"body1"}>
            By saving{" "}
            <span className={"text-green-300"}>
              {currencyFormatter.format(recurringAmount)}
            </span>{" "}
            each month, the deadline will move to{" "}
            {monthYearFormatter.format(reestimatedDeadline)}
          </Typography>

          <Typography variant={"body1"} sx={{ fontWeight: "bold" }}>
            Or...
          </Typography>

          <Typography variant={"body1"}>
            You would now have to save{" "}
            <span className={"text-green-300"}>
              {currencyFormatter.format(reestimatedSavingAmount)}
            </span>{" "}
            each month to meet the same deadline of{" "}
            {monthYearFormatter.format(new Date(deadline))}
          </Typography>
        </Grid>

        {/*Buttons*/}
        <Grid xs={12}>
          <div className={"flex justify-end gap-1"}>
            <Tooltip title={`Changes the recurring amount to ${reestimatedSavingAmount}`}>
              <Button
                variant={"contained"}
              >{`Accept ${currencyFormatter.format(reestimatedSavingAmount)}`}</Button>
            </Tooltip>

            <Tooltip
              title={`Changes the deadline to ${monthYearFormatter.format(reestimatedDeadline)}`}
            >
              <Button variant={"contained"}>{`Change deadline`}</Button>
            </Tooltip>
          </div>
        </Grid>
      </Grid>
    </div>
  );
}
