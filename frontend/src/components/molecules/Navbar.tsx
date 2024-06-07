import {
  Button,
  Divider,
  Drawer,
  IconButton,
  useMediaQuery,
  useTheme,
} from "@mui/material";
import HomeIcon from "@mui/icons-material/Home";
import MenuIcon from "@mui/icons-material/Menu";
import AccessTimeFilledIcon from "@mui/icons-material/AccessTimeFilled";
import NotificationImportantIcon from "@mui/icons-material/NotificationImportant";
import LabelIcon from "@mui/icons-material/Label";
import AutoStoriesIcon from "@mui/icons-material/AutoStories";
import SettingsIcon from "@mui/icons-material/Settings";
import LogoutIcon from "@mui/icons-material/Logout";
import { ReactNode, useState } from "react";
import { Logo } from "../atoms";
import { setIsAuthenticated } from "../../store";
import { useNavigate } from "@tanstack/react-router";
import { useDispatch } from "react-redux";
import api from "../../api";
import { useMutation, useQuery } from "@tanstack/react-query";
import { Credentials, User } from "../../types";

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
  const theme = useTheme();
  const mdUp: boolean = useMediaQuery(theme.breakpoints.up("md"));

  const getUserQuery = useQuery({
    queryKey: ["user"],
    queryFn: () => api.getUser(),
    retry: false,
    staleTime: 1000,
    refetchOnWindowFocus: false,
  });

  const user: User | undefined = getUserQuery.data?.data;

  const dispatch = useDispatch();
  const navigate = useNavigate();

  const logoutMutation = useMutation({
    mutationFn: api.logout,
    onSuccess: () => {
      dispatch(setIsAuthenticated(false));
      navigate({ to: "/login" })
        .then(() => {})
        .catch((err) => {
          console.error("Error navigating to /login", err);
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

  return (
    <>
      <div
        className={
          mdUp
            ? "hidden"
            : //   mx-[-40px] to override the padding of the parent div. We don't want the navbar to be centered like the
              //   contents of a page
              "flex p-4 bg-white-100 flex-row justify-items-center mx-[-40px] mb-2.5"
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

      <Drawer anchor="right" open={open} onClose={() => setOpen(false)}>
        <nav
          hidden={mdUp}
          style={{ backgroundColor: "white" }}
          className="flex flex-col h-screen w-44"
        >
          <div className="flex items-center p-4 justify-center w-full">
            <Logo variant="h5" />
          </div>

          <div className="pl-3">
            <Button sx={buttonStyle} startIcon={<HomeIcon sx={customWidth} />}>
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
              startIcon={<AutoStoriesIcon sx={customWidth} />}
            >
              Savings
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

      <div>
        <nav
          style={{ backgroundColor: "white" }}
          className={
            mdUp
              ? "flex flex-col h-screen w-[180px] sticky top-0 mr-3"
              : "hidden"
          }
        >
          <div className="flex items-center p-4 justify-center w-full">
            <Logo variant="h5" />
          </div>

          <div className="pl-3">
            <Button sx={buttonStyle} startIcon={<HomeIcon sx={customWidth} />}>
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
              startIcon={<AutoStoriesIcon sx={customWidth} />}
            >
              Savings
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
