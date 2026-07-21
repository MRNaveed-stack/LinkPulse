import apiClient from './client';
import type { PublicProfile } from '../types/public';

export interface UpdateProfileRequest {
  display_name: string;
  username: string;
  email: string;
  bio: string;
  avatar_url: string;
}

export interface Profile {
  user_id: string;
  display_name: string;
  bio: string;
  avatar_url: string;
  created_at: string;
  updated_at: string;
}

export interface ProfileResponse {
  message: string;
  profile: Profile;
}

export interface ChangePasswordRequest {
  current_password: string;
  new_password: string;
}

export const getProfile = (username: string) =>
  apiClient.get<PublicProfile>(`/u/${username}`).then((res) => res.data);

export const updateProfile = (data: UpdateProfileRequest) =>
  apiClient.put<ProfileResponse>('/profile', data).then((res) => res.data);

export const changePassword = (data: ChangePasswordRequest) =>
  apiClient.put<{ message: string }>('/profile/password', data).then((res) => res.data);

export const deleteAccount = () =>
  apiClient.delete<{ message: string }>('/profile').then((res) => res.data);
