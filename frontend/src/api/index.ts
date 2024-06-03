import { login, logout, signUp } from "./auth.ts";
import { getUser } from "./money-api.ts";

const api = {
  signUp,
  getUser,
  login,
  logout,
};

export default api;
