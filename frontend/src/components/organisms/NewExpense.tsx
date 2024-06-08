import Grid from "@mui/material/Unstable_Grid2";
import {
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
import CalendarTodayIcon from "@mui/icons-material/CalendarToday";
import { useState } from "react";
import { CategorySelect } from "./CategorySelect.tsx";
import { useQuery } from "@tanstack/react-query";
import api from "../../api";
import { User } from "../../types";
import { DatePicker } from "@mui/x-date-pickers";
import dayjs, { Dayjs } from "dayjs";
import { Button } from "../atoms";

type ExpenseTypeOption = {
  label: string;
  value: string;
};

export function NewExpense() {
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

  const [amount, setAmount] = useState<number>();
  const [title, setTitle] = useState<string>("");
  const [description, setDescription] = useState<string>("");
  const [category, setCategory] = useState<string>("");
  const [type, setType] = useState<string>("");
  const [recurringDay, setRecurringDay] = useState<number>(1);
  const [date, setDate] = useState<Dayjs | null>(dayjs());

  const getUser = useQuery({
    queryKey: ["user"],
    queryFn: () => api.getUser(),
  });

  const user: User | undefined = getUser.data?.data;

  return (
    <Grid
      container
      spacing={2}
      bgcolor={"white.main"}
      borderRadius="1rem"
      width={"700px"}
      style={{ border: "1px solid black" }}
      p="1.5rem"
    >
      <Grid xs={12}>
        <Typography variant={"h4"}>New Expense</Typography>
        <Divider />
      </Grid>

      {/*Title*/}
      <Grid xs={12}>
        <TextField
          margin={"none"}
          name={"title"}
          value={title}
          fullWidth={true}
          type={"text"}
          label={"Title"}
          variant={"outlined"}
          required
          onChange={(e) => setTitle(e.target.value)}
        />
      </Grid>

      {/*Left side*/}
      <Grid xs={6}>
        {/*Amount*/}
        <TextField
          margin={"normal"}
          name={"amount"}
          value={amount}
          fullWidth={true}
          type={"number"}
          label={"Amount"}
          variant={"outlined"}
          required
          onChange={(e) => setAmount(Number(e.target.value))}
        />

        {/*Date*/}
        <DatePicker
          label="Date"
          sx={{ marginTop: "10px" }}
          value={date}
          onChange={(newDate) => setDate(newDate)}
        />

        {/*Description*/}
        <TextField
          margin={"normal"}
          name={"description"}
          value={description}
          multiline
          minRows={3}
          maxRows={6}
          fullWidth={true}
          type={"text"}
          label={"Description (optional)"}
          variant={"outlined"}
          size={"medium"}
          required
          onChange={(e) => setDescription(e.target.value)}
        />
      </Grid>

      {/*Right side*/}
      <Grid xs={6}>
        {/*Category*/}
        {user && user.categories && (
          <div className={"mb-2"}>
            <CategorySelect
              categories={user.categories}
              selected={category === "" ? [] : [category]}
              onSelectedUpdate={(selected) => setCategory(selected.join(""))}
              width={"400px"}
              label={"Category(optional)"}
            />
          </div>
        )}

        {/*Type*/}
        <>
          <FormControl>
            <FormLabel id="expense-type-radio-buttons-group-label">
              Type
            </FormLabel>
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
                  control={<Radio />}
                  label={expenseType.label}
                />
              ))}
            </RadioGroup>
          </FormControl>

          <FormControl fullWidth disabled={type !== "recurring"}>
            <InputLabel id="recurrent-expense-day-select-label">
              Every
            </InputLabel>
            <Select
              labelId="recurrent-expense-day-select-label"
              id="recurrent-expense-select"
              value={recurringDay}
              label="Day"
              startAdornment={
                <CalendarTodayIcon sx={{ marginRight: "10px" }} />
              }
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
          >
            Cancel
          </Button>
          <Button
            sx={{ fontSize: "16px", marginLeft: "0.5rem" }}
            variant={"contained"}
          >
            Save
          </Button>
        </div>
      </Grid>
    </Grid>
  );
}
