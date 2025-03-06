import { Dialog, ErrorSnackbar } from "../molecules";
import {
  Box,
  Divider,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  TextField,
  Typography,
} from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import React, { FormEvent, useState } from "react";
import { v4 as uuidv4 } from "uuid";
import {
  queryKeys,
  useGetPeriodsInfinite,
  useGetSavingGoalsInfinite,
} from "../../queries";
import { Saving, SavingGoal, SnackAlert } from "../../types";
import * as yup from "yup";
import { ValidationError } from "yup";
import { Button } from "../atoms";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import api from "../../api";

type NewSavingProps = {
  open: boolean;
  onClose: () => void;
  onAlert: (alert?: SnackAlert) => void;
};

export function NewSaving({ open, onClose, onAlert }: NewSavingProps) {
  const labelId: string = uuidv4();
  const validationSchema = yup.object({
    amount: yup.number().required("Amount is required").moreThan(0, "Amount is required"),
    period: yup.string().required("Period is required"),
    saving_goal_id: yup.string().optional(),
  });

  const [period, setPeriod] = useState<string>("");
  const [amount, setAmount] = useState<number | null>(null);
  const [savingGoal, setSavingGoal] = useState<string>("");

  const getPeriodsQuery = useGetPeriodsInfinite();
  const periods: string[] = (() => {
    if (getPeriodsQuery.data) {
      return getPeriodsQuery.data.pages
        .map((page) => page.periods)
        .flat()
        .map((p) => p.name);
    }

    return [];
  })();

  const getSavingGoalsQuery = useGetSavingGoalsInfinite();
  const savingGoals: SavingGoal[] = (() => {
    if (getSavingGoalsQuery.data) {
      return getSavingGoalsQuery.data.pages
        .map((page) => page.saving_goals)
        .flat()
        .map((p) => p);
    }

    return [];
  })();

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
        .invalidateQueries({ queryKey: [queryKeys.SAVINGS] })
        .then(() => {})
        .catch((e) => {
          console.error("Error invalidating savings query", e);
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

  function handlePeriodsMenuScroll(e: React.UIEvent<HTMLDivElement, UIEvent>) {
    const { scrollTop, clientHeight, scrollHeight } = e.currentTarget;
    if (
      scrollTop + clientHeight >= scrollHeight - 5 &&
      !(getPeriodsQuery.isFetching || getPeriodsQuery.isFetchingNextPage)
    ) {
      getPeriodsQuery
        .fetchNextPage()
        .then(() => {})
        .catch((e) => {
          console.error("Error fetching more periods", e);
        });
    }
  }

  function handleSavingGoalsMenuScroll(e: React.UIEvent<HTMLDivElement, UIEvent>) {
    const { scrollTop, clientHeight, scrollHeight } = e.currentTarget;
    if (
      scrollTop + clientHeight >= scrollHeight - 5 &&
      !(getSavingGoalsQuery.isFetching || getSavingGoalsQuery.isFetchingNextPage)
    ) {
      getSavingGoalsQuery
        .fetchNextPage()
        .then(() => {})
        .catch((e) => {
          console.error("Error fetching more saving goals", e);
        });
    }
  }

  function createSaving(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();

    const saving: Saving = {
      saving_id: "",
      username: "",
      amount: amount as number,
      period: period,
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

  return (
    <Dialog open={open} onClose={onClose}>
      {getPeriodsQuery.isError && (
        <ErrorSnackbar
          openProp={getPeriodsQuery.isError}
          title={"Error fetching periods"}
          message={getPeriodsQuery.error.message}
        />
      )}

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
            <FormControl sx={{ width: "100%" }}>
              <InputLabel id={labelId}>Period</InputLabel>

              <Select
                labelId={labelId}
                id={"Period"}
                MenuProps={{
                  slotProps: {
                    paper: {
                      onScroll: handlePeriodsMenuScroll,
                    },
                  },
                  PaperProps: {
                    sx: {
                      maxHeight: 150,
                    },
                  },
                }}
                label={"Period"}
                value={periods.length > 0 ? period : ""}
                onChange={(e) => setPeriod(e.target.value)}
              >
                {Array.isArray(periods) &&
                  periods.map((p) => (
                    <MenuItem key={p} id={p} value={p}>
                      {p}
                    </MenuItem>
                  ))}
              </Select>
            </FormControl>
          </Grid>

          {/* Saving goal */}
          <Grid xs={12}>
            <FormControl sx={{ width: "100%" }}>
              <InputLabel id={labelId}>Saving goal</InputLabel>

              <Select
                labelId={labelId}
                id={"Goal"}
                MenuProps={{
                  slotProps: {
                    paper: {
                      onScroll: handleSavingGoalsMenuScroll,
                    },
                  },
                  PaperProps: {
                    sx: {
                      maxHeight: 150,
                    },
                  },
                }}
                label={"Goal"}
                value={savingGoals.length > 0 ? savingGoal : ""}
                onChange={(e) => setSavingGoal(e.target.value)}
              >
                {Array.isArray(savingGoals) &&
                  savingGoals.map((sg) => (
                    <MenuItem
                      key={sg.saving_goal_id}
                      id={sg.saving_goal_id}
                      value={sg.saving_goal_id}
                    >
                      {sg.name}
                    </MenuItem>
                  ))}
              </Select>
            </FormControl>
          </Grid>

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
