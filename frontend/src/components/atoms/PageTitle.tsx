import { Typography } from "@mui/material";
import { ReactNode } from "react";

type PageTitleProps = {
  children: ReactNode;
};

export function PageTitle({ children }: PageTitleProps) {
  return (
    <Typography variant={"h3"} marginBottom={"1rem"}>
      {children}
    </Typography>
  );
}
