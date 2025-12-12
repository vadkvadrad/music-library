import { useParams, Link, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { musicApi } from '../api/music';
import { profileApi } from '../api/profile';
import Card from '../components/Card';
import Button from '../components/Button';
import { formatDate } from '../utils/format';
import { Mic, Disc, Edit } from 'lucide-react';

export default function ArtistPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const { data: artist, isLoading, error } = useQuery({
    queryKey: ['artist', id],
    queryFn: () => musicApi.getArtist(id!),
    enabled: !!id,
  });

  const { data: profile } = useQuery({
    queryKey: ['profile'],
    queryFn: profileApi.getProfile,
  });

  const canEdit = profile?.artist && profile.artist.id === Number(id);

  if (isLoading) {
    return <div className="text-center py-12">Загрузка...</div>;
  }

  if (error || !artist) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-red-700">
        Ошибка при загрузке артиста
      </div>
    );
  }

  return (
    <div>
      <div className="mb-8">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-4">
            <Mic className="h-12 w-12 text-primary-600" />
            <div>
              <h1 className="text-4xl font-bold">{artist.name}</h1>
              {artist.formation_year && (
                <p className="text-gray-600 mt-1">
                  Основан: {formatDate(artist.formation_year)}
                </p>
              )}
            </div>
          </div>
          {canEdit && (
            <Button
              variant="primary"
              onClick={() => navigate('/create-artist')}
            >
              <Edit className="h-4 w-4 mr-2" />
              Редактировать
            </Button>
          )}
        </div>
        {artist.description && (
          <p className="text-gray-700 text-lg mt-4">{artist.description}</p>
        )}
      </div>

      {artist.albums && artist.albums.length > 0 && (
        <div>
          <h2 className="text-2xl font-semibold mb-4 flex items-center gap-2">
            <Disc className="h-6 w-6" />
            Альбомы
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {artist.albums.map((album) => (
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
        <div className="text-center py-12 text-gray-500">
          У этого артиста пока нет альбомов
        </div>
      )}
    </div>
  );
}

