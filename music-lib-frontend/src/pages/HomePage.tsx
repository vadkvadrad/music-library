import { Link } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { useQuery } from '@tanstack/react-query';
import { profileApi } from '../api/profile';
import { Music, Search, Plus, User, Edit } from 'lucide-react';
import Button from '../components/Button';
import Card from '../components/Card';
import { formatDate } from '../utils/format';
import { Key, ReactElement, JSXElementConstructor, ReactNode } from 'react';

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

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
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
            <div className="card text-center">
              {hasArtist && profile?.artist?.id ? (
                <>
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
                </>
              ) : (
                <>
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
                </>
              )}
            </div>

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
          <h2 className="text-3xl font-bold mb-6">Мой артист</h2>
          <div className="mb-8">
            <div className="flex items-center gap-4 mb-4">
              <h3 className="text-2xl font-semibold">{profile.artist.name}</h3>
              <Link to={`/artist/${profile.artist.id}`}>
                <Button variant="primary">
                  <Edit className="h-4 w-4 mr-2" />
                  Редактировать
                </Button>
              </Link>
            </div>
            {profile.artist.description && (
              <p className="text-gray-700 mb-4">{profile.artist.description}</p>
            )}
          </div>

          {profile.artist.albums && profile.artist.albums.length > 0 && (
            <div>
              <h3 className="text-2xl font-semibold mb-4">Мои альбомы</h3>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {profile.artist.albums.map((album: { id: Key | null | undefined; cover_art_url: string | undefined; title: string | number | boolean | ReactElement<any, string | JSXElementConstructor<any>> | Iterable<ReactNode> | null | undefined; release_date: string; }) => (
                  <Link key={album.id} to={`/album/${album.id}`}>
                    <Card>
                      {album.cover_art_url && (
                        <img
                          src={album.cover_art_url}
                          alt={album.title}
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
            <div className="text-center py-8 text-gray-500">
              <p className="mb-4">У вас пока нет альбомов</p>
              <Link to="/create-album">
                <Button variant="primary">Создать альбом</Button>
              </Link>
            </div>
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

