import {useMutation} from '@tanstack/react-query';
import {register} from '../api/auth';
import {useAuthStore} from '../store/authStore';
import type { RegisterRequest } from '../types/auth';

export const useRegister = () => {
    const {login: storeLogin} = useAuthStore();
    return useMutation({
        mutationFn: (data : RegisterRequest) => register(data),
        onSuccess : (response) => {
            const {access_token, refresh_token} = response.data;
            const user = JSON.parse(atob(access_token.split('.')[1]));
            storeLogin(
                {id:user.sub, username:user.username, email:user.email},
                {accessToken: access_token, refreshToken: refresh_token}
            );
        }
    })
}