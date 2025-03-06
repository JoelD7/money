import { Alert, AlertTitle, Snackbar as MuiSnackbar } from "@mui/material";

type SnackbarProps = {
  open: boolean;
  title: string;
  message?: string;
  severity?: "error" | "success" | "info" | "warning";
  onClose: () => void;
};

export function Snackbar({
  open,
  title,
  message,
  onClose,
  severity = "info",
}: SnackbarProps) {
  return (
    <MuiSnackbar
      open={open}
      autoHideDuration={6000}
      anchorOrigin={{ vertical: "top", horizontal: "right" }}
      onClose={onClose}
    >
      <Alert
        onClose={onClose}
        severity={severity}
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
    </MuiSnackbar>
  );
}
