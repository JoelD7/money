import {Typography, useMediaQuery, useTheme} from "@mui/material";

export function MoneyBanner() {
  const theme = useTheme();
  const lgUp: boolean = useMediaQuery(theme.breakpoints.up("lg"));

  return (
    <div
      className={
        lgUp
          ? "flex items-center justify-center h-lvh bg-[#024511] rounded-e-3xl"
          : "hidden"
      }
    >
      <div>
        <div className="flex items-center justify-center">
          <img
            className="w-1/6"
            src="https://money-static-files.s3.amazonaws.com/images/dollar.png"
            alt="dollar_logo"
          />
          <Typography color={"white.main"} variant={"h2"} ml="5px">
            Money
          </Typography>
        </div>
        <div className={"flex justify-center"}>
          <Typography variant={"h6"} color={"white.main"}>
            Finance tracker
          </Typography>
        </div>
      </div>
    </div>
  );
}
