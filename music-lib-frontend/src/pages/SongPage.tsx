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
    return (
      <div className="text-center py-12">
        <div className="inline-block w-8 h-8 border-4 border-primary-600/30 border-t-primary-600 rounded-full animate-spin" />
      </div>
    );
  }

  if (error || !song) {
    return (
      <Card className="bg-red-500/20 border border-red-500/30">
        <p className="text-red-400">Ошибка при загрузке песни</p>
      </Card>
    );
  }

  return (
    <div className="animate-fade-in">
      <div className="mb-8">
        <div className="flex items-center gap-4 mb-6">
          <div className="p-3 rounded-xl bg-primary-600/20">
            <Music className="h-10 w-10 text-primary-400" />
          </div>
          <div>
            <h1 className="text-4xl font-bold text-white">{song.title}</h1>
            <p className="text-gray-400 mt-2">
              Длительность: {formatDuration(song.duration)}
            </p>
          </div>
        </div>
        {song.album_id && (
          <Link
            to={`/album/${song.album_id}`}
            className="inline-flex items-center gap-2 text-primary-400 hover:text-primary-300 transition-colors"
          >
            <Disc className="h-4 w-4" />
            Перейти к альбому
          </Link>
        )}
      </div>

      {song.lyrics && song.lyrics.text && song.lyrics.text.length > 0 && (
        <Card className="mb-8">
          <div className="flex items-center gap-3 mb-6">
            <div className="p-2 rounded-lg bg-primary-600/20">
              <Music className="h-5 w-5 text-primary-400" />
            </div>
            <h2 className="text-2xl font-semibold text-white">Текст песни</h2>
          </div>
          <div className="space-y-6">
            {song.lyrics.text.map((couplet, index) => (
              <div 
                key={index} 
                className="text-gray-300 whitespace-pre-line leading-relaxed text-lg"
              >
                {couplet.couplet}
              </div>
            ))}
          </div>
        </Card>
      )}

      {(!song.lyrics || !song.lyrics.text || song.lyrics.text.length === 0) && (
        <Card className="mb-8 text-center py-12">
          <div className="p-4 rounded-xl bg-gray-800/50 w-fit mx-auto mb-4">
            <Music className="h-12 w-12 text-gray-500" />
          </div>
          <p className="text-gray-400">Текст песни отсутствует</p>
        </Card>
      )}

      {song.file_path && (
        <Card>
          <h2 className="text-xl font-semibold mb-4 text-white">Файл</h2>
          <a
            href={song.file_path}
            target="_blank"
            rel="noopener noreferrer"
            className="text-primary-400 hover:text-primary-300 transition-colors break-all"
          >
            {song.file_path}
          </a>
        </Card>
      )}
    </div>
  );
}

