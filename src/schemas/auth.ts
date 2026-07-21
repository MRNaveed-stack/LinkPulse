import {z} from 'zod';

export const loginSchema = z.object({
    email : z.email('Invalid email'),
    password : z.string().min(6,'Password must be at least 6 characters'),
});

export const registerSchema = z.object({
    username: z.string().min(3,'Username must be at least 3 characters'),
    email : z.email('Invalid email'),
    password : z.string().min(6,'Password must be at least 6 characters')
});

export const forgotPasswordSchema = z.object({
    email : z.email('Invalid email'),
});

export const resetPasswordSchema = z.object({
    token: z.string().min(1,'Token is required'),
    new_password : z.string().min(6, 'Password must be at least 6 characters')
});