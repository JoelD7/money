import { Alert, AlertTitle, Snackbar } from "@mui/material";

type ErrorSnackbarProps = {
  open: boolean;
  message: string;
  extraDetails?: string;
  onClose: () => void;
};

export function ErrorSnackbar({
  open,
  message,
  extraDetails,
  onClose,
}: ErrorSnackbarProps) {
  return (
    <div>
      <Snackbar
        open={open}
        autoHideDuration={6000}
        anchorOrigin={{ vertical: "top", horizontal: "right" }}
        onClose={onClose}
      >
        <Alert
          onClose={onClose}
          severity="error"
          variant="filled"
          sx={{ width: "100%" }}
        >
          {extraDetails ? (
            <>
              <AlertTitle>{message}</AlertTitle>
              {extraDetails}
            </>
          ) : (
            <>{message}</>
          )}
        </Alert>
      </Snackbar>
    </div>
  );
}
