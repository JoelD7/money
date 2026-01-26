import { Autocomplete, CircularProgress, FormControl, InputLabel, MenuItem, Select, TextField } from "@mui/material";
import React from "react";
import { useGetPeriodsInfinite } from "../../queries";
import { v4 as uuidv4 } from "uuid";
import { Period } from "../../types";
import { ErrorSnackbar } from "./ErrorSnackbar.tsx";

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

  function handlePeriodsMenuScroll(e: React.UIEvent<HTMLUListElement, UIEvent>) {
    const { scrollTop, clientHeight, scrollHeight } = e.currentTarget;
    if (
      scrollTop + clientHeight >= scrollHeight - 5 &&
      !(getPeriodsQuery.isFetching || getPeriodsQuery.isFetchingNextPage)
    ) {
      getPeriodsQuery
        .fetchNextPage()
        .then(() => { })
        .catch((e) => {
          console.error("Error fetching more periods", e);
        });
    }
  }

  function showErrorSnackbar() {
    if (getPeriodsQuery.isError && getPeriodsQuery.error.response) {
      return getPeriodsQuery.error.response.status !== 404;
    }

    return getPeriodsQuery.isError;
  }

  return (
    <>
      {showErrorSnackbar() && (
        <ErrorSnackbar
          openProp={getPeriodsQuery.isError}
          title={"Error fetching periods"}
          message={getPeriodsQuery.error ? getPeriodsQuery.error.message : ""}
        />
      )}

      <Autocomplete
        sx={{ width: "100%" }}
        isOptionEqualToValue={(option, value) => option.name === value.name}
        getOptionLabel={(option) => option.name}
        onChange={(_, newValue) => {
          if (newValue) {
            console.log("newValue", newValue);
          }
        }}
        options={periods}
        loading={getPeriodsQuery.isFetching}
        ListboxProps={{
          sx: {
            maxHeight: 150,
          },
          onScroll: handlePeriodsMenuScroll,
        }}
        renderInput={(params) => (
          <TextField
            {...params}
            label="Period"

            InputProps={{
              ...params.InputProps,
              endAdornment: (
                <>
                  {getPeriodsQuery.isFetching ? <CircularProgress color="inherit" size={20} /> : null}
                  {params.InputProps.endAdornment}
                </>
              ),
            }}
          />
        )}
      />

    </>
  );
}
