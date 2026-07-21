export interface RegisterRequest {
    username : string;
    email : string;
    password : string;
}

export interface LoginRequest {
    email : string;
    password: string;
}

export interface AuthResponse{
    message: string;
    access_token: string;
    refresh_token: string;
    token?: string;
}

export interface User {
    id : string;
    username : string;
    email : string;
    avatar?: string;
    plan?: string;
    created_at?: string;
}
export interface ForgotPasswordRequest {
    email : string;
}
export interface ResetPasswordRequest {
    token : string;
    new_password: string;
}

export interface GoogleTokenRequest {
    access_token: string;
}

