// Auth types
export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
}

export interface VerifyRequest {
  session_id: string;
  code: string;
}

export interface LoginResponse {
  jwt_token: string;
}

export interface RegisterResponse {
  session_id: string;
}

export interface VerifyResponse {
  jwt_token: string;
}

// Music types
export interface Artist {
  id: number;
  name: string;
  description: string;
  formation_year: string;
  albums?: Album[];
}

export interface Album {
  id: number;
  title: string;
  release_date: string;
  cover_art_url: string;
  songs?: Song[];
  artist_id?: number;
}

export interface Song {
  id: number;
  title: string;
  album_id: number;
  duration: number;
  file_path: string;
  lyrics?: Lyrics;
}

export interface Lyrics {
  text: Couplet[];
}

export interface Couplet {
  number: number;
  couplet: string;
}

export interface Genre {
  id: number;
  name: string;
}

// Profile types
export interface Profile {
  user_id: number;
  user_name: string;
  user_email: string;
  bio: string;
  avatar_url: string;
  created_at: string;
  updated_at: string;
  artist?: {
    id: number;
    name: string;
    description: string;
    formation_year: string;
    albums?: Album[];
  };
}

export interface NewProfileRequest {
  bio: string;
  avatar_url: string;
}

// Request types
export interface NewArtistRequest {
  artist_name: string;
  description: string;
  formation_year: string;
}

export interface UpdateArtistRequest {
  artist_name?: string;
  description?: string;
  formation_year?: string;
}

export interface NewAlbumRequest {
  title: string;
  release_date: string;
  cover_art_url: string;
}

export interface NewSongRequest {
  title: string;
  genres: { genre_id: number }[];
  duration_sec: number;
  file_path: string;
  lyrics: {
    text: { couplet: string }[];
  };
}

export interface NewGenreRequest {
  genre_name: string;
}

export interface UpdateGenreRequest {
  genre_name_update: string;
}

// Search types
export interface Pagination {
  limit: number;
  offset: number;
  total: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: Pagination;
}

export interface SearchResult {
  artist?: PaginatedResponse<Artist>;
  album?: PaginatedResponse<Album>;
  song?: PaginatedResponse<Song>;
}

// Error types
export interface ApiError {
  error: string;
  message?: string;
}

