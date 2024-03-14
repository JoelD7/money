import {Alert, AlertTitle, Box, Grid, TextField, Typography} from "@mui/material";
import {Button, MoneyBanner, MoneyBannerMobile} from "../components";
import {useMutation} from "@tanstack/react-query";
import {AxiosError} from "axios";
import {useState} from "react";
import {InputError} from "../types";
import { api } from "../api";

export function Login() {
  const mutation = useMutation({
    mutationFn: api.login,
    onSuccess: () => {
      setErrResponse("");
    },
    onError: (error) => {
      if (error) {
        const err = error as AxiosError;
        setErrResponse(err.response?.data as string);
      }
    },
  });

  const[username, setUsername] = useState<string>("");
  const[password, setPassword] = useState<string>("");

  const [inputErr, setInputErr] = useState<InputError>({
    username: "",
    password: "",
  });

  const [errResponse, setErrResponse] = useState<string>("");

  function login() {
    if (!validateInput()) {
      return;
    }

    mutation.mutate({
      username: username,
      password: password,
    });
  }

  function validateInput(): boolean {
    if (inputErr.username !== "" && inputErr.password !== "") {
      return true
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
    return isValid;
  }

  return (
      <Grid container>
        {/*Green background logo*/}
        <Grid lg={6}>
          <MoneyBanner/>
        </Grid>

        {/*Form and title*/}
        <Grid xs={12} lg={6}>
          {/*Title*/}
          <MoneyBannerMobile/>

          {/*Form*/}
          <Box component="form" height={"100vh"} autoComplete="on">
            <Grid container marginTop={10} justifyContent={"center"}>
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
                      onChange={(e)=> setUsername(e.target.value)}
                      onBlur={validateInput}
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
                      onChange={(e)=> setPassword(e.target.value)}
                      onBlur={validateInput}
                  />
                </div>
              </Grid>

              {/*Button*/}
              <Grid maxWidth={"645px"} xs={12} md={9} paddingTop={1}>
                <div className={"w-11/12 m-auto max-w-[645px]"}>
                  <Button
                      variant={"contained"}
                      loading={mutation.isPending}
                      fullWidth={true}
                      onClick={login}
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
                </div>
              </Grid>
            </Grid>
          </Box>
        </Grid>


      </Grid>
  )
}
