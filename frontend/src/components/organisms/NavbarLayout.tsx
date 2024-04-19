import { ReactNode } from "react";
import { Navbar } from "../molecules";
import { Container, useMediaQuery } from "@mui/material";
import { theme } from "../../assets";

type NavbarLayoutProps = {
  children: ReactNode;
};

export function NavbarLayout({ children }: NavbarLayoutProps) {
  const mdUp: boolean = useMediaQuery(theme.breakpoints.up("md"));
  const containerStyles = {
    backgroundColor: "#fafafa",
    width: "auto",
  };

  return (
    <>
      <Navbar />
      <Container
        sx={
          mdUp
            ? { marginLeft: "11rem", ...containerStyles }
            : { ...containerStyles }
        }
        maxWidth={false}
      >
        <div className={"flex max-w-[1200px]"}>{children}</div>
      </Container>
    </>
  );
}
