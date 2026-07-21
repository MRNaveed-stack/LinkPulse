import apiClient from './client';
import type { PublicProfile } from '../types/public';

export const getPublicProfile = (username: string) =>
  apiClient.get<PublicProfile>(`/u/${username}`).then((res) => res.data);