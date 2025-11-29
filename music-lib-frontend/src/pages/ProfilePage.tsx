import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { profileApi } from '../api/profile';
import { useForm } from 'react-hook-form';
import { NewProfileRequest } from '../types';
import Card from '../components/Card';
import Button from '../components/Button';
import Input from '../components/Input';
import { User, Save } from 'lucide-react';
import { formatDate } from '../utils/format';

export default function ProfilePage() {
  const queryClient = useQueryClient();
  const [isEditing, setIsEditing] = useState(false);

  const { data: profile, isLoading } = useQuery({
    queryKey: ['profile'],
    queryFn: profileApi.getProfile,
  });

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<NewProfileRequest>({
    defaultValues: profile,
  });

  const mutation = useMutation({
    mutationFn: profileApi.createProfile,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['profile'] });
      setIsEditing(false);
    },
  });

  const onSubmit = (data: NewProfileRequest) => {
    mutation.mutate(data);
  };

  if (isLoading) {
    return <div className="text-center py-12">Загрузка...</div>;
  }

  if (!profile && !isEditing) {
    return (
      <div>
        <h1 className="text-3xl font-bold mb-6">Профиль</h1>
        <Card>
          <p className="text-gray-600 mb-4">У вас еще нет профиля</p>
          <Button onClick={() => setIsEditing(true)} variant="primary">
            Создать профиль
          </Button>
        </Card>
      </div>
    );
  }

  return (
    <div>
      <div className="flex items-center gap-4 mb-6">
        <User className="h-8 w-8 text-primary-600" />
        <h1 className="text-3xl font-bold">Профиль</h1>
      </div>

      {isEditing ? (
        <Card>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Биография
              </label>
              <textarea
                {...register('bio', {
                  required: 'Биография обязательна',
                })}
                className="input min-h-[100px]"
                rows={4}
              />
              {errors.bio && (
                <p className="mt-1 text-sm text-red-600">{errors.bio.message}</p>
              )}
            </div>

            <Input
              label="URL аватара"
              type="url"
              {...register('avatar_url', {
                required: 'URL аватара обязателен',
              })}
              error={errors.avatar_url?.message}
            />

            <div className="flex gap-4">
              <Button
                type="submit"
                variant="primary"
                isLoading={mutation.isPending}
              >
                <Save className="h-4 w-4 mr-2" />
                Сохранить
              </Button>
              <Button
                type="button"
                variant="secondary"
                onClick={() => {
                  setIsEditing(false);
                  reset();
                }}
              >
                Отмена
              </Button>
            </div>
          </form>
        </Card>
      ) : (
        <Card>
          {profile?.avatar_url && (
            <img
              src={profile.avatar_url}
              alt="Avatar"
              className="w-32 h-32 rounded-full object-cover mb-4"
            />
          )}
          {profile?.artist && (
            <div className="mb-4">
              <h2 className="text-xl font-semibold mb-2">Артист</h2>
              <p className="text-gray-700 font-medium">{profile.artist.name}</p>
            </div>
          )}
          <div className="mb-4">
            <h2 className="text-xl font-semibold mb-2">Биография</h2>
            <p className="text-gray-700">{profile?.bio}</p>
          </div>
          {profile?.created_at && (
            <p className="text-gray-500 text-sm">
              Создан: {formatDate(profile.created_at)}
            </p>
          )}
          
          {profile?.artist && (
            <div className="mt-6 p-4 bg-primary-50 rounded-lg border border-primary-200">
              <h3 className="text-lg font-semibold mb-2">Ваш артист</h3>
              <p className="text-gray-700 mb-2">
                <strong>{profile.artist.name}</strong>
              </p>
              {profile.artist.description && (
                <p className="text-gray-600 text-sm mb-2">
                  {profile.artist.description}
                </p>
              )}
              <a
                href={`/artist/${profile.artist.id}`}
                className="text-primary-600 hover:underline text-sm"
              >
                Посмотреть профиль артиста →
              </a>
            </div>
          )}
          
          <Button
            onClick={() => {
              setIsEditing(true);
              reset(profile);
            }}
            variant="primary"
            className="mt-4"
          >
            Редактировать
          </Button>
        </Card>
      )}
    </div>
  );
}

