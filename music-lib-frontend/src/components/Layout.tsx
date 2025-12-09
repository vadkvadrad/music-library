import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { useQuery } from '@tanstack/react-query';
import { profileApi } from '../api/profile';
import { Search, User, LogOut, Music, Plus } from 'lucide-react';

interface LayoutProps {
  children: React.ReactNode;
}

export default function Layout({ children }: LayoutProps) {
  const { isAuthenticated, setToken } = useAuth();
  const navigate = useNavigate();

  const { data: profile } = useQuery({
    queryKey: ['profile'],
    queryFn: profileApi.getProfile,
    enabled: isAuthenticated,
  });

  const hasArtist = !!profile?.artist;

  const handleLogout = () => {
    setToken(null);
    navigate('/login');
  };

  return (
    <div className="min-h-screen bg-gradient-dark">
      <nav className="bg-dark-light/80 backdrop-blur-xl border-b border-gray-800/50 sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center space-x-6">
              <Link to="/" className="flex items-center space-x-2 group">
                <div className="p-2 rounded-xl bg-gradient-primary group-hover:scale-110 transition-transform duration-300">
                  <Music className="h-5 w-5 text-white" />
                </div>
                <span className="text-xl font-bold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">
                  Music Library
                </span>
              </Link>
              <Link
                to="/search"
                className="nav-link"
              >
                <Search className="h-5 w-5" />
                <span>Поиск</span>
              </Link>
            </div>
            <div className="flex items-center space-x-3">
              {isAuthenticated ? (
                <>
                  {!hasArtist && (
                    <Link
                      to="/create-artist"
                      className="nav-link"
                    >
                      <Plus className="h-5 w-5" />
                      <span>Создать артиста</span>
                    </Link>
                  )}
                  <Link
                    to="/profile"
                    className="nav-link"
                  >
                    <User className="h-5 w-5" />
                    <span>Профиль</span>
                  </Link>
                  <button
                    onClick={handleLogout}
                    className="nav-link text-red-400 hover:text-red-300 hover:bg-red-500/10"
                  >
                    <LogOut className="h-5 w-5" />
                    <span>Выйти</span>
                  </button>
                </>
              ) : (
                <>
                  <Link
                    to="/login"
                    className="btn btn-secondary"
                  >
                    Войти
                  </Link>
                  <Link
                    to="/register"
                    className="btn btn-primary"
                  >
                    Регистрация
                  </Link>
                </>
              )}
            </div>
          </div>
        </div>
      </nav>
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {children}
      </main>
    </div>
  );
}

