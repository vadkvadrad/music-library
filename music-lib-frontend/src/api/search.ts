import { apiClient } from './client';
import { SearchResult } from '../types';

export const searchApi = {
  search: async (
    query: string,
    types: string[] = ['artist', 'album', 'song'],
    limit: number = 10,
    offset: number = 0
  ): Promise<SearchResult> => {
    const response = await apiClient.get<SearchResult>('/search', {
      params: {
        q: query,
        type: types.join(','),
        limit,
        offset,
      },
    });
    return response.data;
  },
};

