import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { musicApi } from '../api/music';
import { NewGenreRequest } from '../types';
import Card from '../components/Card';
import Button from '../components/Button';
import Input from '../components/Input';
import { useForm } from 'react-hook-form';
import { Music, Plus } from 'lucide-react';

export default function GenresPage() {
  const queryClient = useQueryClient();
  const [error, setError] = useState<string>('');

  const { data: genres, isLoading } = useQuery({
    queryKey: ['genres'],
    queryFn: musicApi.getAllGenres,
  });

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<NewGenreRequest>();

  const mutation = useMutation({
    mutationFn: musicApi.createGenre,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['genres'] });
      reset();
      setError('');
    },
    onError: (err: any) => {
      setError(err.response?.data?.error || 'Ошибка при создании жанра');
    },
  });

  const onSubmit = (data: NewGenreRequest) => {
    setError('');
    mutation.mutate(data);
  };

  if (isLoading) {
    return (
      <div className="text-center py-12">
        <div className="inline-block w-8 h-8 border-4 border-primary-600/30 border-t-primary-600 rounded-full animate-spin" />
      </div>
    );
  }

  return (
    <div className="animate-fade-in">
      <div className="flex items-center gap-3 mb-8">
        <div className="p-2 rounded-xl bg-primary-600/20">
          <Music className="h-8 w-8 text-primary-400" />
        </div>
        <h1 className="text-3xl font-bold text-white">Жанры</h1>
      </div>

      <Card className="mb-8">
        <h2 className="text-2xl font-semibold mb-4 text-white">Создать новый жанр</h2>
        {error && (
          <div className="mb-4 p-4 bg-red-500/20 border border-red-500/30 rounded-xl text-red-400">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <Input
            label="Название жанра"
            {...register('genre_name', {
              required: 'Название жанра обязательно',
            })}
            error={errors.genre_name?.message}
            placeholder="Например: Rock, Pop, Jazz"
          />

          <Button
            type="submit"
            variant="primary"
            isLoading={mutation.isPending}
          >
            <Plus className="h-4 w-4 mr-2" />
            Создать жанр
          </Button>
        </form>
      </Card>

      <div>
        <h2 className="text-2xl font-semibold mb-4 text-white">Все жанры</h2>
        {genres && genres.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {genres.map((genre) => (
              <Card key={genre.id} className="card-hover">
                <div className="flex items-center gap-3">
                  <div className="p-2 rounded-lg bg-primary-600/20">
                    <Music className="h-5 w-5 text-primary-400" />
                  </div>
                  <div>
                    <h3 className="text-lg font-semibold text-white">{genre.name}</h3>
                    <p className="text-sm text-gray-400">ID: {genre.id}</p>
                  </div>
                </div>
              </Card>
            ))}
          </div>
        ) : (
          <Card className="text-center py-12">
            <Music className="h-16 w-16 text-gray-600 mx-auto mb-4" />
            <p className="text-gray-400 text-lg">Пока нет жанров</p>
          </Card>
        )}
      </div>
    </div>
  );
}

