import { z } from 'zod';

export const ProfileUpdateSchema = z.object({
  display_name: z.string().min(2, 'Display name must be at least 2 characters').max(100, 'Display name must be under 100 characters'),
  username: z.string()
    .min(2, 'Username must be at least 2 characters')
    .max(50, 'Username must be under 50 characters')
    .regex(/^[a-zA-Z0-9_]+$/, 'Username can only contain letters, numbers, and underscores'),
  email: z.string().email('Invalid email address'),
  bio: z.string().max(200, 'Bio must be under 200 characters').optional().or(z.literal('')),
  avatar_url: z.string().url('Invalid URL format').optional().or(z.literal('')),
});

export const PasswordChangeSchema = z.object({
  current_password: z.string().min(8, 'Password must be at least 8 characters'),
  new_password: z.string()
    .min(8, 'New password must be at least 8 characters')
    .regex(/[A-Z]/, 'Password must contain at least one uppercase letter')
    .regex(/[a-z]/, 'Password must contain at least one lowercase letter')
    .regex(/[0-9]/, 'Password must contain at least one number'),
  confirm_password: z.string().min(8, 'Password confirmation must be at least 8 characters'),
}).refine((data) => data.new_password === data.confirm_password, {
  message: 'Passwords do not match',
  path: ['confirm_password'],
});
