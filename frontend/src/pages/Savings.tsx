import {
  Button,
  Container,
  Navbar,
  NewSaving,
  PageTitle,
  SavingGoalsTable,
  SavingsTable,
  Snackbar,
} from "../components";
import Grid from "@mui/material/Unstable_Grid2";
import { Box, Typography } from "@mui/material";
import { SnackAlert } from "../types";
import { useState } from "react";

export function Savings() {
  const [open, setOpen] = useState(false);
  const [key, setKey] = useState(0);
  const [alert, setAlert] = useState<SnackAlert>({
    open: false,
    type: "success",
    title: "",
  });

  function openNewSavingDialog() {
    setOpen(true);
  }

  return (
    <Container>
      <Navbar />

      <Snackbar
        open={alert.open}
        title={alert.title}
        message={alert.message}
        severity={alert.type}
        onClose={() => setAlert({ ...alert, open: false })}
      />

      <Grid container position={"relative"} spacing={1} marginTop={"20px"}>
        {/*Title and summary*/}
        <Grid xs={12}>
          <PageTitle>Savings</PageTitle>
        </Grid>

        {/*Saving goals table*/}
        <Grid xs={12}>
          <SavingGoalsTable />
        </Grid>

        {/*Latest savings*/}
        <Grid xs={12} mt={"2rem"}>
          {/*Title and buttons*/}
          <Typography variant={"h5"}>Latest savings</Typography>
          <div className={"pt-4"}>
            <Button variant={"contained"} onClick={() => openNewSavingDialog()}>
              New saving
            </Button>
            <Box height={"fit-content"} paddingTop={"10px"}>
              <SavingsTable />
            </Box>
          </div>
        </Grid>
      </Grid>

      <NewSaving
        key={key}
        open={open}
        onClose={() => {
          setOpen(false);
          setKey(key + 1);
        }}
        onAlert={(a) => {
          if (a) {
            setAlert({ ...a });
          }
        }}
      />
    </Container>
  );
}
