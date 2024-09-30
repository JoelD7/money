import { APIError, Income, Period, SnackAlert } from "../../types";
import { Box, Dialog, Divider, TextField, Typography } from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";
import { FormEvent, useState } from "react";
import { DatePicker } from "@mui/x-date-pickers";
import dayjs, { Dayjs } from "dayjs";
import { Button } from "../atoms";
import {useMutation, useQueryClient} from "@tanstack/react-query";
import api from "../../api";
import { AxiosError } from "axios";
import * as yup from "yup";
import { ValidationError } from "yup";
import {queryKeys, useGetPeriod} from "../../queries";

type NewIncomeProps = {
  onClose: () => void;
  onAlert: (alert?: SnackAlert) => void;
  open: boolean;
};

export function NewIncome({ onClose, open, onAlert }: NewIncomeProps) {
  const dialogStyle = {
    "& .MuiDialog-paper": {
      width: "705px",
      height: "390px",
      maxWidth: "100%",
    },
  };

  const getPeriod = useGetPeriod();
  const period: Period | undefined = getPeriod.data?.data;

  const queryClient = useQueryClient()

  const ciMutation = useMutation({
    mutationFn: api.createIncome,
    onSuccess: () => {
      onAlert({
        open: true,
        type: "success",
        message: "Income created successfully",
      });
      onClose();

      queryClient.invalidateQueries({queryKey: [queryKeys.PERIOD_STATS]});
    },
    onError: (error) => {
      if (error) {
        const err = error as AxiosError;
        const responseError = err.response?.data as APIError;
        onAlert({
          open: true,
          type: "error",
          message: responseError.message as string,
        });
      }
    },
  });

  const [amount, setAmount] = useState<number | null>();
  const [name, setName] = useState<string>("");
  const [date, setDate] = useState<Dayjs | null>(dayjs());
  const [notes, setNotes] = useState<string>("");

  const validationSchema = yup.object({
    name: yup.string().required("Name is required"),
    amount: yup
      .number()
      .required("Amount is required")
      .moreThan(0, "Amount is required"),
    created_date: yup.date().required("Date is required"),
  });

  function createIncome(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const income: Income = {
      amount: amount as number,
      name: name,
      created_date: date ? date.format("") : "",
      notes: notes,
      period: period ? period.period : "",
    };

    try {
      validationSchema.validateSync(income);
      ciMutation.mutate(income);
    } catch (e) {
      const err = e as ValidationError;
      onAlert({ open: true, type: "error", message: err.errors[0] });
    }
  }

  return (
    <Dialog open={open} onClose={onClose} fullWidth sx={dialogStyle}>
      <Box component={"form"} onSubmit={createIncome} height={"100%"}>
        <Grid
          container
          spacing={1}
          bgcolor={"white.main"}
          borderRadius="1rem"
          width={"700px"}
          height={"100%"}
          p="1.5rem"
        >
          <Grid xs={12}>
            <Typography variant={"h4"}>New Income</Typography>
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

          {/*Date*/}
          <Grid xs={6}>
            <DatePicker
              label="Date"
              sx={{ width: "100%" }}
              value={date}
              onChange={(newDate) => setDate(newDate)}
            />
          </Grid>

          {/*Name*/}
          <Grid xs={6}>
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

          {/*Notes*/}
          <Grid xs={6}>
            <TextField
              name={"notes"}
              value={notes}
              multiline
              minRows={3}
              maxRows={6}
              fullWidth={true}
              type={"text"}
              label={"Notes (optional)"}
              variant={"outlined"}
              size={"medium"}
              onChange={(e) => setNotes(e.target.value)}
            />
          </Grid>

          <Grid xs={12} alignSelf={"end"}>
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
                loading={ciMutation.isPending}
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
