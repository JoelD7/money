import { createSlice } from "@reduxjs/toolkit";

export type authState = {
  isAuthenticated: boolean;
};

const defaultState: authState = {
  isAuthenticated: false,
};

export const authSlice = createSlice({
  name: "auth",
  initialState: defaultState,
  reducers: {
    setIsAuthenticated: (state: authState, action) => {
      state.isAuthenticated = action.payload;
    },
  },
});

// actions
export const { setIsAuthenticated } = authSlice.actions;

export default authSlice.reducer;
