import { Alert, AlertTitle, Snackbar } from "@mui/material";
import {useState} from "react";

type ErrorSnackbarProps = {
  openProp: boolean;
  title: string;
  message?: string;
};

export function ErrorSnackbar({
  openProp,
  title,
  message,
}: ErrorSnackbarProps) {
  const[open, setOpen] = useState<boolean>(openProp)
  return (
    <div>
      <Snackbar
        open={open}
        autoHideDuration={6000}
        anchorOrigin={{ vertical: "top", horizontal: "right" }}
        onClose={()=> setOpen(false)}
      >
        <Alert
          onClose={()=> setOpen(false)}
          severity="error"
          variant="filled"
          sx={{ width: "100%" }}
        >
          {message ? (
            <>
              <AlertTitle>{title}</AlertTitle>
              {message}
            </>
          ) : (
            <>{title}</>
          )}
        </Alert>
      </Snackbar>
    </div>
  );
}
