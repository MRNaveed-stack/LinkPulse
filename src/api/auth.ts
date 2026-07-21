import apiClient from './client';
import type {
    RegisterRequest,
    LoginRequest,
    ForgotPasswordRequest,
    ResetPasswordRequest,
    GoogleTokenRequest,
    AuthResponse,
    User,
} from '../types/auth';

export const register = (data: RegisterRequest) => apiClient.post<AuthResponse>('/auth/register',data);
export const login = (data: LoginRequest) => apiClient.post<AuthResponse>('/auth/login',data);

export const forgotPassword = (data: ForgotPasswordRequest) =>
  apiClient.post<{ message: string; token?: string }>(
    '/auth/forgot-password',
    data
  );

export const resetPassword = (data: ResetPasswordRequest) =>
  apiClient.post<{ message: string }>('/auth/reset-password', data);

export const googleLogin = (data: GoogleTokenRequest) =>
  apiClient.post<AuthResponse>('/auth/google/token', data);

export const refreshToken = (refreshToken: string) =>
  apiClient.post<AuthResponse>('/auth/refresh', {
    refresh_token: refreshToken,
  });

export const getCurrentUser = () =>
  apiClient.get<User>('/me'); 


