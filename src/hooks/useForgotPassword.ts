import { useMutation } from '@tanstack/react-query';
import {forgotPassword} from '../api/auth';
import type { ForgotPasswordRequest } from '../types/auth';

export const useForgotPassword = () => {
    return useMutation({
        mutationFn: (data: ForgotPasswordRequest) => forgotPassword(data),
    });
};