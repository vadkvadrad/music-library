import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { useMutation, useQuery } from '@tanstack/react-query';
import { musicApi } from '../api/music';
import { profileApi } from '../api/profile';
import { NewArtistRequest, UpdateArtistRequest } from '../types';
import Card from '../components/Card';
import Button from '../components/Button';
import Input from '../components/Input';
import { Mic, Edit } from 'lucide-react';

export default function CreateArtistPage() {
  const navigate = useNavigate();
  const [error, setError] = useState<string>('');

  const { data: profile } = useQuery({
    queryKey: ['profile'],
    queryFn: profileApi.getProfile,
  });

  const hasArtist = !!profile?.artist;

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<NewArtistRequest | UpdateArtistRequest>({
    defaultValues: hasArtist && profile?.artist ? {
      artist_name: profile.artist.name,
      description: profile.artist.description,
      formation_year: profile.artist.formation_year.split('T')[0],
    } : undefined,
  });

  useEffect(() => {
    if (hasArtist && profile.artist) {
      reset({
        artist_name: profile.artist.name,
        description: profile.artist.description,
        formation_year: profile.artist.formation_year.split('T')[0],
      });
    }
  }, [hasArtist, profile, reset]);

  const createMutation = useMutation({
    mutationFn: musicApi.createArtist,
    onSuccess: (data) => {
      navigate(`/artist/${data.id}`);
    },
    onError: (err: any) => {
      setError(err.response?.data?.error || 'Ошибка при создании артиста');
    },
  });

  const updateMutation = useMutation({
    mutationFn: (data: UpdateArtistRequest) =>
      musicApi.updateArtist(profile!.artist!.id, data),
    onSuccess: () => {
      navigate(`/artist/${profile!.artist!.id}`);
    },
    onError: (err: any) => {
      setError(err.response?.data?.error || 'Ошибка при обновлении артиста');
    },
  });

  const mutation = hasArtist ? updateMutation : createMutation;

  const onSubmit = (data: NewArtistRequest | UpdateArtistRequest) => {
    setError('');
    mutation.mutate(data as any);
  };

  return (
    <div>
      <div className="flex items-center gap-4 mb-6">
        {hasArtist ? (
          <Edit className="h-8 w-8 text-primary-600" />
        ) : (
          <Mic className="h-8 w-8 text-primary-600" />
        )}
        <h1 className="text-3xl font-bold">
          {hasArtist ? 'Редактировать артиста' : 'Создать артиста'}
        </h1>
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
            {hasArtist ? 'Сохранить изменения' : 'Создать артиста'}
          </Button>
        </form>
      </Card>
    </div>
  );
}

