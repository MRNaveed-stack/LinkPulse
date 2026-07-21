import apiClient from './client';
import type {
  LinkListResponse,
  CreateLinkRequest,
  UpdateLinkRequest,
  ToggleStatusRequest,
  Link,
  LinkStatusResponse,
} from '../types/link';

export const getLinks = () =>
  apiClient.get<LinkListResponse>('/links').then((res) => res.data.links);

export const createLink = (data: CreateLinkRequest) =>
  apiClient.post('/links', data).then((res) => res.data);

export const updateLink = (id: string, data: UpdateLinkRequest) =>
  apiClient.put<{ message: string; link: Link }>(`/links/${id}`, data).then((res) => res.data);

export const deleteLink = (id: string) =>
  apiClient.delete(`/links/${id}`).then((res) => res.data);

export const toggleStatus = (id: string, data: ToggleStatusRequest) =>
  apiClient.patch<LinkStatusResponse>(`/links/${id}/status`, data).then((res) => res.data);