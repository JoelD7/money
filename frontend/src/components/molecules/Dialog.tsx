import { Dialog as MuiDialog, dialogClasses, DialogProps } from "@mui/material";

export function Dialog({ children, ...other }: DialogProps) {
  const styles = {
    [`& .${dialogClasses.paper}`]: {
      padding: "1.5rem",
      maxWidth: "fit-content",
    },
  };

  return (
    <MuiDialog sx={styles} {...other}>
      {children}
    </MuiDialog>
  );
}
