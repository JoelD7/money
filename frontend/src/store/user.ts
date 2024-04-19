import { createSlice } from "@reduxjs/toolkit";
import { User } from "../types/domain.ts";

type userState = {
  user?: User;
};

const defaultState: userState = {

}

export const usersSlice = createSlice({
  name: "users",
  initialState: defaultState,
  reducers: {
    setUser: (state: userState, action) => {
        state.user = action.payload
    },
  },
});


export const { setUser } = usersSlice.actions

export default usersSlice.reducer