import { Link } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { useQuery } from '@tanstack/react-query';
import { profileApi } from '../api/profile';
import { Music, Search, Plus, User, Edit, Mic, Disc } from 'lucide-react';
import Button from '../components/Button';
import Card from '../components/Card';
import { formatDate } from '../utils/format';
import { Album } from '../types';

export default function HomePage() {
  const { isAuthenticated } = useAuth();
  
  const { data: profile } = useQuery({
    queryKey: ['profile'],
    queryFn: profileApi.getProfile,
    enabled: isAuthenticated,
  });
  
  const hasArtist = !!profile?.artist;

  return (
    <div className="text-center">
      <div className="mb-12">
        <div className="flex justify-center mb-6">
          <Music className="h-20 w-20 text-primary-600" />
        </div>
        <h1 className="text-5xl font-bold text-gray-900 mb-4">
          Добро пожаловать в Music Library
        </h1>
        <p className="text-xl text-gray-600 max-w-2xl mx-auto">
          Откройте для себя огромную коллекцию музыки, артистов и альбомов
        </p>
      </div>

      <div className={`grid grid-cols-1 ${isAuthenticated ? (hasArtist ? 'md:grid-cols-2' : 'md:grid-cols-3') : 'md:grid-cols-1'} gap-6 mb-12`}>
        <div className="card text-center">
          <Search className="h-12 w-12 text-primary-600 mx-auto mb-4" />
          <h2 className="text-xl font-semibold mb-2">Поиск</h2>
          <p className="text-gray-600 mb-4">
            Найдите любимых артистов, альбомы и песни
          </p>
          <Link to="/search">
            <Button variant="primary" className="w-full">
              Начать поиск
            </Button>
          </Link>
        </div>

        {isAuthenticated && (
          <>
            {hasArtist && profile?.artist?.id ? (
              <div className="card text-center">
                <Edit className="h-12 w-12 text-primary-600 mx-auto mb-4" />
                <h2 className="text-xl font-semibold mb-2">Редактировать</h2>
                <p className="text-gray-600 mb-4">
                  Управляйте своим артистом, альбомами и песнями
                </p>
                <Link to={`/artist/${profile.artist.id}`}>
                  <Button variant="primary" className="w-full">
                    Редактировать артиста
                  </Button>
                </Link>
              </div>
            ) : (
              <div className="card text-center">
                <Plus className="h-12 w-12 text-primary-600 mx-auto mb-4" />
                <h2 className="text-xl font-semibold mb-2">Создать</h2>
                <p className="text-gray-600 mb-4">
                  Добавьте своего артиста, альбом или песню
                </p>
                <Link to="/create-artist">
                  <Button variant="primary" className="w-full">
                    Создать артиста
                  </Button>
                </Link>
              </div>
            )}

            <div className="card text-center">
              <User className="h-12 w-12 text-primary-600 mx-auto mb-4" />
              <h2 className="text-xl font-semibold mb-2">Профиль</h2>
              <p className="text-gray-600 mb-4">
                Управляйте своим профилем и настройками
              </p>
              <Link to="/profile">
                <Button variant="secondary" className="w-full">
                  Мой профиль
                </Button>
              </Link>
            </div>
          </>
        )}
      </div>

      {isAuthenticated && hasArtist && profile?.artist && (
        <div className="mt-12">
          <div className="flex items-center gap-2 mb-6">
            <Mic className="h-8 w-8 text-primary-600" />
            <h2 className="text-3xl font-bold">Мой артист</h2>
          </div>
          
          <Card className="mb-8">
            <div className="flex items-start justify-between mb-4">
              <div className="flex items-center gap-4">
                <Mic className="h-12 w-12 text-primary-600" />
                <div>
                  <h3 className="text-2xl font-semibold">{profile.artist.name}</h3>
                  {profile.artist.formation_year && (
                    <p className="text-gray-600 mt-1">
                      Основан: {formatDate(profile.artist.formation_year)}
                    </p>
                  )}
                </div>
              </div>
              <Link to={`/artist/${profile.artist.id}`}>
                <Button variant="primary">
                  <Edit className="h-4 w-4 mr-2" />
                  Редактировать
                </Button>
              </Link>
            </div>
            {profile.artist.description && (
              <p className="text-gray-700 text-lg">{profile.artist.description}</p>
            )}
          </Card>

          {profile.artist.albums && profile.artist.albums.length > 0 && (
            <div>
              <div className="flex items-center gap-2 mb-4">
                <Disc className="h-6 w-6 text-primary-600" />
                <h3 className="text-2xl font-semibold">Мои альбомы</h3>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {profile.artist.albums.map((album: Album) => (
                  <Link key={album.id} to={`/album/${album.id}`}>
                    <Card>
                      {album.cover_art_url && (
                        <img
                          src={album.cover_art_url}
                          alt={album.title || 'Альбом'}
                          className="w-full h-48 object-cover rounded-lg mb-4"
                        />
                      )}
                      <h4 className="text-xl font-semibold mb-2">{album.title}</h4>
                      <p className="text-gray-500 text-sm">
                        {formatDate(album.release_date)}
                      </p>
                    </Card>
                  </Link>
                ))}
              </div>
            </div>
          )}

          {(!profile.artist.albums || profile.artist.albums.length === 0) && (
            <Card className="text-center py-8">
              <Disc className="h-12 w-12 text-gray-400 mx-auto mb-4" />
              <p className="text-gray-500 mb-4">У вас пока нет альбомов</p>
              <Link to="/create-album">
                <Button variant="primary">Создать альбом</Button>
              </Link>
            </Card>
          )}
        </div>
      )}

      {!isAuthenticated && (
        <div className="bg-primary-50 rounded-xl p-8">
          <h2 className="text-2xl font-semibold mb-4">
            Присоединяйтесь к сообществу
          </h2>
          <p className="text-gray-700 mb-6">
            Зарегистрируйтесь, чтобы создавать и делиться музыкой
          </p>
          <div className="flex justify-center space-x-4">
            <Link to="/register">
              <Button variant="primary">Регистрация</Button>
            </Link>
            <Link to="/login">
              <Button variant="secondary">Войти</Button>
            </Link>
          </div>
        </div>
      )}
    </div>
  );
}

