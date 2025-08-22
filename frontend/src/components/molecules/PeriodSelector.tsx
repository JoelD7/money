import { FormControl, InputLabel, MenuItem, Select } from "@mui/material";
import React from "react";
import { useGetPeriodsInfinite } from "../../queries";
import { v4 as uuidv4 } from "uuid";
import { Period } from "../../types";

type PeriodSelectorProps = {
  period: string;
  onPeriodChange: (value: string) => void;
  active?: boolean;
};

export function PeriodSelector({ period, active, onPeriodChange }: PeriodSelectorProps) {
  const labelId: string = uuidv4();

  const getPeriodsQuery = useGetPeriodsInfinite({ active });

  const periods: Period[] = (() => {
    if (getPeriodsQuery.data) {
      return getPeriodsQuery.data.pages.map((page) => page.periods).flat();
    }
    return [];
  })();

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

  return (
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
        onChange={(e) => onPeriodChange(e.target.value)}
      >
        {Array.isArray(periods) &&
          periods.map((p) => (
            <MenuItem key={p.period_id} id={p.period_id} value={p.period_id}>
              {p.name}
            </MenuItem>
          ))}
      </Select>
    </FormControl>
  );
}
