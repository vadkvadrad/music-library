import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useForm, useFieldArray } from 'react-hook-form';
import { useMutation } from '@tanstack/react-query';
import { musicApi } from '../api/music';
import { NewSongRequest } from '../types';
import Card from '../components/Card';
import Button from '../components/Button';
import Input from '../components/Input';
import { Music, Plus, Trash2 } from 'lucide-react';

export default function CreateSongPage() {
  const { albumId } = useParams<{ albumId: string }>();
  const navigate = useNavigate();
  const [error, setError] = useState<string>('');

  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
  } = useForm<NewSongRequest>({
    defaultValues: {
      genres: [{ genre_id: 0 }],
      lyrics: {
        text: [{ couplet: '' }],
      },
    },
  });

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

  const mutation = useMutation({
    mutationFn: (data: NewSongRequest) =>
      musicApi.createSong(Number(albumId), data),
    onSuccess: () => {
      navigate(`/album/${albumId}`);
    },
    onError: (err: any) => {
      setError(err.response?.data?.error || 'Ошибка при создании песни');
    },
  });

  const onSubmit = (data: NewSongRequest) => {
    setError('');
    mutation.mutate(data);
  };

  return (
    <div>
      <div className="flex items-center gap-4 mb-6">
        <Music className="h-8 w-8 text-primary-600" />
        <h1 className="text-3xl font-bold">Добавить песню</h1>
      </div>

      <Card>
        {error && (
          <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg text-red-700">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          <Input
            label="Название песни"
            {...register('title', {
              required: 'Название песни обязательно',
            })}
            error={errors.title?.message}
            placeholder="Название песни"
          />

          <Input
            label="Длительность (секунды)"
            type="number"
            {...register('duration_sec', {
              required: 'Длительность обязательна',
              min: { value: 1, message: 'Длительность должна быть больше 0' },
            })}
            error={errors.duration_sec?.message}
          />

          <Input
            label="Путь к файлу"
            {...register('file_path', {
              required: 'Путь к файлу обязателен',
            })}
            error={errors.file_path?.message}
            placeholder="/path/to/song.mp3"
          />

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Жанры
            </label>
            {genreFields.map((field, index) => (
              <div key={field.id} className="flex gap-2 mb-2">
                <Input
                  type="number"
                  {...register(`genres.${index}.genre_id` as const, {
                    required: 'ID жанра обязателен',
                    valueAsNumber: true,
                  })}
                  placeholder="ID жанра"
                />
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
            <label className="block text-sm font-medium text-gray-700 mb-2">
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
                      required: 'Текст куплета обязателен',
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
            isLoading={mutation.isPending}
            className="w-full"
          >
            Создать песню
          </Button>
        </form>
      </Card>
    </div>
  );
}

