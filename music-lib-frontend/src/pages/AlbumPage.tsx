import { useParams, Link, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { musicApi } from '../api/music';
import { useAuth } from '../contexts/AuthContext';
import Card from '../components/Card';
import Button from '../components/Button';
import { formatDate, formatDuration } from '../utils/format';
import { Disc, Music, Edit, Plus } from 'lucide-react';

export default function AlbumPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();

  const { data: album, isLoading, error } = useQuery({
    queryKey: ['album', id],
    queryFn: () => musicApi.getAlbum(id!),
    enabled: !!id,
  });

  const { data: permission } = useQuery({
    queryKey: ['album-permission', id],
    queryFn: () => musicApi.hasAlbumPermission(id!),
    enabled: !!id && isAuthenticated,
  });

  const canEdit = permission?.has_permission || false;

  if (isLoading) {
    return (
      <div className="text-center py-12">
        <div className="inline-block w-8 h-8 border-4 border-primary-600/30 border-t-primary-600 rounded-full animate-spin" />
      </div>
    );
  }

  if (error || !album) {
    return (
      <Card className="bg-red-500/20 border border-red-500/30">
        <p className="text-red-400">Ошибка при загрузке альбома</p>
      </Card>
    );
  }

  return (
    <div className="animate-fade-in">
      <div className="mb-8">
        <div className="flex items-start gap-6">
          {album.cover_art_url && (
            <img
              src={album.cover_art_url}
              alt={album.title}
              className="w-64 h-64 object-cover rounded-2xl shadow-large"
            />
          )}
          <div className="flex-1">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-3">
                <div className="p-2 rounded-xl bg-primary-600/20">
                  <Disc className="h-8 w-8 text-primary-400" />
                </div>
                <h1 className="text-4xl font-bold text-white">{album.title}</h1>
              </div>
              {canEdit && (
                <div className="flex gap-2">
                  <Button
                    variant="secondary"
                    onClick={() => navigate(`/create-album`, { state: { albumId: album.id } })}
                  >
                    <Edit className="h-4 w-4 mr-2" />
                    Редактировать
                  </Button>
                  <Button
                    variant="primary"
                    onClick={() => navigate(`/create-song/${album.id}`)}
                  >
                    <Plus className="h-4 w-4 mr-2" />
                    Добавить песню
                  </Button>
                </div>
              )}
            </div>
            <p className="text-gray-400 text-lg mb-4">
              Дата выпуска: {formatDate(album.release_date)}
            </p>
            {album.artist_id && (
              <Link
                to={`/artist/${album.artist_id}`}
                className="text-primary-400 hover:text-primary-300 transition-colors"
              >
                Посмотреть артиста →
              </Link>
            )}
          </div>
        </div>
      </div>

      {album.songs && album.songs.length > 0 && (
        <div>
          <div className="flex items-center gap-3 mb-6">
            <Music className="h-6 w-6 text-primary-400" />
            <h2 className="text-2xl font-semibold text-white">
              Песни ({album.songs.length})
            </h2>
          </div>
          <div className="space-y-2">
            {album.songs.map((song, index) => (
              <Link key={song.id} to={`/song/${song.id}`}>
                <Card hover>
                  <div className="flex justify-between items-center">
                    <div className="flex items-center gap-4">
                      <span className="text-gray-500 font-mono w-8 text-center">
                        {index + 1}
                      </span>
                      <div>
                        <h3 className="text-lg font-semibold text-white">{song.title}</h3>
                        <p className="text-gray-400 text-sm">
                          {formatDuration(song.duration)}
                        </p>
                      </div>
                    </div>
                    <Music className="h-5 w-5 text-primary-400" />
                  </div>
                </Card>
              </Link>
            ))}
          </div>
        </div>
      )}

      {(!album.songs || album.songs.length === 0) && (
        <Card className="text-center py-12">
          <div className="p-4 rounded-xl bg-gray-800/50 w-fit mx-auto mb-4">
            <Music className="h-12 w-12 text-gray-500" />
          </div>
          <p className="text-gray-400">В этом альбоме пока нет песен</p>
        </Card>
      )}
    </div>
  );
}

