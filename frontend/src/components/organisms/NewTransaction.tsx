import { Alert, AlertTitle, capitalize, Snackbar } from "@mui/material";
import { useState } from "react";
import {SnackAlert, User} from "../../types";
import { NewExpense } from "./NewExpense.tsx";
import { NewIncome } from "./NewIncome.tsx";

type NewTransactionProps = {
  open: boolean;
  onClose: () => void;
  type: "income" | "expense";
  user?: User;
};

export function NewTransaction({ open, onClose, type , user}: NewTransactionProps) {
  const [alert, setAlert] = useState<SnackAlert>({
    open: false,
    type: "success",
    title: "",
  });
  const [key, setKey] = useState<number>(0);

  function handleClose() {
    setKey(key + 1);
    onClose();
  }

  function handleAlert(alert?: SnackAlert) {
    if (alert) {
      setAlert(alert);
    }
  }

  return (
    <>
      <Snackbar
        open={alert.open}
        onClose={() => setAlert({ ...alert, open: false })}
        autoHideDuration={6000}
        anchorOrigin={{ vertical: "top", horizontal: "right" }}
      >
        <Alert variant={"filled"} severity={alert.type}>
          <AlertTitle>{capitalize(alert.type)}</AlertTitle>
          {alert.title}
        </Alert>
      </Snackbar>

      {type === "income" ? (
        <NewIncome
          key={key}
          open={open}
          onClose={handleClose}
          onAlert={handleAlert}
        />
      ) : (
        <NewExpense
          key={key}
          user={user}
          open={open}
          onClose={handleClose}
          onAlert={handleAlert}
        />
      )}
    </>
  );
}
