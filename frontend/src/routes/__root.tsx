import { createRootRoute, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/router-devtools";
import { ThemeProvider } from "@mui/material";
import { theme } from "../assets";
import { Provider } from "react-redux";
import { persistor, store } from "../store";
import { PersistGate } from "redux-persist/integration/react";

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
      <ThemeProvider theme={theme}>
        <Provider store={store}>
          <PersistGate loading={null} persistor={persistor}>
            <Outlet />
          </PersistGate>
        </Provider>
      </ThemeProvider>

      <TanStackRouterDevtools />
    </>
  );
}
