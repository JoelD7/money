import { createRootRoute, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/router-devtools";
import { ThemeProvider, useMediaQuery } from "@mui/material";
import { theme } from "../assets";
import { Provider } from "react-redux";
import { persistor, store } from "../store";
import { PersistGate } from "redux-persist/integration/react";
import { LocalizationProvider } from "@mui/x-date-pickers";
import { AdapterDayjs } from "@mui/x-date-pickers/AdapterDayjs";
import "dayjs/locale/en-gb";

declare module "@mui/material/styles" {
  interface Palette {
    white: Palette["primary"];
    red: Palette["primary"];
    blue: Palette["primary"];
    gray: Palette["primary"];
    darkGreen: Palette["primary"];
    darkerGray: Palette["primary"];
  }

  interface PaletteOptions {
    white?: PaletteOptions["primary"];
    red?: PaletteOptions["primary"];
    blue?: PaletteOptions["primary"];
    gray?: PaletteOptions["primary"];
    darkGreen?: PaletteOptions["primary"];
    darkerGray: PaletteOptions["primary"];
  }

  interface PaletteColor {
    darker?: string;
  }

  interface SimplePaletteColorOptions {
    darker?: string;
  }
}

declare module "@mui/material/Button" {
  interface ButtonPropsColorOverrides {
    white: true;
    red: true;
    blue: true;
    gray: true;
    darkGreen: true;
    darkerGray: true;
  }
}

export const Route = createRootRoute({
  component: () => <Root />,
});

function Root() {
  return (
    <>
      <LocalizationProvider dateAdapter={AdapterDayjs} adapterLocale={"en-gb"}>
        <ThemeProvider theme={theme}>
          <Provider store={store}>
            <PersistGate loading={null} persistor={persistor}>
              <div className={"bg-zinc-50 h-full pb-14 min-h-screen"}>
                <PageContent />
              </div>
            </PersistGate>
          </Provider>
        </ThemeProvider>
      </LocalizationProvider>

      <TanStackRouterDevtools />
    </>
  );
}

function PageContent() {
  const mdUp: boolean = useMediaQuery(theme.breakpoints.up("md"));

  if (mdUp) {
    return (
      <div
        className={"flex justify-center"}
        style={mdUp ? {} : { flexDirection: "column" }}
      >
        <Outlet />
      </div>
    );
  }

  return (
    <div
      className={"flex max-w-[1600px]"}
      style={mdUp ? {} : { flexDirection: "column" }}
    >
      <div className={"px-10"}>{<Outlet />}</div>
    </div>
  );
}
