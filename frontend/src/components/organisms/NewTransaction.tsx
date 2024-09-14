import { Alert, AlertTitle, capitalize, Snackbar } from "@mui/material";
import { useState } from "react";
import { SnackAlert } from "../../types";
import { NewExpense } from "./NewExpense.tsx";
import { NewIncome } from "./NewIncome.tsx";

type NewTransactionProps = {
  open: boolean;
  onClose: () => void;
  type: "income" | "expense";
};

export function NewTransaction({ open, onClose, type }: NewTransactionProps) {
  const [alert, setAlert] = useState<SnackAlert>({
    open: false,
    type: "success",
    message: "",
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
          {alert.message}
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
          open={open}
          onClose={handleClose}
          onAlert={handleAlert}
        />
      )}
    </>
  );
}
