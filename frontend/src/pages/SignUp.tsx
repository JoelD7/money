import {
  Alert,
  AlertTitle,
  Box,
  Link,
  TextField,
  Typography,
} from "@mui/material";
import { Button } from "../components";
import { Colors } from "../assets";
import { useMutation } from "@tanstack/react-query";
import { api } from "../api";
import { ChangeEvent, useState } from "react";
import { SignUpUser } from "../types";

type InputError = {
  username?: string;
  password?: string;
};

export function SignUp() {
  const mutation = useMutation({
    mutationFn: api.signUp,
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
    <Box component="form" autoComplete="on">
      {/*Title*/}
      <div className={"h-[12rem] flex items-center"}>
        <div>
          <div className="flex items-center justify-center">
            <img
              className="w-1/6"
              src="https://money-static-files.s3.amazonaws.com/images/dollar.png"
              alt="dollar_logo"
            />
            <Typography color={"darkGreen.main"} variant={"h2"} ml="5px">
              Money
            </Typography>
          </div>
          <div className={"flex justify-center"}>
            <Typography variant={"h6"} color={"darkGreen.main"}>
              Finance tracker
            </Typography>
          </div>
        </div>
      </div>

      {/*Input fields*/}
      <div className={"w-11/12 m-auto"}>
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

      {/*Button*/}
      <div className={"w-11/12 m-auto pt-10"}>
        <Button variant={"contained"} fullWidth={true} onClick={signUp}>
          Sign up
        </Button>

        {mutation.isError && (
          <div className={"p-2"}>
            <Alert severity="error">
              <AlertTitle>Error</AlertTitle>
              {mutation.error.message}
            </Alert>
          </div>
        )}

        <Typography textAlign={"center"}>
          Already signed up?{" "}
          <Link color={Colors.BLUE_DARK} target={"_blank"} href={"/login"}>
            Login
          </Link>
        </Typography>
      </div>
    </Box>
  );
}
