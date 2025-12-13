import { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { musicApi } from '../api/music';
import { NewAlbumRequest, UpdateAlbumRequest } from '../types';
import Card from '../components/Card';
import Button from '../components/Button';
import Input from '../components/Input';
import { Disc } from 'lucide-react';

export default function CreateAlbumPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const queryClient = useQueryClient();
  const [error, setError] = useState<string>('');
  
  const albumId = location.state?.albumId;
  const isEditing = !!albumId;

  const { data: album } = useQuery({
    queryKey: ['album', albumId],
    queryFn: () => musicApi.getAlbum(albumId!.toString()),
    enabled: isEditing && !!albumId,
  });

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<NewAlbumRequest | UpdateAlbumRequest>({
    defaultValues: album,
  });

  useEffect(() => {
    if (album) {
      reset({
        title: album.title,
        release_date: album.release_date.split('T')[0], // Форматируем дату для input type="date"
        cover_art_url: album.cover_art_url,
      });
    }
  }, [album, reset]);

  const createMutation = useMutation({
    mutationFn: musicApi.createAlbum,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['profile'] });
      navigate('/');
    },
    onError: (err: any) => {
      setError(err.response?.data?.error || 'Ошибка при создании альбома');
    },
  });

  const updateMutation = useMutation({
    mutationFn: (data: UpdateAlbumRequest) => musicApi.updateAlbum(albumId!, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['album', albumId] });
      queryClient.invalidateQueries({ queryKey: ['profile'] });
      navigate(`/album/${albumId}`);
    },
    onError: (err: any) => {
      setError(err.response?.data?.error || 'Ошибка при обновлении альбома');
    },
  });

  const onSubmit = (data: NewAlbumRequest | UpdateAlbumRequest) => {
    setError('');
    if (isEditing) {
      updateMutation.mutate(data);
    } else {
      createMutation.mutate(data as NewAlbumRequest);
    }
  };

  const isLoading = createMutation.isPending || updateMutation.isPending;

  return (
    <div className="animate-fade-in">
      <div className="flex items-center gap-3 mb-8">
        <div className="p-2 rounded-xl bg-primary-600/20">
          <Disc className="h-6 w-6 text-primary-400" />
        </div>
        <h1 className="text-3xl font-bold text-white">{isEditing ? 'Редактировать альбом' : 'Создать альбом'}</h1>
      </div>

      <Card>
        {error && (
          <div className="mb-6 p-4 bg-red-500/20 border border-red-500/30 rounded-xl text-red-400">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <Input
            label="Название альбома"
            {...register('title', {
              required: !isEditing ? 'Название альбома обязательно' : false,
            })}
            error={errors.title?.message}
            placeholder="Название альбома"
          />

          <Input
            label="Дата выпуска"
            type="date"
            {...register('release_date', {
              required: !isEditing ? 'Дата выпуска обязательна' : false,
            })}
            error={errors.release_date?.message}
          />

          <Input
            label="URL обложки"
            type="url"
            {...register('cover_art_url', {
              required: !isEditing ? 'URL обложки обязателен' : false,
            })}
            error={errors.cover_art_url?.message}
            placeholder="https://example.com/cover.jpg"
          />

          <Button
            type="submit"
            variant="primary"
            isLoading={isLoading}
            className="w-full"
          >
            {isEditing ? 'Сохранить изменения' : 'Создать альбом'}
          </Button>
        </form>
      </Card>
    </div>
  );
}

