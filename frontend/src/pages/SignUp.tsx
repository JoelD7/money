import {
  Alert,
  AlertTitle,
  Box,
  Grid,
  Link,
  TextField,
  Typography,
} from "@mui/material";
import {Button, MoneyBanner, MoneyBannerMobile} from "../components";
import { Colors } from "../assets";
import { useMutation } from "@tanstack/react-query";
import { api } from "../api";
import { ChangeEvent, useState } from "react";
import {InputError, SignUpUser} from "../types";
import { AxiosError } from "axios";

export function SignUp() {
  const mutation = useMutation({
    mutationFn: api.signUp,
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

  const [signUpUser, setSignUpUser] = useState<SignUpUser>({
    username: "",
    password: "",
    fullname: "",
  });
  const [inputErr, setInputErr] = useState<InputError>({
    username: "",
    password: "",
  });

  const [errResponse, setErrResponse] = useState<string>("");

  function onInputChange(
    e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
  ) {
    resetError(e);
    setSignUpUser({
      ...signUpUser,
      [e.target.name]: e.target.value,
    });
  }

  function resetError(e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) {
    setInputErr({
      ...inputErr,
      [e.target.name]: "",
    });
  }

  function signUp() {
    if (!validateInput()) {
      return;
    }

    mutation.mutate({
      username: signUpUser.username,
      password: signUpUser.password,
      fullname: signUpUser.fullname,
    });
  }

  function validateInput(): boolean {
    let isValid = true;
    const errObj: InputError = { ...inputErr };

    if (signUpUser.username === "") {
      errObj.username = "Username is required";

      isValid = false;
    }

    if (signUpUser.password === "") {
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
            <Grid container marginTop={20} justifyContent={"center"}>
              {/*Input fields*/}
              <Grid xs={12} md={9}>
                <div className={"w-11/12 m-auto max-w-[645px]"}>
                  <Typography textAlign={"center"} variant={"h4"}>
                    Sign up
                  </Typography>

                  <TextField
                      margin={"normal"}
                      name={"fullname"}
                      value={signUpUser.fullname}
                      fullWidth={true}
                      label={"Full name"}
                      variant={"outlined"}
                      onChange={onInputChange}
                  />
                  <TextField
                      autoComplete={"on"}
                      margin={"normal"}
                      name={"username"}
                      value={signUpUser.username}
                      fullWidth={true}
                      type={"email"}
                      label={"Email"}
                      variant={"outlined"}
                      error={inputErr.username !== ""}
                      helperText={inputErr.username}
                      required
                      onChange={onInputChange}
                  />
                  <TextField
                      autoComplete={"on"}
                      margin={"normal"}
                      name={"password"}
                      value={signUpUser.password}
                      fullWidth={true}
                      type={"password"}
                      label={"Password"}
                      variant={"outlined"}
                      error={inputErr.password !== ""}
                      helperText={inputErr.password}
                      required={true}
                      onChange={onInputChange}
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
                      onClick={signUp}
                  >
                    Sign up
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
                    Already signed up?{" "}
                    <Link color={Colors.BLUE_DARK} target={"_blank"} href={"/login"}>
                      Login
                    </Link>
                  </Typography>
                </div>
              </Grid>
            </Grid>
        </Box>
      </Grid>


    </Grid>
  );
}
