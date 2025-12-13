import { Routes, Route, Navigate } from 'react-router-dom';
import { useAuth } from './contexts/AuthContext';
import Layout from './components/Layout';
import HomePage from './pages/HomePage';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import VerifyPage from './pages/VerifyPage';
import SearchPage from './pages/SearchPage';
import ArtistPage from './pages/ArtistPage';
import AlbumPage from './pages/AlbumPage';
import SongPage from './pages/SongPage';
import ProfilePage from './pages/ProfilePage';
import CreateArtistPage from './pages/CreateArtistPage';
import CreateAlbumPage from './pages/CreateAlbumPage';
import CreateSongPage from './pages/CreateSongPage';
import GenresPage from './pages/GenresPage';

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuth();
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" />;
}

function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />
      <Route path="/verify" element={<VerifyPage />} />
      <Route
        path="/"
        element={
          <Layout>
            <HomePage />
          </Layout>
        }
      />
      <Route
        path="/search"
        element={
          <Layout>
            <SearchPage />
          </Layout>
        }
      />
      <Route
        path="/artist/:id"
        element={
          <Layout>
            <ArtistPage />
          </Layout>
        }
      />
      <Route
        path="/album/:id"
        element={
          <Layout>
            <AlbumPage />
          </Layout>
        }
      />
      <Route
        path="/song/:id"
        element={
          <Layout>
            <SongPage />
          </Layout>
        }
      />
      <Route
        path="/profile"
        element={
          <ProtectedRoute>
            <Layout>
              <ProfilePage />
            </Layout>
          </ProtectedRoute>
        }
      />
      <Route
        path="/create-artist"
        element={
          <ProtectedRoute>
            <Layout>
              <CreateArtistPage />
            </Layout>
          </ProtectedRoute>
        }
      />
      <Route
        path="/create-album"
        element={
          <ProtectedRoute>
            <Layout>
              <CreateAlbumPage />
            </Layout>
          </ProtectedRoute>
        }
      />
      <Route
        path="/create-song/:albumId"
        element={
          <ProtectedRoute>
            <Layout>
              <CreateSongPage />
            </Layout>
          </ProtectedRoute>
        }
      />
      <Route
        path="/genres"
        element={
          <ProtectedRoute>
            <Layout>
              <GenresPage />
            </Layout>
          </ProtectedRoute>
        }
      />
    </Routes>
  );
}

export default App;

