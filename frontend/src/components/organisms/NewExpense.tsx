import {
  Box,
  Divider,
  FormControl,
  FormControlLabel,
  FormLabel,
  InputLabel,
  MenuItem,
  Radio,
  RadioGroup,
  Select,
  TextField,
  Typography,
} from "@mui/material";
import { FormEvent, useState } from "react";
import { APIError, Category, Expense, ExpenseType, SnackAlert, User } from "../../types";
import * as yup from "yup";
import { ValidationError } from "yup";
import Grid from "@mui/material/Unstable_Grid2";
import { DatePicker } from "@mui/x-date-pickers";
import { CategorySelect } from "./CategorySelect.tsx";
import CalendarTodayIcon from "@mui/icons-material/CalendarToday";
import { Button } from "../atoms";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import api from "../../api";
import { AxiosError } from "axios";
import dayjs, { Dayjs } from "dayjs";
import {Dialog, PeriodSelector} from "../molecules";
import { expensesQueryKeys } from "../../queries";
import {PERIOD_STATS} from "../../queries/keys";

type ExpenseTypeOption = {
  label: string;
  value: ExpenseType;
};

type NewExpenseProps = {
  open: boolean;
  user?: User;
  onClose: () => void;
  onAlert: (alert?: SnackAlert) => void;
};

export function NewExpense({ onClose, open, onAlert, user }: NewExpenseProps) {
  const expenseTypes: ExpenseTypeOption[] = [
    {
      label: "Regular",
      value: "regular",
    },
    {
      label: "Recurring",
      value: "recurring",
    },
  ];

  const [amount, setAmount] = useState<number | null>();
  const [name, setName] = useState<string>("");
  const [notes, setNotes] = useState<string>("");
  const [category, setCategory] = useState<Category>();
  const [type, setType] = useState<string>("regular");
  const [recurringDay, setRecurringDay] = useState<number>(1);
  const [date, setDate] = useState<Dayjs | null>(dayjs());
  const[period, setPeriod] = useState<string>((user && user.current_period) ? user.current_period : "");

  const queryClient = useQueryClient();

  const ceMutation = useMutation({
    mutationFn: api.createExpense,
    onSuccess: () => {
      onAlert({
        open: true,
        type: "success",
        title: "Expense created successfully",
      });

      queryClient
        .invalidateQueries({ queryKey: [...expensesQueryKeys.all] })
        .then(null, (error) => {
          console.error("Error invalidating expenses query", error);
        });

      queryClient
          .invalidateQueries({ queryKey: [PERIOD_STATS] })
          .then(null, (error) => {
            console.error("Error invalidating period stats query", error);
          });

      onClose();
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
    amount: yup.number().required("Amount is required").moreThan(0, "Amount is required"),
    created_date: yup.date().required("Date is required"),
    category_id: yup.string().required("Category is required"),
    type: yup.string().oneOf(["regular", "recurring"]),
    recurringDay: yup.number().when("type", {
      is: "recurring",
      then: (schema) => schema.required(),
      otherwise: (schema) => schema.optional(),
    }),
  });

  function createExpense(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const expense: Expense = {
      username: user ? user.username : "",
      expense_id: "",
      name: name,
      amount: amount ? amount : 0,
      notes: notes,
      category_id: category ? category.id : "",
      type: type as ExpenseType,
      created_date: date ? date.format("") : "",
      period_id: (user && user.current_period) ? user.current_period : "",
    };

    try {
      validationSchema.validateSync(expense);
      ceMutation.mutate(expense);
    } catch (e) {
      const err = e as ValidationError;
      onAlert({ open: true, type: "error", title: err.errors[0] });
    }
  }

  function onCategoryChange(selected: Category[]) {
    setCategory(selected[0]);
  }

  return (
    <Dialog open={open} onClose={onClose} fullWidth>
      <Box component="form" onSubmit={createExpense}>
        <Grid
          container
          spacing={2}
          bgcolor={"white.main"}
          borderRadius="1rem"
          width={"700px"}
        >
          <Grid xs={12}>
            <Typography variant={"h4"}>New Expense</Typography>
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
              // required
              onChange={(e) => setName(e.target.value)}
            />
          </Grid>

          {/*Left side*/}
          <Grid xs={6}>
            {/*Amount*/}
            <TextField
              margin={"normal"}
              sx={{ marginTop: "0px" }}
              name={"amount"}
              value={amount || ""}
              fullWidth={true}
              type={"number"}
              label={"Amount"}
              variant={"outlined"}
              // required
              onChange={(e) => setAmount(Number(e.target.value))}
            />

            {/*Date*/}
            <DatePicker
              label="Date"
              sx={{ marginTop: "10px" }}
              value={date}
              onChange={(newDate) => setDate(newDate)}
            />

            {/*Notes*/}
            <TextField
              margin={"normal"}
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

          {/*Right side*/}
          <Grid xs={6}>
            {/*Category*/}
            {user && user.categories && (
              <div className={"mb-2"}>
                <CategorySelect
                  categories={user.categories}
                  selected={category ? [category] : []}
                  onSelectedUpdate={onCategoryChange}
                  width={"400px"}
                  label={"Category(optional)"}
                />
              </div>
            )}

            <PeriodSelector period={period} onPeriodChange={setPeriod} active />

            {/*Type*/}
            <>
              <FormControl>
                <FormLabel id="expense-type-radio-buttons-group-label">Type</FormLabel>
                <RadioGroup
                  row
                  aria-labelledby="expense-type-radio-buttons-group-label"
                  name="row-radio-buttons-group"
                  onChange={(e) => setType(e.target.value)}
                >
                  {expenseTypes.map((expenseType) => (
                    <FormControlLabel
                      key={expenseType.value}
                      value={expenseType.value}
                      checked={type === expenseType.value}
                      control={<Radio />}
                      label={expenseType.label}
                    />
                  ))}
                </RadioGroup>
              </FormControl>

              <FormControl fullWidth disabled={type !== "recurring"}>
                <InputLabel id="recurrent-expense-day-select-label">Every</InputLabel>
                <Select
                  labelId="recurrent-expense-day-select-label"
                  id="recurrent-expense-select"
                  value={recurringDay}
                  label="Day"
                  startAdornment={<CalendarTodayIcon sx={{ marginRight: "10px" }} />}
                  onChange={(e) => setRecurringDay(Number(e.target.value))}
                  MenuProps={{ PaperProps: { sx: { maxHeight: 250 } } }}
                >
                  {Array.from({ length: 30 }, (_, i) => i + 1).map((day) => (
                    <MenuItem key={day} value={day}>
                      {day}th
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </>
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
                loading={ceMutation.isPending}
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
