import { Alert, AlertTitle, capitalize, Snackbar } from "@mui/material";
import { useState } from "react";
import {SnackAlert, User} from "../../types";
import { NewExpenseDialog } from "./NewExpenseDialog.tsx";

type NewExpenseProps = {
  open: boolean;
  user?: User;
  onClose: () => void;
};

export function NewExpense({ open, onClose, user }: NewExpenseProps) {
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

      <NewExpenseDialog
        key={key}
        open={open}
        user={user}
        onClose={handleClose}
        onAlert={handleAlert}
      />
    </>
  );
}
