import { SyntheticEvent, useEffect, useState } from "react";
import { Alert, AlertTitle, Snackbar } from "@mui/material";
import { useQueryClient } from "@tanstack/react-query";

export function BackgroundRefetchErrorSnackbar() {
  const [open, setOpen] = useState(false);
  const queryClient = useQueryClient();

  useEffect(() => {
    const unsubscribe = queryClient.getQueryCache().subscribe((event) => {
      if (
        event.query.state.fetchStatus === "idle" &&
        event.query.state.status === "error"
      ) {
        setOpen(true);
      }
    });

    return () => unsubscribe();
  }, [queryClient]);

  const handleClose = (_event: Event | SyntheticEvent, reason: string) => {
    if (reason === "clickaway") {
      return;
    }
    setOpen(false);
  };

  return (
    <Snackbar
      open={open}
      autoHideDuration={6000}
      onClose={handleClose}
      anchorOrigin={{ vertical: "top", horizontal: "right" }}
    >
      <Alert severity="warning">
        <AlertTitle>Background refetch failed</AlertTitle>
        Data might be outdated
      </Alert>
    </Snackbar>
  );
}
