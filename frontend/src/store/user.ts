import { createSlice } from "@reduxjs/toolkit";
import { Period, User } from "../types/domain.ts";

type userState = {
  user?: User;
  // selectedPeriod is the period about which all finances throughout the app are displayed
  selectedPeriod?: Period
};

const defaultState: userState = {

}

export const usersSlice = createSlice({
  name: "users",
  initialState: defaultState,
  reducers: {
    setUser: (state: userState, action) => {
      state.user = action.payload
      if (action.payload?.current_period) {
        state.selectedPeriod = action.payload.current_period;
      }
    },
    setSelectedPeriod: (state: userState, action) => {
      state.selectedPeriod = action.payload;
    }
  },
});


export const { setUser, setSelectedPeriod } = usersSlice.actions

export default usersSlice.reducer