import { Button as MuiButton, ButtonProps, SxProps, Theme } from "@mui/material";
import { forwardRef, ReactNode } from "react";
import LoadingButton from "@mui/lab/LoadingButton";
import SaveIcon from "@mui/icons-material/Save";

type CustomButtonProps = {
  children: ReactNode;
  loading?: boolean;
} & ButtonProps;

export const Button = forwardRef<HTMLButtonElement, CustomButtonProps>((props, ref) => {
  const { sx, loading, variant, size, children, ...other } = props;
  let styles: SxProps<Theme> | undefined = {
    textTransform: "capitalize",
    borderRadius: "10px",
    ...sx,
  };

  if (variant === "outlined") {
    styles = {
      ...styles,
      "&.MuiButton-root": {
        backgroundColor: "#ffffff",
      },
    };
  }

  if (size === "large") {
    styles = {
      ...styles,
      fontSize: "18px",
    };
  }

  return (
    <>
      {loading ? (
        <LoadingButton
          sx={styles}
          loading
          loadingPosition="start"
          startIcon={<SaveIcon />}
          variant={variant}
          {...other}
        ></LoadingButton>
      ) : (
        <MuiButton ref={ref} sx={styles} variant={variant} {...other}>
          {children}
        </MuiButton>
      )}
    </>
  );
});
