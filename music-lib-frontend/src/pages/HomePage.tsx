import { Link } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { Music, Search, Plus, User } from 'lucide-react';
import Button from '../components/Button';

export default function HomePage() {
  const { isAuthenticated } = useAuth();

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

