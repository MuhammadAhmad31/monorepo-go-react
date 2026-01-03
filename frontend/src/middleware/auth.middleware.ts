import { redirect, type ContextOptions } from "@tanstack/react-router";

export const AuthMiddleware = (context: ContextOptions<any, any>) => {
  const authToken = sessionStorage.getItem("authToken");
  const filePath = context.location.pathname;
  try {
    if (!authToken && filePath !== "/") {
      return redirect({
        to: "/",
        search: context.location.search,
      });
    }

    if (authToken && filePath === "/") {
      return redirect({
        to: "/dashboard",
        search: context.location.search,
      });
    }
  } catch (error) {
    if (!authToken) {
      return redirect({
        to: "/",
        search: context.location.search,
      });
    }
  }

    return;
};
