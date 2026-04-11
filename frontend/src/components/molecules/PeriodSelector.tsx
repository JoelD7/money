import { Autocomplete, AutocompleteInputChangeReason, CircularProgress, TextField } from "@mui/material";
import React, { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { useGetPeriodsInfinite } from "../../queries";
import { setAllPeriods } from "../../store/user.ts";
import { Period } from "../../types";
import { ErrorSnackbar } from "./ErrorSnackbar.tsx";

type PeriodSelectorProps = {
  period?: Period;
  onPeriodChange: (period?: Period) => void;
  active?: boolean;
};

export function PeriodSelector({ period, active, onPeriodChange }: PeriodSelectorProps) {
  const [searchPeriodName, setSearchPeriodName] = useState<string>("");

  const dispatch = useDispatch()

  const getPeriodsQuery = useGetPeriodsInfinite({ pageSize: 100, active });
  const allPeriods = useSelector((state: any) => state.usersReducer.allPeriods);
  const periods: Period[] = (() => {
    if (getPeriodsQuery.data) {
      return getPeriodsQuery.data.pages.map((page) => page.periods).flat();
    }
    return [];
  })();

  useEffect(() => {
    if (shouldUpdateAllPeriods()) {
      dispatch(setAllPeriods(periods))
    }
  }, [periods, dispatch, getPeriodsQuery.data]);

  function shouldUpdateAllPeriods(): boolean {
    // When this flag is true, then we fetch ONLY the active periods from the server. We don't want to update the
    // allPeriods state in that case.
    if (active) {
      return false
    }

    if (!allPeriods) {
      return true
    }

    return allPeriods.length === 0
      // This might mean that the user created a new period and we need to update the allPeriods state.
      || allPeriods.length !== periods.length
  }


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
        getOptionKey={(period) => period.period_id}
        filterOptions={(period) => {
          return period.filter((p) => p.name.toLowerCase().startsWith(searchPeriodName.toLowerCase()));
        }}
        isOptionEqualToValue={(option, value) => {
          if (!value || !option) {
            return false;
          }

          return option.name === value.name
        }}
        getOptionLabel={(option) => getOptionLabel(option)}
        onInputChange={(_, newValue: string, reason: AutocompleteInputChangeReason) => {
          if (reason !== "reset") {
            setSearchPeriodName(newValue);
          }
        }}
        onChange={(_, newValue: Period | null) => {
          if (newValue) {
            onPeriodChange(newValue);
          } else {
            onPeriodChange(undefined);
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
