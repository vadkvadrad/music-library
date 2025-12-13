import { useState, useEffect } from 'react';
import { useParams, useNavigate, useLocation } from 'react-router-dom';
import { useForm, useFieldArray } from 'react-hook-form';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { musicApi } from '../api/music';
import { NewSongRequest, UpdateSongRequest } from '../types';
import Card from '../components/Card';
import Button from '../components/Button';
import Input from '../components/Input';
import { Music, Plus, Trash2 } from 'lucide-react';

export default function CreateSongPage() {
  const { albumId } = useParams<{ albumId: string }>();
  const navigate = useNavigate();
  const location = useLocation();
  const queryClient = useQueryClient();
  const [error, setError] = useState<string>('');

  const songId = location.state?.songId;
  const isEditing = !!songId;

  const { data: song } = useQuery({
    queryKey: ['song', songId],
    queryFn: () => musicApi.getSong(songId!.toString()),
    enabled: isEditing && !!songId,
  });

  const { data: allGenres } = useQuery({
    queryKey: ['genres'],
    queryFn: musicApi.getAllGenres,
  });

  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
    reset,
  } = useForm<NewSongRequest | UpdateSongRequest>({
    defaultValues: {
      genres: [{ genre_id: 0 }],
      lyrics: {
        text: [{ couplet: '' }],
      },
    },
  });

  useEffect(() => {
    if (song) {
      reset({
        title: song.title,
        duration_sec: song.duration,
        file_path: song.file_path,
        genres: song.genres && song.genres.length > 0
          ? song.genres.map((g) => ({ genre_id: g.id }))
          : [{ genre_id: 0 }],
        lyrics: {
          text: song.lyrics?.text?.map((c) => ({ couplet: c.couplet })) || [{ couplet: '' }],
        },
      });
    }
  }, [song, reset]);

  const {
    fields: genreFields,
    append: appendGenre,
    remove: removeGenre,
  } = useFieldArray({
    control,
    name: 'genres',
  });

  const {
    fields: coupletFields,
    append: appendCouplet,
    remove: removeCouplet,
  } = useFieldArray({
    control,
    name: 'lyrics.text',
  });

  const createMutation = useMutation({
    mutationFn: (data: NewSongRequest) =>
      musicApi.createSong(Number(albumId), data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['profile'] });
      if (albumId) {
        navigate(`/album/${albumId}`);
      }
    },
    onError: (err: any) => {
      setError(err.response?.data?.error || 'Ошибка при создании песни');
    },
  });

  const updateMutation = useMutation({
    mutationFn: (data: UpdateSongRequest) => musicApi.updateSong(songId!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['song', songId] });
      queryClient.invalidateQueries({ queryKey: ['profile'] });
      navigate(`/song/${songId}`);
    },
    onError: (err: any) => {
      setError(err.response?.data?.error || 'Ошибка при обновлении песни');
    },
  });

  const onSubmit = (data: NewSongRequest | UpdateSongRequest) => {
    setError('');
    if (isEditing) {
      updateMutation.mutate(data);
    } else {
      createMutation.mutate(data as NewSongRequest);
    }
  };

  const isLoading = createMutation.isPending || updateMutation.isPending;

  return (
    <div className="animate-fade-in">
      <div className="flex items-center gap-4 mb-6">
        <div className="p-2 rounded-xl bg-primary-600/20">
          <Music className="h-8 w-8 text-primary-400" />
        </div>
        <h1 className="text-3xl font-bold text-white">{isEditing ? 'Редактировать песню' : 'Добавить песню'}</h1>
      </div>

      <Card>
        {error && (
          <div className="mb-6 p-4 bg-red-500/20 border border-red-500/30 rounded-xl text-red-400">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          <Input
            label="Название песни"
            {...register('title', {
              required: !isEditing ? 'Название песни обязательно' : false,
            })}
            error={errors.title?.message}
            placeholder="Название песни"
          />

          <Input
            label="Длительность (секунды)"
            type="number"
            {...register('duration_sec', {
              required: !isEditing ? 'Длительность обязательна' : false,
              min: { value: 1, message: 'Длительность должна быть больше 0' },
              valueAsNumber: true,
            })}
            error={errors.duration_sec?.message}
          />

          <Input
            label="Путь к файлу"
            {...register('file_path', {
              required: !isEditing ? 'Путь к файлу обязателен' : false,
            })}
            error={errors.file_path?.message}
            placeholder="/path/to/song.mp3"
          />

          <div>
            <label className="block text-sm font-medium text-white mb-2">
              Жанры
            </label>
            {genreFields.map((field, index) => (
              <div key={field.id} className="mb-2">
                <div className="flex gap-2">
                  <select
                    {...register(`genres.${index}.genre_id` as const, {
                      required: !isEditing ? 'Жанр обязателен' : false,
                      validate: (value) => {
                        if (!isEditing && value === 0) {
                          return 'Выберите жанр';
                        }
                        return true;
                      },
                      valueAsNumber: true,
                    })}
                    className="input flex-1"
                  >
                    <option value={0}>Выберите жанр</option>
                    {allGenres?.map((genre) => (
                      <option key={genre.id} value={genre.id}>
                        {genre.name}
                      </option>
                    ))}
                  </select>
                  {genreFields.length > 1 && (
                    <Button
                      type="button"
                      variant="danger"
                      onClick={() => removeGenre(index)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  )}
                </div>
                {errors.genres?.[index]?.genre_id && (
                  <p className="text-red-400 text-sm mt-1">
                    {errors.genres[index]?.genre_id?.message}
                  </p>
                )}
              </div>
            ))}
            <Button
              type="button"
              variant="secondary"
              onClick={() => appendGenre({ genre_id: 0 })}
            >
              <Plus className="h-4 w-4 mr-2" />
              Добавить жанр
            </Button>
          </div>

          <div>
            <label className="block text-sm font-medium text-white mb-2">
              Текст песни (куплеты)
            </label>
            {coupletFields.map((field, index) => (
              <div key={field.id} className="mb-4">
                <label className="block text-xs text-gray-500 mb-1">
                  Куплет {index + 1}
                </label>
                <div className="flex gap-2">
                  <textarea
                    {...register(`lyrics.text.${index}.couplet` as const, {
                      required: !isEditing ? 'Текст куплета обязателен' : false,
                    })}
                    className="input min-h-[80px] flex-1"
                    placeholder="Текст куплета..."
                  />
                  {coupletFields.length > 1 && (
                    <Button
                      type="button"
                      variant="danger"
                      onClick={() => removeCouplet(index)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  )}
                </div>
              </div>
            ))}
            <Button
              type="button"
              variant="secondary"
              onClick={() => appendCouplet({ couplet: '' })}
            >
              <Plus className="h-4 w-4 mr-2" />
              Добавить куплет
            </Button>
          </div>

          <Button
            type="submit"
            variant="primary"
            isLoading={isLoading}
            className="w-full"
          >
            {isEditing ? 'Сохранить изменения' : 'Создать песню'}
          </Button>
        </form>
      </Card>
    </div>
  );
}

