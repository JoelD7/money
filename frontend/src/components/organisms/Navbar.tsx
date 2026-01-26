import {
  Alert,
  AlertTitle,
  Button,
  capitalize,
  Divider,
  Drawer,
  IconButton,
  Snackbar,
  useMediaQuery,
  useTheme,
} from "@mui/material";
import HomeIcon from "@mui/icons-material/Home";
import MenuIcon from "@mui/icons-material/Menu";
import AccessTimeFilledIcon from "@mui/icons-material/AccessTimeFilled";
import MonetizationOnIcon from "@mui/icons-material/MonetizationOn";
import SavingsIcon from "@mui/icons-material/Savings";
import NotificationImportantIcon from "@mui/icons-material/NotificationImportant";
import LabelIcon from "@mui/icons-material/Label";
import SettingsIcon from "@mui/icons-material/Settings";
import LogoutIcon from "@mui/icons-material/Logout";
import { ReactNode, useState } from "react";
import { Logo } from "../atoms";
import { setIsAuthenticated } from "../../store";
import { useNavigate } from "@tanstack/react-router";
import { useDispatch } from "react-redux";
import api from "../../api";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { APIError, Credentials, SnackAlert, User } from "../../types";
import { useGetUser } from "../../queries";
import { PeriodSelector } from "../molecules";
import { USER } from "../../queries/keys";
import { AxiosError } from "axios";
import { ValidationError } from "yup";

type NavbarProps = {
  children?: ReactNode;
};

export function Navbar({ children }: NavbarProps) {
  const customWidth = {
    "&.MuiSvgIcon-root": {
      width: "28px",
      height: "28px",
      fill: "#024511",
    },
  };

  const buttonStyle = {
    margin: "5px 0px",
    color: "gray.dark",
    textTransform: "capitalize",
    "&.MuiButton-root": {
      justifyContent: "flex-start",
      width: "100%",
    },
  };

  const [open, setOpen] = useState<boolean>(false);
  const [selectedPeriod, setSelectedPeriod] = useState<string>("");
  const [alert, setAlert] = useState<SnackAlert>({
    open: false,
    type: "success",
    title: "",
  });

  const theme = useTheme();
  const mdUp: boolean = useMediaQuery(theme.breakpoints.up("md"));

  const getUser = useGetUser();

  const user: User | undefined = getUser.data;

  const queryClient = useQueryClient();
  const dispatch = useDispatch();
  const navigate = useNavigate();

  const logoutMutation = useMutation({
    mutationFn: api.logout,
    onSuccess: () => {
      dispatch(setIsAuthenticated(false));
      navigate({ to: "/login" })
        .then(() => { })
        .catch((err) => {
          console.error("Error navigating to /login", err);
        });
    },
  });

  const patchUserMu = useMutation({
    mutationFn: api.patchUser,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [USER] }).then(null, (error) => {
        console.error("Error invalidating users query", error);
      });
    },
    onError: (error) => {
      if (!error) {
        return
      }

      const err = error as AxiosError;
      let title = "Error updating user's current period"
      try {
        const responseError = err.response?.data as APIError;
        title = responseError.message as string
      } catch (e) {
        console.warn("Unexpected error type:", e)
      }

      setAlert({
        open: true,
        type: "error",
        title: title,
      });
    },
  });

  function logout() {
    let username = "";
    if (user) {
      username = user.username;
    }
    const credentials: Credentials = { username: username, password: "" };
    logoutMutation.mutate(credentials);
  }

  function goToHome() {
    const curLocation = window.location.pathname;

    navigate({ to: "/" })
      .then(() => {
        if (curLocation === "/") {
          //Reload the page to reset state
          window.location.reload();
        }
      })
      .catch((err) => {
        console.error("Error navigating to /", err);
      });
  }

  function goToIncome() {
    let route = "/income";
    if (user) {
      route = `/income?period=${user.current_period}`;
    }

    navigate({ to: route })
      .then(() => { })
      .catch((err) => {
        console.error("Error navigating to /income", err);
      });
  }

  function goToSavings() {
    console.error("Navigating to /savings");
    navigate({ to: "/savings" })
      .then(() => { })
      .catch((err) => {
        console.error("Error navigating to /savings", err);
      });
  }

  function onSelectedPeriodChange(newPeriodID: string) {
    setSelectedPeriod(newPeriodID);

    const userToUpdate: User = {
      username: user ? user.username : "",
      remainder: 0,
      current_period: newPeriodID,
    };

    try {
      patchUserMu.mutate(userToUpdate);
    } catch (e) {
      const err = e as ValidationError;
      setAlert({ open: true, type: "error", title: err.errors[0] });
    }
  }

  return (
    <>
      <Snackbar
        open={alert.open}
        onClose={() => setAlert({ ...alert, open: false })}
        autoHideDuration={6000}
        anchorOrigin={{ vertical: "top", horizontal: "right" }}
      >
        <Alert variant={"filled"} severity={alert.type}>
          <AlertTitle>{capitalize(alert.type)}</AlertTitle>
          {alert.title}
        </Alert>
      </Snackbar>

      {/* Mobile menubar */}
      <div
        className={
          mdUp
            ? "hidden"
            : "flex p-4 bg-white-100 flex-row justify-items-center mb-2.5 fixed top-0 z-10 left-0 w-full"
        }
      >
        {children ? children : <Logo />}

        <div className="ml-auto mr-3">
          <IconButton title="Home" sx={{ margin: "5px 0px" }}>
            <HomeIcon sx={customWidth} />
          </IconButton>

          <IconButton sx={{ marginLeft: "15px" }} onClick={() => setOpen(true)}>
            <MenuIcon sx={customWidth} />
          </IconButton>
        </div>
      </div>

      {/*TODO: add back button when route isn't Home*/}
      {/*Title and go back*/}
      {/*<Grid xs={12}>*/}
      {/*    <IconButton>*/}
      {/*        /!*@ts-ignore*!/*/}
      {/*        <ArrowCircleLeftIcon sx={backButtonStyle} color={"darkGreen"}/>*/}
      {/*    </IconButton>*/}
      {/*</Grid>*/}

      {/* Mobile menu drawer */}
      <Drawer anchor="right" open={open} onClose={() => setOpen(false)}>
        <nav
          hidden={mdUp}
          style={{ backgroundColor: "white" }}
          className="flex flex-col h-screen w-44"
        >
          <div className="flex items-center p-4 justify-center w-full">
            <Logo variant="h5" />
          </div>

          <div className={"p-2 w-full"}>
            <PeriodSelector
              period={selectedPeriod}
              onPeriodChange={onSelectedPeriodChange}
            />
          </div>

          <div className="p-3">
            <Divider />
          </div>

          <div className="pl-3">
            <Button
              sx={buttonStyle}
              startIcon={<HomeIcon sx={customWidth} />}
              onClick={() => goToHome()}
            >
              Home
            </Button>

            <Button
              sx={buttonStyle}
              startIcon={<AccessTimeFilledIcon sx={customWidth} />}
            >
              History
            </Button>

            <Button
              sx={buttonStyle}
              startIcon={<NotificationImportantIcon sx={customWidth} />}
            >
              Notifications
            </Button>

            <Button sx={buttonStyle} startIcon={<LabelIcon sx={customWidth} />}>
              Categories
            </Button>

            <Button
              sx={buttonStyle}
              startIcon={<SavingsIcon sx={customWidth} />}
              onClick={() => goToSavings()}
            >
              Savings
            </Button>

            <Button
              sx={buttonStyle}
              startIcon={<MonetizationOnIcon sx={customWidth} />}
              onClick={() => goToIncome()}
            >
              Income
            </Button>
          </div>

          <Divider sx={{ width: "60%", margin: "20px auto" }} />

          <div className="pl-3 h-full">
            <Button
              sx={{ ...buttonStyle, margin: "0px" }}
              startIcon={<SettingsIcon sx={customWidth} />}
            >
              Settings
            </Button>

            <Button
              onClick={() => logout()}
              sx={{ ...buttonStyle, marginTop: "auto", marginBottom: "20px" }}
              startIcon={<LogoutIcon sx={customWidth} />}
            >
              Logout
            </Button>
          </div>
        </nav>
      </Drawer>

      {/* Desktop menu */}
      <div className={"h-[100%] fixed bottom-0 left-0"}>
        <nav
          style={{ backgroundColor: "white" }}
          className={
            mdUp ? "flex flex-col h-screen w-[180px] sticky top-0 mr-3" : "hidden"
          }
        >
          <div className="flex items-center p-4 justify-center w-full">
            <Logo variant="h5" />
          </div>

          {/*Period selector*/}
          <div className={"p-2 w-full"}>
            <PeriodSelector
              period={selectedPeriod}
              onPeriodChange={onSelectedPeriodChange}
            />
          </div>

          <div className="p-3">
            <Divider />
          </div>

          <div className="pl-3">
            <Button
              sx={buttonStyle}
              startIcon={<HomeIcon sx={customWidth} />}
              onClick={() => goToHome()}
            >
              Home
            </Button>

            <Button
              sx={buttonStyle}
              startIcon={<AccessTimeFilledIcon sx={customWidth} />}
            >
              History
            </Button>

            <Button
              sx={buttonStyle}
              startIcon={<NotificationImportantIcon sx={customWidth} />}
            >
              Notifications
            </Button>

            <Button sx={buttonStyle} startIcon={<LabelIcon sx={customWidth} />}>
              Categories
            </Button>

            <Button
              sx={buttonStyle}
              startIcon={<SavingsIcon sx={customWidth} />}
              onClick={() => goToSavings()}
            >
              Savings
            </Button>

            <Button
              sx={buttonStyle}
              startIcon={<MonetizationOnIcon sx={customWidth} />}
              onClick={() => goToIncome()}
            >
              Income
            </Button>
          </div>

          <Divider sx={{ width: "60%", margin: "20px auto" }} />

          <div className="pl-3 h-full">
            <Button
              sx={{ ...buttonStyle, margin: "0px" }}
              startIcon={<SettingsIcon sx={customWidth} />}
            >
              Settings
            </Button>

            <Button
              onClick={() => logout()}
              sx={{ ...buttonStyle, marginTop: "auto", marginBottom: "20px" }}
              startIcon={<LogoutIcon sx={customWidth} />}
            >
              Logout
            </Button>
          </div>
        </nav>
      </div>
    </>
  );
}
