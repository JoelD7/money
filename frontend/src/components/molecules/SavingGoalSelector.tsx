import {
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  SelectChangeEvent,
} from "@mui/material";
import React from "react";
import { v4 as uuidv4 } from "uuid";
import { useGetSavingGoalsInfinite } from "../../queries";
import { SavingGoal } from "../../types";

type SavingGoalSelectorProps = {
  onSavingGoalChange: (savingGoalId: string) => void;
  savingGoalID: string;
};

export function SavingGoalSelector({
  onSavingGoalChange,
  savingGoalID,
}: SavingGoalSelectorProps) {
  const labelId: string = uuidv4();

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

  function handleSavingGoalChange(e: SelectChangeEvent) {
    onSavingGoalChange(e.target.value);
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

  return (
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
        value={savingGoals.length > 0 ? savingGoalID : ""}
        onChange={handleSavingGoalChange}
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
  );
}
