import { combineReducers, configureStore } from "@reduxjs/toolkit";
import usersReducer from "./user.ts";
import authReducer from "./auth.ts";
import { persistReducer, persistStore } from "redux-persist";
import storage from "redux-persist/lib/storage";
import {
  FLUSH,
  PAUSE,
  PERSIST,
  PURGE,
  REGISTER,
  REHYDRATE,
} from "redux-persist/es/constants";

export * from "./user.ts";
export * from "./auth.ts";

const persistConfig = {
  key: "root",
  storage,
};

const baseReducer = combineReducers({ usersReducer, authReducer });

const persistedReducer = persistReducer(persistConfig, baseReducer);

export const store = configureStore({
  reducer: persistedReducer,
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        // To prevent the error "A non serializable value was detected in state" from redus-persis
        ignoredActions: [FLUSH, REHYDRATE, PAUSE, PERSIST, PURGE, REGISTER],
      },
    }),
});

export const persistor = persistStore(store);
