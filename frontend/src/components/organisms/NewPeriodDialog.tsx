import { Dialog } from "../molecules";
import Grid from "@mui/material/Unstable_Grid2";
import {
  Box,
  Divider,
  FormControlLabel,
  FormGroup,
  IconButton,
  Switch,
  TextField,
  Tooltip,
  Typography,
} from "@mui/material";
import { FormEvent, useState } from "react";
import dayjs, { Dayjs } from "dayjs";
import { DatePicker } from "@mui/x-date-pickers";
import { Button } from "../atoms";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import api from "../../api";
import { AxiosError } from "axios";
import { APIError, Period, SnackAlert, User } from "../../types";
import * as yup from "yup";
import { ValidationError } from "yup";
import HelpIcon from "@mui/icons-material/Help";
import { useGetUser } from "../../queries";
import { USER } from "../../queries/keys";

type NewPeriodDialogProps = {
  open: boolean;
  onClose: () => void;
  onAlert: (alert?: SnackAlert) => void;
};

export function NewPeriodDialog({ open, onClose, onAlert }: NewPeriodDialogProps) {
  const [startDate, setStartDate] = useState<Dayjs | null>(dayjs());
  const [endDate, setEndDate] = useState<Dayjs | null>(dayjs().add(1, "month"));
  const [name, setName] = useState("");
  const [isCurrent, setIsCurrent] = useState(true);
  const getUser = useGetUser();
  const user: User | undefined = getUser.data;

  const queryClient = useQueryClient();

  const currentPeriodExplainer =
    "If you set this period as 'current', all expenses, savings and income will be " +
    "created with this as their period by default. The data on the home page will also be updated to reflect " +
    "the calculations based on this period.";

  const patchUserMu = useMutation({
    mutationFn: api.patchUser,
    onSuccess: () => {
      onAlert({
        open: true,
        type: "success",
        title: "Period created successfully",
      });
      onClose();

      queryClient
        .invalidateQueries({ queryKey: [USER] })
        .then(null, (error) => {
          console.error("Error invalidating users query", error);
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

  const newPeriodMu = useMutation({
    mutationFn: api.createPeriod,
    onSuccess: (res) => {
      if (isCurrent) {
        updateUserCurrentPeriod(res.data.period_id as string);
      } else {
        onAlert({
          open: true,
          type: "success",
          title: "Period created successfully",
        });
        onClose();
      }
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
    start_date: yup.string().required("Start date is required"),
    end_date: yup.string().required("End date is required"),
  });

  function createPeriod(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const period: Period = {
      username: "",
      name: name,
      start_date: startDate ? startDate.format("") : "",
      end_date: endDate ? endDate.format("") : "",
      period_id: "",
    };

    try {
      validationSchema.validateSync(period);
      newPeriodMu.mutate(period);
    } catch (e) {
      const err = e as ValidationError;
      onAlert({ open: true, type: "error", title: err.errors[0] });
    }
  }

  function updateUserCurrentPeriod(periodID: string) {
    const userToUpdate: User = {
      username: user ? user.username : "",
      remainder: 0,
      current_period: periodID,
    }

    try {
      patchUserMu.mutate(userToUpdate);
    } catch (e) {
      const err = e as ValidationError;
      onAlert({ open: true, type: "error", title: err.errors[0] });
    }
  }

  return (
    <Dialog open={open} onClose={onClose}>
      <Box component={"form"} onSubmit={createPeriod}>
        <Grid
          container
          spacing={2}
          bgcolor={"white.main"}
          borderRadius="1rem"
          width={"500px"}
          p="1.5rem"
        >
          <Grid xs={12}>
            <Typography variant={"h4"}>Create new period</Typography>
            <Divider />
          </Grid>

          {/*Set as current period*/}
          <Grid xs={12}>
            <div className={"flex"}>
              <FormGroup>
                <FormControlLabel
                  control={
                    <Switch
                      checked={isCurrent}
                      onChange={(e) => setIsCurrent(e.target.checked)}
                    />
                  }
                  label="Set as current"
                />
              </FormGroup>

              <Tooltip title={currentPeriodExplainer}>
                <IconButton>
                  <HelpIcon />
                </IconButton>
              </Tooltip>
            </div>
          </Grid>

          <Grid xs={12}>
            <TextField
              margin={"none"}
              name={"name"}
              value={name}
              fullWidth={true}
              type={"text"}
              label={"Name"}
              variant={"outlined"}
              onChange={(e) => setName(e.target.value)}
            />
          </Grid>

          <Grid xs={6}>
            <DatePicker
              label="Start date"
              sx={{ width: "100%" }}
              value={startDate}
              disablePast
              onChange={(newDate) => setStartDate(newDate)}
            />
          </Grid>

          <Grid xs={6}>
            <DatePicker
              label="End date"
              sx={{ width: "100%" }}
              value={endDate}
              disablePast
              minDate={startDate === null ? undefined : startDate}
              onChange={(newDate) => setEndDate(newDate)}
            />
          </Grid>

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
                loading={newPeriodMu.isPending}
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
