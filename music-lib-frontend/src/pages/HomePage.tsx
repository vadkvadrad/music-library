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
    <div className="text-center animate-fade-in">
      <div className="mb-16">
        <div className="flex justify-center mb-6">
          <div className="p-4 rounded-2xl bg-gradient-primary shadow-glow animate-scale-in">
            <Music className="h-16 w-16 text-white" />
          </div>
        </div>
        <h1 className="text-5xl md:text-6xl font-bold mb-4 bg-gradient-to-r from-white via-gray-100 to-gray-300 bg-clip-text text-transparent animate-slide-up">
          Добро пожаловать в Music Library
        </h1>
        <p className="text-xl text-gray-400 max-w-2xl mx-auto animate-slide-up">
          Откройте для себя огромную коллекцию музыки, артистов и альбомов
        </p>
      </div>

      <div className={`grid grid-cols-1 ${isAuthenticated ? (hasArtist ? 'md:grid-cols-2' : 'md:grid-cols-3') : 'md:grid-cols-1'} gap-6 mb-12`}>
        <Card hover className="text-center">
          <div className="p-3 rounded-xl bg-primary-600/20 w-fit mx-auto mb-4">
            <Search className="h-8 w-8 text-primary-400" />
          </div>
          <h2 className="text-xl font-semibold mb-2 text-white">Поиск</h2>
          <p className="text-gray-400 mb-6">
            Найдите любимых артистов, альбомы и песни
          </p>
          <Link to="/search">
            <Button variant="primary" className="w-full">
              Начать поиск
            </Button>
          </Link>
        </Card>

        {isAuthenticated && (
          <>
            {hasArtist && profile?.artist?.id ? (
              <Card hover className="text-center">
                <div className="p-3 rounded-xl bg-primary-600/20 w-fit mx-auto mb-4">
                  <Edit className="h-8 w-8 text-primary-400" />
                </div>
                <h2 className="text-xl font-semibold mb-2 text-white">Редактировать</h2>
                <p className="text-gray-400 mb-6">
                  Управляйте своим артистом, альбомами и песнями
                </p>
                <Link to={`/artist/${profile.artist.id}`}>
                  <Button variant="primary" className="w-full">
                    Редактировать артиста
                  </Button>
                </Link>
              </Card>
            ) : (
              <Card hover className="text-center">
                <div className="p-3 rounded-xl bg-primary-600/20 w-fit mx-auto mb-4">
                  <Plus className="h-8 w-8 text-primary-400" />
                </div>
                <h2 className="text-xl font-semibold mb-2 text-white">Создать</h2>
                <p className="text-gray-400 mb-6">
                  Добавьте своего артиста, альбом или песню
                </p>
                <Link to="/create-artist">
                  <Button variant="primary" className="w-full">
                    Создать артиста
                  </Button>
                </Link>
              </Card>
            )}

            <Card hover className="text-center">
              <div className="p-3 rounded-xl bg-primary-600/20 w-fit mx-auto mb-4">
                <User className="h-8 w-8 text-primary-400" />
              </div>
              <h2 className="text-xl font-semibold mb-2 text-white">Профиль</h2>
              <p className="text-gray-400 mb-6">
                Управляйте своим профилем и настройками
              </p>
              <Link to="/profile">
                <Button variant="secondary" className="w-full">
                  Мой профиль
                </Button>
              </Link>
            </Card>
          </>
        )}
      </div>

      {isAuthenticated && hasArtist && profile?.artist && (
        <div className="mt-16 animate-slide-up">
          <div className="flex items-center gap-3 mb-8">
            <div className="p-2 rounded-xl bg-primary-600/20">
              <Mic className="h-6 w-6 text-primary-400" />
            </div>
            <h2 className="text-3xl font-bold text-white">Мой артист</h2>
          </div>
          
          <Card className="mb-8">
            <div className="flex items-start justify-between mb-4">
              <div className="flex items-center gap-4">
                <div className="p-3 rounded-xl bg-gradient-primary">
                  <Mic className="h-10 w-10 text-white" />
                </div>
                <div>
                  <h3 className="text-2xl font-semibold text-white mb-1">{profile.artist.name}</h3>
                  {profile.artist.formation_year && (
                    <p className="text-gray-400 text-sm">
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
              <p className="text-gray-300 text-lg leading-relaxed">{profile.artist.description}</p>
            )}
          </Card>

          {profile.artist.albums && profile.artist.albums.length > 0 && (
            <div>
              <div className="flex items-center gap-3 mb-6">
                <Disc className="h-6 w-6 text-primary-400" />
                <h3 className="text-2xl font-semibold text-white">Мои альбомы</h3>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {profile.artist.albums.map((album: Album) => (
                  <Link key={album.id} to={`/album/${album.id}`}>
                    <Card hover>
                      {album.cover_art_url && (
                        <img
                          src={album.cover_art_url}
                          alt={album.title || 'Альбом'}
                          className="w-full h-48 object-cover rounded-xl mb-4 shadow-medium"
                        />
                      )}
                      <h4 className="text-xl font-semibold mb-2 text-white">{album.title}</h4>
                      <p className="text-gray-400 text-sm">
                        {formatDate(album.release_date)}
                      </p>
                    </Card>
                  </Link>
                ))}
              </div>
            </div>
          )}

          {(!profile.artist.albums || profile.artist.albums.length === 0) && (
            <Card className="text-center py-12">
              <div className="p-4 rounded-xl bg-gray-800/50 w-fit mx-auto mb-4">
                <Disc className="h-12 w-12 text-gray-500" />
              </div>
              <p className="text-gray-400 mb-6">У вас пока нет альбомов</p>
              <Link to="/create-album">
                <Button variant="primary">Создать альбом</Button>
              </Link>
            </Card>
          )}
        </div>
      )}

      {!isAuthenticated && (
        <Card className="text-center py-12 animate-slide-up">
          <h2 className="text-3xl font-bold mb-4 text-white">
            Присоединяйтесь к сообществу
          </h2>
          <p className="text-gray-400 mb-8 text-lg">
            Зарегистрируйтесь, чтобы создавать и делиться музыкой
          </p>
          <div className="flex justify-center gap-4">
            <Link to="/register">
              <Button variant="primary" size="lg">Регистрация</Button>
            </Link>
            <Link to="/login">
              <Button variant="secondary" size="lg">Войти</Button>
            </Link>
          </div>
        </Card>
      )}
    </div>
  );
}

