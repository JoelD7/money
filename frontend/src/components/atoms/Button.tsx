import {
  Button as MuiButton,
  ButtonProps,
  SxProps,
  Theme,
} from "@mui/material";
import { ReactNode } from "react";
import LoadingButton from "@mui/lab/LoadingButton";
import SaveIcon from "@mui/icons-material/Save";

type CustomButtonProps = {
  children: ReactNode;
  loading?: boolean;
} & ButtonProps;

export function Button(props: CustomButtonProps) {
  const { sx, loading, variant, children, ...other } = props;
  let styles: SxProps<Theme> | undefined = {
    textTransform: "capitalize",
    borderRadius: "1rem",
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
        <MuiButton
          sx={styles}
          /*@ts-ignore*/
          variant={variant}
          {...other}
        >
          {children}
        </MuiButton>
      )}
    </>
  );
}
