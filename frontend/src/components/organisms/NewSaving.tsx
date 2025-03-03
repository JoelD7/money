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
import React, { useState } from "react";
import { v4 as uuidv4 } from "uuid";
import { useGetPeriodsInfinite } from "../../queries";

type NewSavingProps = {
  open: boolean;
  onClose: () => void;
};

export function NewSaving({ open, onClose }: NewSavingProps) {
  const labelId: string = uuidv4();

  const [period, setPeriod] = useState<string>("");
  const [amount, setAmount] = useState<number | null>(null);

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

  function handleMenuScroll(e: React.UIEvent<HTMLDivElement, UIEvent>) {
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
        onSubmit={() => {}}
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
            <FormControl sx={{ width: "150px" }}>
              <InputLabel id={labelId}>Period</InputLabel>

              <Select
                labelId={labelId}
                id={"Period"}
                MenuProps={{
                  slotProps: {
                    paper: {
                      onScroll: handleMenuScroll,
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
          <Grid xs={6}>
            <FormControl sx={{ width: "150px" }}>
              <InputLabel id={labelId}>Saving goal</InputLabel>

              <Select
                labelId={labelId}
                id={"Goal"}
                MenuProps={{
                  slotProps: {
                    paper: {
                      onScroll: handleMenuScroll,
                    },
                  },
                  PaperProps: {
                    sx: {
                      maxHeight: 150,
                    },
                  },
                }}
                label={"Goal"}
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
        </Grid>
      </Box>
    </Dialog>
  );
}
