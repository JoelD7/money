import {
  Alert,
  AlertTitle,
  capitalize,
  CircularProgress as MuiCircularProgress,
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
import { currencyFormatter, monthYearFormatter, utils } from "../../utils";
import { ChangeEvent, useEffect, useRef, useState } from "react";
import { SavingGoal, SnackAlert } from "../../types";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import api from "../../api";
import { savingGoalKeys, useGetSavingGoal } from "../../queries/saving_goals.ts";

type RecurringSavingProps = {
  savingGoalID: string;
};

export function RecurringSaving({ savingGoalID }: RecurringSavingProps) {
  const containerClasses = "flex flex-col paper p-4 h-full";
  const getSavingGoalQuery = useGetSavingGoal(savingGoalID);
  const savingGoal: SavingGoal | undefined = getSavingGoalQuery.data;

  // Default value is 1 to avoid posible division by zero
  const [recurringAmount, setRecurringAmount] = useState<number | undefined>(
    savingGoal ? savingGoal.recurring_amount : 0,
  );
  // This is a copy of recurringAmount except when it's 0. In that case holds the previous value before the change.
  const recurringAmountRef = useRef<number>(recurringAmount || 1);
  const [toggleEditView, setToggleEditView] = useState<boolean>(false);
  const [alert, setAlert] = useState<SnackAlert>({
    open: false,
    type: "success",
    title: "",
  });

  useEffect(() => {
    //This is needed so that recurringAmount is initialized with the correct value when the component is mounted again
    // after the query returns. This should only happen once.
    if (
      savingGoal &&
      savingGoal.recurring_amount !== undefined &&
      recurringAmount === 0
    ) {
      setRecurringAmount(savingGoal.recurring_amount);
      recurringAmountRef.current = savingGoal.recurring_amount;
    }
  }, [savingGoal, recurringAmount]);

  //Always use recurringAmountRef to avoid possible division by zero
  const reestimatedDeadline: Date = utils.estimateDeadlineFromRecurringAmount(
    recurringAmountRef.current,
    savingGoal,
  );
  const reestimatedSavingAmount: number = utils.estimateSavingAmount(savingGoal);

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
        .invalidateQueries({ queryKey: savingGoalKeys.single(savingGoalID) })
        .then(() => {
          // Close the edit view AFTER the query is updated so that the correct info box is rendered
          setToggleEditView(false);
        })
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
    // This should never happen because I'm rendering a loading or error component if this is undefined. I do this check
    // to make TS happy
    if (!savingGoal) {
      return <></>;
    }

    if (!savingGoal.is_recurring) return <InfoBoxNoRecurringSaving />;

    const reestimatedDeadlineString = monthYearFormatter.format(reestimatedDeadline);
    const deadlineString = monthYearFormatter.format(new Date(savingGoal.deadline));

    const props: InfoBoxProps = {
      reestimatedDeadline: reestimatedDeadline,
      recurringAmount: recurringAmount || 0,
      reestimatedSavingAmount: reestimatedSavingAmount,
      savingGoalRecurringAmount: savingGoal.recurring_amount
        ? savingGoal.recurring_amount
        : 0,
      deadline: new Date(savingGoal.deadline),
      onChangeRecurringAmount: handleChangeRecurringAmount,
      onChangeDeadline: handleChangeDeadline,
    };

    if (
      reestimatedDeadlineString === deadlineString &&
      toggleEditView &&
      recurringAmount !== savingGoal.recurring_amount
    ) {
      return <InfoBoxMeetDeadlineAmount {...props} />;
    }

    if (reestimatedDeadlineString === deadlineString) {
      return <InfoBoxKeepItUp {...props} />;
    }

    if (toggleEditView) {
      return <InfoBoxEdit {...props} />;
    }

    return <InfoBoxBehind {...props} />;
  }

  function handleChangeRecurringAmount(newAmount: number) {
    if (!savingGoal) return;

    mutateSavingGoal.mutate({
      ...savingGoal,
      recurring_amount: newAmount,
    });

    setToggleEditView(false);
    setRecurringAmount(newAmount);
    recurringAmountRef.current = newAmount;
  }

  function handleChangeDeadline() {
    if (!savingGoal) return;

    mutateSavingGoal.mutate({
      ...savingGoal,
      recurring_amount: recurringAmount,
      deadline: reestimatedDeadline.toISOString(),
    });
  }

  function handleRecurringAmountChange(event: ChangeEvent<HTMLInputElement>) {
    if (!savingGoal) return;

    const value = Number(event.target.value);
    setRecurringAmount(value);
    if (value > 0) {
      recurringAmountRef.current = value;
    }

    if (value !== savingGoal.recurring_amount) {
      setToggleEditView(true);
    }
  }

  if (getSavingGoalQuery.isPending || savingGoal === undefined) {
    return (
      <div className={`${containerClasses} items-center justify-center`}>
        <MuiCircularProgress size={"7rem"} />
      </div>
    );
  }

  return (
    <div className={containerClasses}>
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
  savingGoalRecurringAmount: number;
  reestimatedSavingAmount: number;
  deadline: Date;
  onChangeRecurringAmount: (amount: number) => void;
  onChangeDeadline?: () => void;
  loadingButton?: boolean;
};

function InfoBoxKeepItUp(props: InfoBoxProps) {
  const { deadline } = props;
  return (
    <Grid container spacing={4} height={"100%"}>
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
              onClick={() => onChangeRecurringAmount(reestimatedSavingAmount)}
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
    savingGoalRecurringAmount,
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

        {savingGoalRecurringAmount !== reestimatedSavingAmount && (
          <>
            <Typography variant={"body1"} sx={{ fontWeight: "bold" }}>
              Or...
            </Typography>

            <Typography variant={"body1"}>
              Save at least{" "}
              <span className={"text-green-300"}>
                {currencyFormatter.format(reestimatedSavingAmount)}
              </span>{" "}
              each month to meet the same deadline of{" "}
              {monthYearFormatter.format(new Date(deadline))}
            </Typography>
          </>
        )}
      </Grid>

      {/*Buttons*/}
      <Grid xs={12}>
        <div className={"flex justify-end gap-1"}>
          {savingGoalRecurringAmount !== reestimatedSavingAmount && (
            <Tooltip title={`Changes the recurring amount to ${reestimatedSavingAmount}`}>
              <Button
                variant={"contained"}
                loading={loadingButton}
                onClick={() => onChangeRecurringAmount(reestimatedSavingAmount)}
              >{`Accept ${currencyFormatter.format(reestimatedSavingAmount)}`}</Button>
            </Tooltip>
          )}

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
  const { recurringAmount, onChangeRecurringAmount, loadingButton } = props;
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
          <Tooltip title={`Changes the recurring amount to ${recurringAmount}`}>
            <Button
              variant={"contained"}
              loading={loadingButton}
              onClick={() => onChangeRecurringAmount(recurringAmount)}
            >{`Save ${currencyFormatter.format(recurringAmount)}`}</Button>
          </Tooltip>
        </div>
      </Grid>
    </Grid>
  );
}
