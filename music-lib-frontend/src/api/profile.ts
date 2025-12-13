import { apiClient } from './client';
import { Profile, NewProfileRequest, UpdateProfileRequest } from '../types';

export const profileApi = {
  getProfile: async (): Promise<Profile> => {
    const response = await apiClient.get<Profile>('/profile');
    return response.data;
  },

  createProfile: async (data: NewProfileRequest): Promise<void> => {
    await apiClient.post('/profile', data);
  },

  updateProfile: async (data: UpdateProfileRequest): Promise<void> => {
    await apiClient.patch('/profile', data);
  },
};

