import { Autocomplete, AutocompleteInputChangeReason, CircularProgress, TextField } from "@mui/material";
import { useDebounce } from "@uidotdev/usehooks";
import React, { useState } from "react";
import { useGetPeriodsInfinite } from "../../queries";
import { Period } from "../../types";
import { ErrorSnackbar } from "./ErrorSnackbar.tsx";

type PeriodSelectorProps = {
  period?: Period;
  onPeriodChange: (period: Period) => void;
  active?: boolean;
};

export function PeriodSelector({ period, active, onPeriodChange }: PeriodSelectorProps) {
  const [searchPeriodName, setSearchPeriodName] = useState<string>("");
  const debouncedSearch = useDebounce(searchPeriodName, 500);

  const getPeriodsQuery = useGetPeriodsInfinite({ active, periodName: debouncedSearch });

  let periods: Period[] = (() => {
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

  function getOptionLabel(option: Period) {
    if (!option || periods.length === 0) {
      return "";
    }

    return option.name;
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
        value={period}
        filterOptions={(x) => x}
        isOptionEqualToValue={(option, value) => option.name === value.name}
        getOptionLabel={(option) => getOptionLabel(option)}
        onInputChange={(_, newValue: string, reason: AutocompleteInputChangeReason) => {
          if (reason !== "reset") {
            setSearchPeriodName(newValue);
          }
        }}
        onChange={(_, newValue: Period | null) => {
          if (newValue) {
            onPeriodChange(newValue);
          }
        }}
        options={periods}
        loading={getPeriodsQuery.isLoading}
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
                  {getPeriodsQuery.isFetching ? (
                    <CircularProgress color="inherit" size={20} />
                  ) : null}
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
