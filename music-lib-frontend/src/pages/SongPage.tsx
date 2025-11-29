import { useParams, Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { musicApi } from '../api/music';
import Card from '../components/Card';
import { formatDuration } from '../utils/format';
import { Music, Disc } from 'lucide-react';

export default function SongPage() {
  const { id } = useParams<{ id: string }>();

  const { data: song, isLoading, error } = useQuery({
    queryKey: ['song', id],
    queryFn: () => musicApi.getSong(id!),
    enabled: !!id,
  });

  if (isLoading) {
    return <div className="text-center py-12">Загрузка...</div>;
  }

  if (error || !song) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-red-700">
        Ошибка при загрузке песни
      </div>
    );
  }

  return (
    <div>
      <div className="mb-8">
        <div className="flex items-center gap-4 mb-4">
          <Music className="h-12 w-12 text-primary-600" />
          <div>
            <h1 className="text-4xl font-bold">{song.title}</h1>
            <p className="text-gray-600 mt-1">
              Длительность: {formatDuration(song.duration)}
            </p>
          </div>
        </div>
        {song.album_id && (
          <Link
            to={`/album/${song.album_id}`}
            className="text-primary-600 hover:underline flex items-center gap-2"
          >
            <Disc className="h-4 w-4" />
            Перейти к альбому
          </Link>
        )}
      </div>

      {song.lyrics && song.lyrics.text && song.lyrics.text.length > 0 && (
        <Card className="mb-8">
          <h2 className="text-2xl font-semibold mb-4">Текст песни</h2>
          <div className="space-y-4">
            {song.lyrics.text.map((couplet, index) => (
              <div key={index} className="text-gray-700 whitespace-pre-line">
                {couplet.couplet}
              </div>
            ))}
          </div>
        </Card>
      )}

      {song.file_path && (
        <Card>
          <h2 className="text-xl font-semibold mb-2">Файл</h2>
          <a
            href={song.file_path}
            target="_blank"
            rel="noopener noreferrer"
            className="text-primary-600 hover:underline"
          >
            {song.file_path}
          </a>
        </Card>
      )}
    </div>
  );
}

