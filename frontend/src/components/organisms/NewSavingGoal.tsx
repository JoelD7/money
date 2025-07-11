import {
  Box,
  Divider,
  FormControlLabel,
  FormGroup,
  Switch,
  TextField,
  Typography,
} from "@mui/material";
import { FormEvent, useState } from "react";
import Grid from "@mui/material/Unstable_Grid2";
import { DatePicker } from "@mui/x-date-pickers";
import dayjs, { Dayjs } from "dayjs";
import { Button } from "../atoms";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import api from "../../api";
import { AxiosError } from "axios";
import { APIError, SavingGoal, SnackAlert } from "../../types";
import * as yup from "yup";
import { ValidationError } from "yup";
import { Dialog } from "../molecules";
import { savingGoalKeys } from "../../queries/saving_goals.ts";

type NewSavingGoalProps = {
  open: boolean;
  onClose: () => void;
  onAlert: (alert?: SnackAlert) => void;
};

export function NewSavingGoal({ open, onClose, onAlert }: NewSavingGoalProps) {
  const [name, setName] = useState<string>("");
  const [target, setTarget] = useState<number | null>(null);
  const [recurringAmount, setRecurringAmount] = useState<number | null>(null);
  const [recurringChecked, setRecurringChecked] = useState<boolean>(false);
  const [deadline, setDeadline] = useState<Dayjs | null>(dayjs());

  const queryClient = useQueryClient();

  const createSavingGoalMutation = useMutation({
    mutationFn: api.createSavingGoal,
    onSuccess: () => {
      onAlert({
        open: true,
        type: "success",
        title: "Saving goal created successfully",
      });
      onClose();
      queryClient
        .invalidateQueries({ queryKey: [...savingGoalKeys.all] })
        .then(() => {})
        .catch((e) => {
          console.error("Error invalidating saving goals query", e);
        });
    },
    onError: (error) => {
      if (error) {
        const err = error as AxiosError;
        const responseError = err.response?.data as APIError;
        onAlert({
          open: true,
          type: "error",
          title: responseError.message as string,
        });
      }
    },
  });

  const validationSchema = yup.object({
    name: yup.string().required("Name is required"),
    target: yup.number().required("Target is required").moreThan(0, "Target is required"),
    deadline: yup.date().required("Deadline is required"),
  });

  function createSavingGoal(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();

    let savingGoalDeadline: string = "";
    if (deadline) {
      const date = deadline.toDate();
      date.setHours(23, 59, 59);
      savingGoalDeadline = date.toISOString();
    }

    const savingGoal: SavingGoal = {
      name: name,
      target: target ? target : 0,
      deadline: savingGoalDeadline,
      saving_goal_id: "",
      progress: 0,
      username: "",
    };

    if (recurringChecked) {
      savingGoal.is_recurring = true;
      savingGoal.recurring_amount = recurringAmount ? recurringAmount : 0;
    }

    try {
      validationSchema.validateSync(savingGoal);
      createSavingGoalMutation.mutate(savingGoal);
    } catch (e) {
      const err = e as ValidationError;
      onAlert({ open: true, type: "error", title: err.errors[0] });
    }
  }

  return (
    <Dialog open={open} onClose={onClose} fullWidth>
      <Box
        component="form"
        onSubmit={(e) => createSavingGoal(e)}
        sx={{
          maxWidth: "500px",
        }}
      >
        <Grid container spacing={2}>
          <Grid xs={12}>
            <Typography variant={"h4"}>New Saving Goal</Typography>
            <Divider />
          </Grid>

          {/*Name*/}
          <Grid xs={12}>
            <TextField
              margin={"none"}
              name={"name"}
              value={name}
              fullWidth={true}
              type={"text"}
              label={"Name"}
              variant={"outlined"}
              required
              onChange={(e) => setName(e.target.value)}
            />
          </Grid>

          {/*Target*/}
          <Grid xs={6}>
            <TextField
              margin={"normal"}
              sx={{ marginTop: "0px" }}
              name={"target"}
              value={target || ""}
              fullWidth={true}
              type={"number"}
              label={"Target"}
              variant={"outlined"}
              required
              onChange={(e) => setTarget(Number(e.target.value))}
            />
          </Grid>

          {/*Deadline*/}
          <Grid xs={6}>
            <DatePicker
              label="Deadline"
              value={deadline}
              disablePast
              onChange={(newDate) => setDeadline(newDate)}
              sx={{ width: "100%" }}
            />
          </Grid>

          {/*Recurring goal switch*/}
          <Grid xs={6}>
            <FormGroup>
              <FormControlLabel
                control={
                  <Switch
                    checked={recurringChecked}
                    onChange={(e) => setRecurringChecked(e.target.checked)}
                  />
                }
                label="Recurring"
              />
            </FormGroup>

            {recurringChecked && (
              <TextField
                margin={"normal"}
                sx={{ marginTop: "0px" }}
                name={"Recurring amount"}
                value={recurringAmount || ""}
                fullWidth={true}
                type={"number"}
                label={"Recurring amount"}
                variant={"outlined"}
                required
                onChange={(e) => setRecurringAmount(Number(e.target.value))}
              />
            )}
          </Grid>

          {/*Buttons*/}
          <Grid xs={12}>
            <div className={"flex justify-end"}>
              <Button
                variant={"contained"}
                color={"gray"}
                sx={{ fontSize: "16px" }}
                onClick={() => onClose()}
              >
                Cancel
              </Button>
              <Button
                type={"submit"}
                sx={{ fontSize: "16px", marginLeft: "0.5rem" }}
                variant={"contained"}
                loading={createSavingGoalMutation.isPending}
              >
                Save
              </Button>
            </div>
          </Grid>
        </Grid>
      </Box>
    </Dialog>
  );
}
