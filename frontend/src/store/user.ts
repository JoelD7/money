import { createSlice } from "@reduxjs/toolkit";
import { Period, User } from "../types/domain.ts";

type userState = {
  user?: User;
  // displayPeriod is the period about which all finances throughout the app are displayed
  displayPeriod?: Period
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
    setDisplayPeriod: (state: userState, action) => {
      state.displayPeriod = { ...action.payload };
    }
  },
});


export const { setUser, setDisplayPeriod } = usersSlice.actions

export default usersSlice.reducer