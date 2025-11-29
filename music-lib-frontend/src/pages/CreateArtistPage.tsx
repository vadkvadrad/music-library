import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { useMutation } from '@tanstack/react-query';
import { musicApi } from '../api/music';
import { NewArtistRequest } from '../types';
import Card from '../components/Card';
import Button from '../components/Button';
import Input from '../components/Input';
import { Mic } from 'lucide-react';

export default function CreateArtistPage() {
  const navigate = useNavigate();
  const [error, setError] = useState<string>('');

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<NewArtistRequest>();

  const mutation = useMutation({
    mutationFn: musicApi.createArtist,
    onSuccess: (data) => {
      navigate(`/artist/${data.id}`);
    },
    onError: (err: any) => {
      setError(err.response?.data?.error || 'Ошибка при создании артиста');
    },
  });

  const onSubmit = (data: NewArtistRequest) => {
    setError('');
    mutation.mutate(data);
  };

  return (
    <div>
      <div className="flex items-center gap-4 mb-6">
        <Mic className="h-8 w-8 text-primary-600" />
        <h1 className="text-3xl font-bold">Создать артиста</h1>
      </div>

      <Card>
        {error && (
          <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg text-red-700">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <Input
            label="Название артиста"
            {...register('artist_name', {
              required: 'Название артиста обязательно',
            })}
            error={errors.artist_name?.message}
            placeholder="Imagine Dragons"
          />

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Описание
            </label>
            <textarea
              {...register('description')}
              className="input min-h-[100px]"
              placeholder="Описание артиста..."
            />
          </div>

          <Input
            label="Год основания"
            type="date"
            {...register('formation_year', {
              required: 'Год основания обязателен',
            })}
            error={errors.formation_year?.message}
          />

          <Button
            type="submit"
            variant="primary"
            isLoading={mutation.isPending}
            className="w-full"
          >
            Создать артиста
          </Button>
        </form>
      </Card>
    </div>
  );
}

