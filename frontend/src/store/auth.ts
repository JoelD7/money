import { createSlice } from "@reduxjs/toolkit";

type authState = {
  isAuthenticated: boolean;
};

const defaultState: authState = {
    isAuthenticated: false
}

export const authSlice = createSlice({
    name: "auth",
    initialState: defaultState,
    reducers: {
        setAuthenticated: (state: authState, action) => {
            state.isAuthenticated = action.payload
        },
    },
})

// actions
export const { setAuthenticated } = authSlice.actions

// selectors
export const selectIsAuthenticated = (state: { auth: authState }) => state.auth.isAuthenticated

export default authSlice.reducer