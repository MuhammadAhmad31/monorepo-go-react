import { getGetCurrentUserQueryKey, useLogin } from "@/generated/auth/auth";
import type { AuthResponse, LoginRequest } from "@/generated/models";
import { useQueryClient } from "@tanstack/react-query";
import type { AxiosError } from "axios";
import { toast } from "sonner";
import { useNavigate } from "@tanstack/react-router";

export const useAuth = () => {
   const navigate = useNavigate();
  const queryClient = useQueryClient()
  const mutation = useLogin({
    mutation: {
      onMutate: async () => {
        const toastId = toast.loading("Logging in...");
        return { toastId };
      },
      onSuccess: (data: AuthResponse, _variables, context) => {
        if (context?.toastId) toast.dismiss(context.toastId);
        console.log("LOGIN SUCCESS DATA:", data);

        sessionStorage.setItem("authToken", data.token);
        toast.success("Login successful!");
        queryClient.invalidateQueries({ queryKey: getGetCurrentUserQueryKey() });
        navigate({ to: "/dashboard" });
      },
      onError: (error, _, context) => {
        if (context?.toastId) toast.dismiss(context.toastId);
        const axiosError = error as AxiosError;
        if(axiosError.response){
          toast.error(`Login failed: ${(axiosError.response.data as any).message || axiosError.message}`);
        } else {
          toast.error(`Login failed: ${axiosError.message}`);
        }
      },
    },
  });

  const logout = () => {
    sessionStorage.removeItem("authToken");
    queryClient.clear(); 
    navigate({ to: "/" });
  };


  return {
    login: (data: LoginRequest) => mutation.mutateAsync({ data }),
    isLoggingIn: mutation.isPending,
    isSuccess: mutation.isSuccess,
    isError: mutation.isError,
    error: mutation.error,
    reset: mutation.reset,
    logout,
  };
};
