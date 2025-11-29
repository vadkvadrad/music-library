import { useParams, Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { musicApi } from '../api/music';
import Card from '../components/Card';
import { formatDate, formatDuration } from '../utils/format';
import { Disc, Music } from 'lucide-react';

export default function AlbumPage() {
  const { id } = useParams<{ id: string }>();

  const { data: album, isLoading, error } = useQuery({
    queryKey: ['album', id],
    queryFn: () => musicApi.getAlbum(id!),
    enabled: !!id,
  });

  if (isLoading) {
    return <div className="text-center py-12">Загрузка...</div>;
  }

  if (error || !album) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-red-700">
        Ошибка при загрузке альбома
      </div>
    );
  }

  return (
    <div>
      <div className="mb-8">
        <div className="flex items-start gap-6">
          {album.cover_art_url && (
            <img
              src={album.cover_art_url}
              alt={album.title}
              className="w-64 h-64 object-cover rounded-xl shadow-lg"
            />
          )}
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-4">
              <Disc className="h-8 w-8 text-primary-600" />
              <h1 className="text-4xl font-bold">{album.title}</h1>
            </div>
            <p className="text-gray-600 text-lg mb-2">
              Дата выпуска: {formatDate(album.release_date)}
            </p>
            {album.artist_id && (
              <Link
                to={`/artist/${album.artist_id}`}
                className="text-primary-600 hover:underline"
              >
                Посмотреть артиста
              </Link>
            )}
          </div>
        </div>
      </div>

      {album.songs && album.songs.length > 0 && (
        <div>
          <h2 className="text-2xl font-semibold mb-4 flex items-center gap-2">
            <Music className="h-6 w-6" />
            Песни ({album.songs.length})
          </h2>
          <div className="space-y-2">
            {album.songs.map((song, index) => (
              <Link key={song.id} to={`/song/${song.id}`}>
                <Card>
                  <div className="flex justify-between items-center">
                    <div className="flex items-center gap-4">
                      <span className="text-gray-500 font-mono w-8">
                        {index + 1}
                      </span>
                      <div>
                        <h3 className="text-lg font-semibold">{song.title}</h3>
                        <p className="text-gray-500 text-sm">
                          {formatDuration(song.duration)}
                        </p>
                      </div>
                    </div>
                    <Music className="h-5 w-5 text-primary-600" />
                  </div>
                </Card>
              </Link>
            ))}
          </div>
        </div>
      )}

      {(!album.songs || album.songs.length === 0) && (
        <div className="text-center py-12 text-gray-500">
          В этом альбоме пока нет песен
        </div>
      )}
    </div>
  );
}

