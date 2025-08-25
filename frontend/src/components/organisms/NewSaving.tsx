import { Dialog, PeriodSelector, SavingGoalSelector } from "../molecules";
import { Box, Divider, TextField, Typography } from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import { FormEvent, useState } from "react";
import { Saving, SnackAlert } from "../../types";
import * as yup from "yup";
import { ValidationError } from "yup";
import { Button } from "../atoms";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import api from "../../api";
import { savingsKeys } from "../../api/savings.ts";
import { savingGoalKeys } from "../../queries/saving_goals.ts";

type NewSavingProps = {
  open: boolean;
  onClose: () => void;
  onAlert: (alert?: SnackAlert) => void;
  savingGoalId?: string;
};

export function NewSaving({ open, onClose, onAlert, savingGoalId }: NewSavingProps) {
  const showSavingGoalSelector: boolean = savingGoalId ? savingGoalId === "" : true;
  const validationSchema = yup.object({
    amount: yup.number().required("Amount is required").moreThan(0, "Amount is required"),
    period_id: yup.string().required("Period is required"),
    saving_goal_id: yup.string().optional(),
  });

  const [period, setPeriod] = useState<string>("");
  const [amount, setAmount] = useState<number | null>(null);
  const [savingGoal, setSavingGoal] = useState<string>(savingGoalId ? savingGoalId : "");

  const queryClient = useQueryClient();
  const createSavingMutation = useMutation({
    mutationFn: api.createSaving,
    onSuccess: () => {
      onAlert({
        open: true,
        type: "success",
        title: "Saving created successfully",
      });

      onClose();
      queryClient
        .invalidateQueries({ queryKey: [...savingsKeys.all] })
        .then(() => {})
        .catch((e) => {
          console.error("Error invalidating savings query", e);
        });

      queryClient
        .invalidateQueries({ queryKey: [...savingGoalKeys.all] })
        .then(() => {})
        .catch((e) => {
          console.error("Error invalidating saving goals query", e);
        });
    },
    onError: (error) => {
      if (error) {
        const err = error as Error;
        onAlert({
          open: true,
          type: "error",
          title: err.message,
        });
      }
    },
  });

  function createSaving(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();

    const saving: Saving = {
      saving_id: "",
      username: "",
      amount: amount as number,
      period_id: period,
      saving_goal_id: savingGoal,
    };

    try {
      validationSchema.validateSync(saving);
      createSavingMutation.mutate(saving);
    } catch (e) {
      const err = e as ValidationError;
      onAlert({
        open: true,
        type: "error",
        title: err.errors[0],
      });
      console.error("Error validating saving", err.errors[0]);
    }
  }

  function onSavingGoalChange(savingGoalId: string) {
    setSavingGoal(savingGoalId);
  }

  return (
    <Dialog open={open} onClose={onClose}>
      <Box
        component="form"
        onSubmit={createSaving}
        sx={{
          maxWidth: "500px",
        }}
      >
        <Grid container spacing={2}>
          {/*Title*/}
          <Grid xs={12}>
            <Typography variant={"h4"}>New Saving</Typography>
            <Divider />
          </Grid>

          {/*Amount*/}
          <Grid xs={6}>
            <TextField
              margin={"normal"}
              sx={{ marginTop: "0px" }}
              name={"amount"}
              value={amount || ""}
              fullWidth={true}
              type={"number"}
              label={"Amount"}
              variant={"outlined"}
              required
              onChange={(e) => setAmount(Number(e.target.value))}
            />
          </Grid>

          {/*Period*/}
          <Grid xs={6}>
            <PeriodSelector period={period} onPeriodChange={setPeriod} />
          </Grid>

          {/* Saving goal */}
          {showSavingGoalSelector && (
            <Grid xs={12}>
              <SavingGoalSelector
                savingGoalID={savingGoal}
                onSavingGoalChange={onSavingGoalChange}
              />
            </Grid>
          )}

          {/*Buttons*/}
          <Grid xs={12} alignSelf={"end"}>
            <div className={"flex justify-end"}>
              <Button
                variant={"contained"}
                color={"gray"}
                sx={{ fontSize: "16px" }}
                onClick={onClose}
              >
                Cancel
              </Button>
              <Button
                type={"submit"}
                sx={{ fontSize: "16px", marginLeft: "0.5rem" }}
                variant={"contained"}
                loading={createSavingMutation.isPending}
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
