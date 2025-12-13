import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { searchApi } from '../api/search';
import { Link } from 'react-router-dom';
import Card from '../components/Card';
import Input from '../components/Input';
import Button from '../components/Button';
import { Search as SearchIcon, Music, Disc, Mic } from 'lucide-react';
import { formatDate, formatDuration } from '../utils/format';

export default function SearchPage() {
  const [query, setQuery] = useState('');
  const [searchQuery, setSearchQuery] = useState('');
  const [types, setTypes] = useState<string[]>(['artist', 'album', 'song']);

  const { data, isLoading, error } = useQuery({
    queryKey: ['search', searchQuery, types],
    queryFn: () => searchApi.search(searchQuery, types, 20, 0),
    enabled: searchQuery.length > 0,
  });

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setSearchQuery(query);
  };

  const toggleType = (type: string) => {
    setTypes((prev) =>
      prev.includes(type) ? prev.filter((t) => t !== type) : [...prev, type]
    );
  };

  return (
    <div>
      <h1 className="text-3xl font-bold mb-6">Поиск</h1>

      <form onSubmit={handleSearch} className="mb-6">
        <div className="flex gap-4">
          <div className="flex-1">
            <Input
              placeholder="Поиск артистов, альбомов, песен..."
              value={query}
              onChange={(e) => setQuery(e.target.value)}
            />
          </div>
          <Button type="submit" variant="primary">
            <SearchIcon className="h-5 w-5 mr-2" />
            Найти
          </Button>
        </div>
      </form>

      <div className="flex gap-2 mb-6">
        {['artist', 'album', 'song'].map((type) => (
          <button
            key={type}
            onClick={() => toggleType(type)}
            className={`px-4 py-2 rounded-lg transition-colors ${
              types.includes(type)
                ? 'bg-primary-600 text-white'
                : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
            }`}
          >
            {type === 'artist' && 'Артисты'}
            {type === 'album' && 'Альбомы'}
            {type === 'song' && 'Песни'}
          </button>
        ))}
      </div>

      {isLoading && <div className="text-center py-12">Загрузка...</div>}

      {error && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-red-700">
          Ошибка при выполнении поиска
        </div>
      )}

      {data && searchQuery && (
        <div className="space-y-8">
          {data.artist && data.artist.data && data.artist.data.length > 0 && (
            <div>
              <div className="flex items-center gap-2 mb-4">
                <Mic className="h-6 w-6 text-primary-600" />
                <h2 className="text-2xl font-semibold">
                  Артисты ({data.artist.pagination.total})
                </h2>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {data.artist.data.map((artist) => (
                  <Link key={artist.id} to={`/artist/${artist.id}`}>
                    <Card>
                      <h3 className="text-xl font-semibold mb-2">{artist.name}</h3>
                      <p className="text-gray-600 text-sm line-clamp-2">
                        {artist.description}
                      </p>
                      {artist.formation_year && (
                        <p className="text-gray-500 text-xs mt-2">
                          Основан: {formatDate(artist.formation_year)}
                        </p>
                      )}
                    </Card>
                  </Link>
                ))}
              </div>
            </div>
          )}

          {data.album && data.album.data && data.album.data.length > 0 && (
            <div>
              <div className="flex items-center gap-2 mb-4">
                <Disc className="h-6 w-6 text-primary-600" />
                <h2 className="text-2xl font-semibold">
                  Альбомы ({data.album.pagination.total})
                </h2>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {data.album.data.map((album) => (
                  <Link key={album.id} to={`/album/${album.id}`}>
                    <Card>
                      {album.cover_art_url && (
                        <img
                          src={album.cover_art_url}
                          alt={album.title}
                          className="w-full h-48 object-cover rounded-lg mb-4"
                        />
                      )}
                      <h3 className="text-xl font-semibold mb-2">{album.title}</h3>
                      <p className="text-gray-500 text-sm">
                        {formatDate(album.release_date)}
                      </p>
                    </Card>
                  </Link>
                ))}
              </div>
            </div>
          )}

          {data.song && data.song.data && data.song.data.length > 0 && (
            <div>
              <div className="flex items-center gap-2 mb-4">
                <Music className="h-6 w-6 text-primary-600" />
                <h2 className="text-2xl font-semibold">
                  Песни ({data.song.pagination.total})
                </h2>
              </div>
              <div className="space-y-2">
                {data.song.data.map((song) => (
                  <Link key={song.id} to={`/song/${song.id}`}>
                    <Card>
                      <div className="flex justify-between items-center">
                        <div>
                          <h3 className="text-lg font-semibold">{song.title}</h3>
                          <p className="text-gray-500 text-sm">
                            Длительность: {formatDuration(song.duration)}
                          </p>
                        </div>
                        <Music className="h-6 w-6 text-primary-600" />
                      </div>
                    </Card>
                  </Link>
                ))}
              </div>
            </div>
          )}

          {(!data.artist || !data.artist.data || data.artist.data.length === 0) &&
            (!data.album || !data.album.data || data.album.data.length === 0) &&
            (!data.song || !data.song.data || data.song.data.length === 0) && (
              <div className="text-center py-12 text-gray-500">
                Ничего не найдено
              </div>
            )}
        </div>
      )}

      {!searchQuery && (
        <div className="text-center py-12 text-gray-500">
          Введите запрос для поиска
        </div>
      )}
    </div>
  );
}

