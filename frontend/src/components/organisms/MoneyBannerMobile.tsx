import {Typography, useMediaQuery, useTheme} from "@mui/material";

export function MoneyBannerMobile() {
  const theme = useTheme();
  const lgUp: boolean = useMediaQuery(theme.breakpoints.up("lg"));

  return (
      <div className={lgUp ? "hidden": "h-[12rem] flex items-center justify-center"}>
        <div>
          <div className="flex items-center justify-center">
            <img
                className="w-1/6"
                src="https://money-static-files.s3.amazonaws.com/images/dollar.png"
                alt="dollar_logo"
            />
            <Typography color={"darkGreen.main"} variant={"h2"} ml="5px">
              Money
            </Typography>
          </div>
          <div className={"flex justify-center"}>
            <Typography variant={"h6"} color={"darkGreen.main"}>
              Finance tracker
            </Typography>
          </div>
        </div>
      </div>
  )
}
