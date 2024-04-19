import { createTheme, Theme } from "@mui/material";
import { Colors } from "./index.ts";

export const theme: Theme = createTheme({
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