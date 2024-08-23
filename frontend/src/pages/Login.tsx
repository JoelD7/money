import {
  Alert,
  AlertTitle,
  Box,
  Link,
  TextField,
  Typography,
  useMediaQuery,
} from "@mui/material";
import { Button, MoneyBanner, MoneyBannerMobile } from "../components";
import { useMutation } from "@tanstack/react-query";
import { AxiosError, AxiosResponse } from "axios";
import { ChangeEvent, FormEvent, useState } from "react";
import { APIError, InputError } from "../types";
import api from "../api";
import { Colors, theme } from "../assets";
import Grid from "@mui/material/Unstable_Grid2";
import { useNavigate } from "@tanstack/react-router";
import { AccessTokenResponse } from "../types/other.ts";
import { keys } from "../utils";
import { useDispatch } from "react-redux";
import { setIsAuthenticated } from "../store";

export function Login() {
  const navigate = useNavigate({ from: "/login" });
  const mdUp: boolean = useMediaQuery(theme.breakpoints.up("md"));

  const dispatch = useDispatch();

  const mutation = useMutation({
    mutationFn: api.login,
    onSuccess: (res: AxiosResponse) => {
      const loginResponse: AccessTokenResponse = res.data;
      localStorage.setItem(keys.ACCESS_TOKEN, loginResponse.accessToken);

      setErrResponse("");

      dispatch(setIsAuthenticated(true));

      navigate({ to: "/" })
        .then(() => {})
        .catch((err) => {
          console.log("Error navigation to /", err);
        });
    },
    onError: (error) => {
      if (error) {
        const err = error as AxiosError;
        const responseError = err.response?.data as APIError;

        setErrResponse(responseError.message as string);
      }
    },
  });

  const [username, setUsername] = useState<string>("");
  const [password, setPassword] = useState<string>("");

  const [inputErr, setInputErr] = useState<InputError>({
    username: "",
    password: "",
  });

  const [errResponse, setErrResponse] = useState<string>("");

  function login(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    mutation.reset();

    if (!validateInput()) {
      return;
    }

    mutation.mutate({
      username: username,
      password: password,
    });
  }

  function onUsernameChange(e: ChangeEvent<HTMLInputElement>) {
    setUsername(e.target.value);
    setInputErr({ ...inputErr, username: "" });
  }

  function onPasswordChange(e: ChangeEvent<HTMLInputElement>) {
    setPassword(e.target.value);
    setInputErr({ ...inputErr, password: "" });
  }

  function validateInput(): boolean {
    if (inputErr.username !== "" && inputErr.password !== "") {
      return false;
    }

    let isValid = true;
    const errObj: InputError = {
      username: "",
      password: "",
    };

    if (username === "") {
      errObj.username = "Username is required";

      isValid = false;
    }

    if (password === "") {
      errObj.password = "Password is required";

      isValid = false;
    }

    setInputErr(errObj);
    setErrResponse("");
    return isValid;
  }

  return (
    <div>
      <Grid
        container
        width={"100vw"}
        sx={
          mdUp
            ? { backgroundColor: "#ffffff" }
            : { backgroundColor: "#ffffff", marginLeft: "-40px" }
        }
      >
        {/*Green background logo*/}
        <Grid lg={6}>
          <MoneyBanner />
        </Grid>

        {/*Form and title*/}
        <Grid xs={12} lg={6}>
          {/*Title*/}
          <MoneyBannerMobile />

          {/*Form*/}
          <Box
            component="form"
            onSubmit={login}
            autoComplete="on"
            maxWidth={"645px"}
            margin={"auto"}
          >
            <Grid container marginTop={20} justifyContent={"center"}>
              {/*Input fields*/}
              <Grid xs={12} md={9}>
                <div className={"w-11/12 m-auto max-w-[645px]"}>
                  <Typography textAlign={"center"} variant={"h4"}>
                    Login
                  </Typography>

                  <TextField
                    autoComplete={"on"}
                    margin={"normal"}
                    name={"username"}
                    value={username}
                    fullWidth={true}
                    type={"email"}
                    label={"Email"}
                    variant={"outlined"}
                    error={inputErr.username !== ""}
                    helperText={inputErr.username}
                    required
                    onChange={onUsernameChange}
                  />
                  <TextField
                    autoComplete={"on"}
                    margin={"normal"}
                    name={"password"}
                    value={password}
                    fullWidth={true}
                    type={"password"}
                    label={"Password"}
                    variant={"outlined"}
                    error={inputErr.password !== ""}
                    helperText={inputErr.password}
                    required={true}
                    onChange={onPasswordChange}
                  />
                </div>
              </Grid>

              {/*Button*/}
              <Grid xs={12} md={9} paddingTop={1}>
                <div className={"w-11/12 m-auto max-w-[645px]"}>
                  <Button
                    variant={"contained"}
                    loading={mutation.isPending}
                    type={"submit"}
                    fullWidth={true}
                  >
                    Login
                  </Button>

                  {mutation.isError && (
                    <div className={"p-2"}>
                      <Alert severity="error">
                        <AlertTitle>Error</AlertTitle>
                        {errResponse ? errResponse : mutation.error.message}
                      </Alert>
                    </div>
                  )}

                  {mutation.isSuccess && (
                    <div className={"p-2"}>
                      <Alert severity="success">
                        <AlertTitle>Success</AlertTitle>
                        {"Account successfully created."}
                      </Alert>
                    </div>
                  )}

                  <Typography textAlign={"center"} marginTop={"5px"}>
                    Don't have an account?{" "}
                    <Link color={Colors.BLUE_DARK} href={"/signup"}>
                      Sign up
                    </Link>
                  </Typography>
                </div>
              </Grid>
            </Grid>
          </Box>
        </Grid>
      </Grid>
    </div>
  );
}
