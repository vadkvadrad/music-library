import { useParams, Link, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { musicApi } from '../api/music';
import { useAuth } from '../contexts/AuthContext';
import Card from '../components/Card';
import Button from '../components/Button';
import { formatDate } from '../utils/format';
import { Mic, Disc, Edit, Plus } from 'lucide-react';

export default function ArtistPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();

  const { data: artist, isLoading, error } = useQuery({
    queryKey: ['artist', id],
    queryFn: () => musicApi.getArtist(id!),
    enabled: !!id,
  });

  const { data: permission } = useQuery({
    queryKey: ['artist-permission', id],
    queryFn: () => musicApi.hasArtistPermission(id!),
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

  if (error || !artist) {
    return (
      <div className="mb-6 p-4 bg-red-500/20 border border-red-500/30 rounded-xl text-red-400">
        Ошибка при загрузке артиста
      </div>
    );
  }

  return (
    <div className="animate-fade-in">
      <div className="mb-8">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-4">
            <Mic className="h-12 w-12 text-primary-400" />
            <div>
              <h1 className="text-4xl font-bold text-white">{artist.name}</h1>
              {artist.formation_year && (
                <p className="text-gray-400 mt-1">
                  Основан: {formatDate(artist.formation_year)}
                </p>
              )}
            </div>
          </div>
          {canEdit && (
            <div className="flex gap-3">
              <Button
                variant="secondary"
                onClick={() => navigate('/create-album')}
              >
                <Plus className="h-4 w-4 mr-2" />
                Добавить альбом
              </Button>
              <Button
                variant="primary"
                onClick={() => navigate('/create-artist')}
              >
                <Edit className="h-4 w-4 mr-2" />
                Редактировать
              </Button>
            </div>
          )}
        </div>
        {artist.description && (
          <Card className="mb-8">
            <p className="text-gray-300 text-lg leading-relaxed">{artist.description}</p>
          </Card>
        )}
      </div>

      {artist.albums && artist.albums.length > 0 && (
        <div>
          <div className="flex items-center gap-2 mb-6">
            <Disc className="h-8 w-8 text-primary-400" />
            <h2 className="text-3xl font-semibold text-white">Альбомы</h2>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {artist.albums.map((album) => (
              <Link key={album.id} to={`/album/${album.id}`}>
                <Card className="card-hover">
                  {album.cover_art_url && (
                    <img
                      src={album.cover_art_url}
                      alt={album.title}
                      className="w-full h-48 object-cover rounded-lg mb-4"
                    />
                  )}
                  <h3 className="text-xl font-semibold mb-2 text-white">{album.title}</h3>
                  <p className="text-gray-400 text-sm">
                    {formatDate(album.release_date)}
                  </p>
                  {album.songs && album.songs.length > 0 && (
                    <p className="text-gray-500 text-sm mt-1">
                      {album.songs.length} {album.songs.length === 1 ? 'песня' : 'песен'}
                    </p>
                  )}
                </Card>
              </Link>
            ))}
          </div>
        </div>
      )}

      {(!artist.albums || artist.albums.length === 0) && (
        <Card className="text-center py-12">
          <Disc className="h-16 w-16 text-gray-600 mx-auto mb-6" />
          <p className="text-gray-400 text-lg mb-6">У этого артиста пока нет альбомов</p>
          {canEdit && (
            <Button
              variant="primary"
              onClick={() => navigate('/create-album')}
            >
              <Plus className="h-4 w-4 mr-2" />
              Добавить альбом
            </Button>
          )}
        </Card>
      )}
    </div>
  );
}

