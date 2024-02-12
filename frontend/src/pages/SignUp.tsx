import { Box, Link, TextField, Typography } from "@mui/material";
import { Button } from "../components";
import { Colors } from "../assets";
import { useMutation } from "@tanstack/react-query";
import { api } from "../api";
import { ChangeEvent, useState } from "react";
import { SignUpUser } from "../types";

type SignUpError = {
  username: string;
  password: string;
  fullname: string;
};

export function SignUp() {
  const mutation = useMutation({
    mutationFn: api.signUp,
  });

  const [signUpUser, setSignUpUser] = useState<SignUpUser | undefined>();
  const [error, setError] = useState<SignUpError | undefined>();

  function onInputChange(
    e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
  ) {
    if (!signUpUser) {
      setSignUpUser({
        [e.target.name]: e.target.value,
      } as SignUpUser);

      return;
    }

    setSignUpUser({
      ...signUpUser,
      [e.target.name]: e.target.value,
    });
  }

  function signUp() {
    if (!signUpUser) {
      setError({
        username: "Email is required",
        password: "Password is required",
        fullname: "",
      });

      return;
    }

    if (signUpUser.username === "" || signUpUser.password === "") {
      setError({
        username: signUpUser.username === "" ? "Email is required" : "",
        password: signUpUser.password === "" ? "Password is required" : "",
        fullname: "",
      });

      return false;
    }

    mutation.mutate({
      username: signUpUser.username,
      password: signUpUser.password,
      fullname: signUpUser.fullname,
    });
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
          value={signUpUser ? signUpUser.fullname : ""}
          fullWidth={true}
          label={"Full name"}
          variant={"outlined"}
          onChange={onInputChange}
        />
        <TextField
          autoComplete={"on"}
          margin={"normal"}
          name={"username"}
          value={signUpUser ? signUpUser.username : ""}
          fullWidth={true}
          type={"email"}
          label={"Email"}
          variant={"outlined"}
          error={error ? error.username === "" : false}
          helperText={error ? error.username : ""}
          required
          onChange={onInputChange}
        />
        <TextField
          autoComplete={"on"}
          margin={"normal"}
          name={"password"}
          value={signUpUser ? signUpUser.password : ""}
          fullWidth={true}
          type={"password"}
          label={"Password"}
          variant={"outlined"}
          error={error ? error.password === "" : false}
          helperText={error ? error.password : ""}
          required={true}
          onChange={onInputChange}
        />
      </div>

      {/*Button*/}
      <div className={"w-11/12 m-auto pt-10"}>
        <Button variant={"contained"} fullWidth={true} onClick={signUp}>
          Sign up
        </Button>
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
