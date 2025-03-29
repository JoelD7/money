import {
  Alert,
  AlertTitle,
  capitalize,
  Snackbar,
  TextField,
  Tooltip,
  Typography,
} from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import { Button, FontAwesomeIcon } from "../atoms";
import {
  faCircleCheck,
  faCircleInfo,
  faTriangleExclamation,
} from "@fortawesome/free-solid-svg-icons";
import { currencyFormatter, monthYearFormatter } from "../../utils";
import { ChangeEvent, useState } from "react";
import { SavingGoal, SnackAlert } from "../../types";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import api from "../../api";
import { savingGoalKeys, useGetSavingGoal } from "../../queries/saving_goals.ts";

type RecurringSavingProps = {
  savingGoalID: string;
};

export function RecurringSaving(props: RecurringSavingProps) {
  // const { savingGoal } = props;
  const savingGoalId = props.savingGoalID;

  const getSavingGoalQuery = useGetSavingGoal(savingGoalId);
  const savingGoal: SavingGoal = getSavingGoalQuery.data ? getSavingGoalQuery.data : {};

  // Default value is 1 to avoid posible division by zero
  const [recurringAmount, setRecurringAmount] = useState<number>(
    savingGoal.recurring_amount ? savingGoal.recurring_amount : 1,
  );
  const [toggleEditView, setToggleEditView] = useState<boolean>(false);
  const [alert, setAlert] = useState<SnackAlert>({
    open: false,
    type: "success",
    title: "",
  });

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

    // Always use the first day of the month for both dates to avoid rounding errors when the cur date is near the
    // end of the month and the deadline is near the start of the month
    const deadlineMonth = new Date(savingGoal.deadline).getMonth();
    const deadlineYear = new Date(savingGoal.deadline).getFullYear();
    const deadline = new Date(deadlineYear, deadlineMonth, 1);
    const currentMonth = new Date().getMonth();
    const currentYear = new Date().getFullYear();
    const current = new Date(currentYear, currentMonth, 1);

    const monthsUntilDeadline = Math.floor(
      (deadline.getTime() - current.getTime()) / (1000 * 60 * 60 * 24 * 30),
    );

    const result = (savingGoal.target - savingGoal.progress) / monthsUntilDeadline;
    return Math.ceil(result * 10) / 10;
  })();

  const queryClient = useQueryClient();
  const mutateSavingGoal = useMutation({
    mutationFn: api.updateSavingGoal,
    onSuccess: () => {
      setAlert({
        ...alert,
        open: true,
        type: "success",
        title: "Saving goal updated successfully",
      });

      queryClient
        .invalidateQueries({ queryKey: savingGoalKeys.single(savingGoal.saving_goal_id) })
        .then()
        .catch((e) => {
          console.error("Error invalidating saving goal query", e);
        });
    },
    onError: () => {
      setAlert({
        ...alert,
        open: true,
        type: "error",
        title: "Error updating saving goal",
      });
    },
  });

  function renderInfoBox() {
    if (!savingGoal.is_recurring) return <InfoBoxNoRecurringSaving />;

    const reestimatedDeadlineString = monthYearFormatter.format(reestimatedDeadline);
    const deadlineString = monthYearFormatter.format(new Date(savingGoal.deadline));

    const props: InfoBoxProps = {
      reestimatedDeadline: reestimatedDeadline,
      recurringAmount: recurringAmount,
      reestimatedSavingAmount: reestimatedSavingAmount,
      deadline: new Date(savingGoal.deadline),
      onChangeRecurringAmount: handleChangeRecurringAmount,
      onChangeDeadline: handleChangeDeadline,
    };

    if (
      reestimatedDeadlineString === deadlineString &&
      savingGoal.recurring_amount === recurringAmount
    ) {
      return <InfoBoxKeepItUp {...props} />;
    }

    if (reestimatedDeadlineString === deadlineString) {
      return <InfoBoxMeetDeadlineAmount {...props} />;
    }

    if (toggleEditView) {
      return <InfoBoxEdit {...props} />;
    }

    return <InfoBoxBehind {...props} />;
  }

  function handleRecurringAmountChange(event: ChangeEvent<HTMLInputElement>) {
    const value = Number(event.target.value);
    setRecurringAmount(value);

    if (value !== savingGoal.recurring_amount) {
      setToggleEditView(true);
    }
  }

  function handleChangeRecurringAmount() {
    setToggleEditView(false);

    mutateSavingGoal.mutate({
      ...savingGoal,
      recurring_amount: reestimatedSavingAmount,
    });
  }

  function handleChangeDeadline() {
    setToggleEditView(false);

    mutateSavingGoal.mutate({
      ...savingGoal,
      deadline: reestimatedDeadline.toISOString(),
    });
  }

  if (getSavingGoalQuery.isPending) {
    return (
      <div className={"paper p-4"}>
        <Typography variant={"h5"} sx={{ fontWeight: "bold" }}>
          Loading...
        </Typography>
      </div>
    );
  }

  if (getSavingGoalQuery.isError) {
    return (
      <div className={"paper p-4"}>
        <Typography variant={"h5"} sx={{ fontWeight: "bold" }}>
          Error loading saving goal...
        </Typography>
      </div>
    );
  }

  return (
    <div className={"flex flex-col paper p-4 h-full"}>
      {/*Alert */}
      <Snackbar
        open={alert.open}
        onClose={() => setAlert({ ...alert, open: false })}
        autoHideDuration={6000}
        anchorOrigin={{ vertical: "top", horizontal: "right" }}
      >
        <Alert variant={"filled"} severity={alert.type}>
          <AlertTitle>{capitalize(alert.type)}</AlertTitle>
          {alert.title}
        </Alert>
      </Snackbar>

      <div className={"flex justify-between w-full"}>
        <Typography variant={"h5"} sx={{ fontWeight: "bold" }}>
          Automatic savings
        </Typography>
      </div>

      {savingGoal.is_recurring && (
        <div className={"flex w-full"}>
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
            onChange={handleRecurringAmountChange}
          />
        </div>
      )}

      {/* Estimation box */}
      <div className={"flex w-full m-auto rounded-xl bg-gray-100 p-4"}>
        {renderInfoBox()}
      </div>
    </div>
  );
}

type InfoBoxProps = {
  reestimatedDeadline: Date;
  recurringAmount: number;
  reestimatedSavingAmount: number;
  deadline: Date;
  onChangeRecurringAmount?: () => void;
  onChangeDeadline?: () => void;
  loadingButton?: boolean;
};

function InfoBoxKeepItUp(props: InfoBoxProps) {
  const { deadline } = props;
  return (
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
  );
}

function InfoBoxNoRecurringSaving() {
  return (
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
          automatically create a new savings entry for this goal with a fixed amount. This
          way, you won’t have to manually add one every month.
        </Typography>
      </Grid>

      {/*Buttons*/}
      <Grid xs={12}>
        <div className={"flex justify-end gap-1"}>
          <Button variant={"contained"}>Set up recurring savings</Button>
        </div>
      </Grid>
    </Grid>
  );
}

function InfoBoxBehind(props: InfoBoxProps) {
  const {
    reestimatedDeadline,
    reestimatedSavingAmount,
    onChangeDeadline,
    onChangeRecurringAmount,
    loadingButton,
  } = props;
  return (
    <Grid container spacing={1} height={"100%"}>
      {/*Information icon*/}
      <Grid xs={1}>
        <div className={"flex h-full items-center"}>
          <FontAwesomeIcon
            colorClassName={"text-yellow-400"}
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
          Change the recurring amount to{" "}
          <span className={"text-green-300"}>
            {currencyFormatter.format(reestimatedSavingAmount)}
          </span>{" "}
          to meet the deadline.
        </Typography>

        <Typography variant={"body1"} sx={{ fontWeight: "bold" }}>
          Or...
        </Typography>

        <Typography variant={"body1"}>
          Change the deadline to {monthYearFormatter.format(reestimatedDeadline)}
        </Typography>
      </Grid>

      {/*Buttons*/}
      <Grid xs={12}>
        <div className={"flex justify-end gap-1"}>
          <Tooltip title={`Changes the recurring amount to ${reestimatedSavingAmount}`}>
            <Button
              variant={"contained"}
              loading={loadingButton}
              onClick={onChangeRecurringAmount}
            >{`Accept ${currencyFormatter.format(reestimatedSavingAmount)}`}</Button>
          </Tooltip>

          <Tooltip
            title={`Changes the deadline to ${monthYearFormatter.format(reestimatedDeadline)}`}
          >
            <Button
              variant={"outlined"}
              loading={loadingButton}
              onClick={onChangeDeadline}
            >{`Save & Change deadline`}</Button>
          </Tooltip>
        </div>
      </Grid>
    </Grid>
  );
}

function InfoBoxEdit(props: InfoBoxProps) {
  const {
    reestimatedDeadline,
    recurringAmount,
    reestimatedSavingAmount,
    deadline,
    onChangeRecurringAmount,
    loadingButton,
    onChangeDeadline,
  } = props;
  return (
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
          Save{" "}
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
              loading={loadingButton}
              onClick={onChangeRecurringAmount}
            >{`Accept ${currencyFormatter.format(reestimatedSavingAmount)}`}</Button>
          </Tooltip>

          <Tooltip
            title={`Changes the deadline to ${monthYearFormatter.format(reestimatedDeadline)}`}
          >
            <Button
              onClick={onChangeDeadline}
              loading={loadingButton}
              variant={"contained"}
            >{`Save & Change deadline`}</Button>
          </Tooltip>
        </div>
      </Grid>
    </Grid>
  );
}

function InfoBoxMeetDeadlineAmount(props: InfoBoxProps) {
  const {
    recurringAmount,
    reestimatedSavingAmount,
    onChangeRecurringAmount,
    loadingButton,
  } = props;
  return (
    <Grid container spacing={1} height={"100%"}>
      {/*Information icon*/}
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
        <Typography variant={"body1"}>
          By saving{" "}
          <span className={"text-green-300"}>
            {currencyFormatter.format(recurringAmount)}
          </span>{" "}
          each month you will meet the deadline!
        </Typography>
      </Grid>

      {/*Buttons*/}
      <Grid xs={12}>
        <div className={"flex justify-end gap-1"}>
          <Tooltip title={`Changes the recurring amount to ${reestimatedSavingAmount}`}>
            <Button
              variant={"contained"}
              loading={loadingButton}
              onClick={onChangeRecurringAmount}
            >{`Save ${currencyFormatter.format(reestimatedSavingAmount)}`}</Button>
          </Tooltip>
        </div>
      </Grid>
    </Grid>
  );
}
