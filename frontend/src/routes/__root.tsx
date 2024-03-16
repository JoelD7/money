import { createRootRoute, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/router-devtools";
import {
  Container,
  createTheme,
  Theme,
  ThemeProvider,
  useMediaQuery,
} from "@mui/material";
import { Colors } from "../assets";

declare module '@mui/material/styles' {
  interface Palette {
    white: Palette['primary'];
    red: Palette['primary'];
    blue: Palette['primary'];
    gray: Palette['primary'];
    darkGreen: Palette['primary'];
    darkerGray: Palette['primary'];
  }

  interface PaletteOptions {
    white?: PaletteOptions['primary'];
    red?: PaletteOptions['primary'];
    blue?: PaletteOptions['primary'];
    gray?: PaletteOptions['primary'];
    darkGreen?: PaletteOptions['primary'];
    darkerGray: PaletteOptions['primary'];
  }

  interface PaletteColor {
    darker?: string;
  }

  interface SimplePaletteColorOptions {
    darker?: string;
  }
}

declare module '@mui/material/Button' {
  interface ButtonPropsColorOverrides {
    white: true;
    red: true;
    blue: true;
    gray: true;
    darkGreen: true;
    darkerGray: true;
  }
}

const theme: Theme = createTheme({
  typography: {
    fontFamily: ["Outfit", "sans-serif"].join(","),
  },
  palette: {
    primary: {
      main: Colors.GREEN,
      darker: Colors.GREEN_DARK,
    },
    secondary: {
      main: Colors.ORANGE,
      contrastText: Colors.WHITE,
    },
    white: {
      main: Colors.WHITE,
      dark: Colors.WHITE_DARK,
      darker: Colors.WHITE_DARKER,
    },
    red: {
      main: Colors.RED,
      dark: Colors.RED_DARK,
      darker: Colors.RED_DARKER,
    },
    blue: {
      main: Colors.BLUE,
      dark: Colors.BLUE_DARK,
      darker: Colors.BLUE_DARKER,
    },
    gray: {
      // Use this color as it is the same as the "bg-zinc-100" Tailwind class
      main: Colors.GRAY,
      dark: Colors.GRAY_DARK,
      darker: Colors.GRAY_DARKER,
      light: Colors.GRAY_LIGHT,
    },
    darkGreen: {
      main: Colors.GREEN_DARK,
    },
    darkerGray: {
      main: Colors.GRAY_DARKER,
    },
  },
});


export const Route = createRootRoute({
  component: () => <Root />,
});

function Root() {
  const mdUp: boolean = useMediaQuery(theme.breakpoints.up("md"));
  const containerStyles = {
    backgroundColor: "#fafafa",
    width: "auto",
  };

  return (
    <>
      <ThemeProvider theme={theme}>
        <div className={"flex max-w-[1200px]"}>
          <Outlet />
        </div>
        <Container
          sx={
            mdUp
              ? { marginLeft: "11rem", ...containerStyles }
              : { ...containerStyles }
          }
          maxWidth={false}
        >

        </Container>
      </ThemeProvider>

      <TanStackRouterDevtools />
    </>
  );
}
