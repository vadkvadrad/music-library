import { apiClient } from './client';
import {
  Artist,
  Album,
  Song,
  NewArtistRequest,
  UpdateArtistRequest,
  NewAlbumRequest,
  UpdateAlbumRequest,
  NewSongRequest,
  UpdateSongRequest,
  Genre,
  NewGenreRequest,
  UpdateGenreRequest,
} from '../types';

export const musicApi = {
  // Artists
  getArtist: async (id: string): Promise<Artist> => {
    const response = await apiClient.get<Artist>(`/artist/${id}`);
    return response.data;
  },

  createArtist: async (data: NewArtistRequest): Promise<Artist> => {
    const response = await apiClient.post<Artist>('/artist', data);
    return response.data;
  },

  updateArtist: async (id: number, data: UpdateArtistRequest): Promise<Artist> => {
    const response = await apiClient.patch<Artist>(`/artist/${id}`, data);
    return response.data;
  },

  hasArtistPermission: async (id: string): Promise<{ has_permission: boolean }> => {
    const response = await apiClient.get<{ has_permission: boolean }>(`/artist/${id}/has-permission`);
    return response.data;
  },

  // Albums
  getAlbum: async (id: string): Promise<Album> => {
    const response = await apiClient.get<Album>(`/album/${id}`);
    return response.data;
  },

  createAlbum: async (data: NewAlbumRequest): Promise<void> => {
    await apiClient.post('/album', data);
  },

  updateAlbum: async (id: number, data: UpdateAlbumRequest): Promise<void> => {
    await apiClient.patch(`/album/${id}`, data);
  },

  hasAlbumPermission: async (id: string): Promise<{ has_permission: boolean }> => {
    const response = await apiClient.get<{ has_permission: boolean }>(`/album/${id}/has-permission`);
    return response.data;
  },

  // Songs
  getSong: async (id: string): Promise<Song> => {
    const response = await apiClient.get<Song>(`/song/${id}`);
    return response.data;
  },

  createSong: async (albumId: number, data: NewSongRequest): Promise<void> => {
    await apiClient.post(`/song/${albumId}`, data);
  },

  updateSong: async (id: number, data: UpdateSongRequest): Promise<void> => {
    await apiClient.patch(`/song/${id}`, data);
  },

  hasSongPermission: async (id: string): Promise<{ has_permission: boolean }> => {
    const response = await apiClient.get<{ has_permission: boolean }>(`/song/${id}/has-permission`);
    return response.data;
  },

  // Genres
  getAllGenres: async (): Promise<Genre[]> => {
    const response = await apiClient.get<Genre[]>('/genre');
    return response.data;
  },

  getGenre: async (id: string): Promise<Genre> => {
    const response = await apiClient.get<Genre>(`/genre/${id}`);
    return response.data;
  },

  createGenre: async (data: NewGenreRequest): Promise<void> => {
    await apiClient.post('/genre', data);
  },

  updateGenre: async (id: number, data: UpdateGenreRequest): Promise<void> => {
    await apiClient.patch(`/genre/${id}`, data);
  },
};

