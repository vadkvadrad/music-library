import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { useMutation } from '@tanstack/react-query';
import { musicApi } from '../api/music';
import { NewAlbumRequest } from '../types';
import Card from '../components/Card';
import Button from '../components/Button';
import Input from '../components/Input';
import { Disc } from 'lucide-react';

export default function CreateAlbumPage() {
  const navigate = useNavigate();
  const [error, setError] = useState<string>('');

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<NewAlbumRequest>();

  const mutation = useMutation({
    mutationFn: musicApi.createAlbum,
    onSuccess: () => {
      navigate('/');
    },
    onError: (err: any) => {
      setError(err.response?.data?.error || 'Ошибка при создании альбома');
    },
  });

  const onSubmit = (data: NewAlbumRequest) => {
    setError('');
    mutation.mutate(data);
  };

  return (
    <div>
      <div className="flex items-center gap-4 mb-6">
        <Disc className="h-8 w-8 text-primary-600" />
        <h1 className="text-3xl font-bold">Создать альбом</h1>
      </div>

      <Card>
        {error && (
          <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg text-red-700">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <Input
            label="Название альбома"
            {...register('title', {
              required: 'Название альбома обязательно',
            })}
            error={errors.title?.message}
            placeholder="Название альбома"
          />

          <Input
            label="Дата выпуска"
            type="date"
            {...register('release_date', {
              required: 'Дата выпуска обязательна',
            })}
            error={errors.release_date?.message}
          />

          <Input
            label="URL обложки"
            type="url"
            {...register('cover_art_url', {
              required: 'URL обложки обязателен',
            })}
            error={errors.cover_art_url?.message}
            placeholder="https://example.com/cover.jpg"
          />

          <Button
            type="submit"
            variant="primary"
            isLoading={mutation.isPending}
            className="w-full"
          >
            Создать альбом
          </Button>
        </form>
      </Card>
    </div>
  );
}

