import { ReactNode } from "react";
import { Typography } from "@mui/material";

type TableHeaderProps = {
  icon?: ReactNode;
  headerName: string;
};

export function TableHeader({ icon, headerName }: TableHeaderProps) {
  return (
    <div className={"flex items-center h-full"}>
      {icon}
      <Typography variant={"h6"} marginLeft={icon ? "5px" : "0px"}>
        {headerName}
      </Typography>
    </div>
  );
}
