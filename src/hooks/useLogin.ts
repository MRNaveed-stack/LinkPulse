import {useMutation} from "@tanstack/react-query";
import {login} from "../api/auth";
import {useAuthStore} from "../store/authStore";
import type { LoginRequest } from "../types/auth";

export const useLogin = () => {
    const {login: storelogin} = useAuthStore();
    return useMutation({
        mutationFn: (data: LoginRequest) => login(data),
        onSuccess : (response) => {
     const {access_token, refresh_token} = response.data;
     const user = JSON.parse(atob(access_token.split('.')[1]));
     storelogin(
        {id:user.sub, username:user.username, email:user.email},
        {accessToken: access_token, refreshToken: refresh_token}
     );
        },
    });
};