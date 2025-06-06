import { ReactNode } from "react";
import { useMediaQuery } from "@mui/material";
import { theme } from "../../assets";

type ContainerProps = {
  children?: ReactNode;
};

export function Container({ children }: ContainerProps) {
  const mdUp: boolean = useMediaQuery(theme.breakpoints.up("md"));

  // Desktop
  if (mdUp) {
    return <div className={"max-w-[1600px] w-[100%] pl-48 pt-8"}>{children}</div>;
  }

  // Mobile
  return <div className={"w-[100%] pt-24"}>{children}</div>;
}
