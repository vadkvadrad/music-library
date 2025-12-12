import { useState, useEffect } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { profileApi } from '../api/profile';
import { useForm } from 'react-hook-form';
import { NewProfileRequest, UpdateProfileRequest } from '../types';
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

  const hasProfile = !!profile;

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<NewProfileRequest | UpdateProfileRequest>({
    defaultValues: profile,
  });

  useEffect(() => {
    if (profile) {
      reset({
        bio: profile.bio,
        avatar_url: profile.avatar_url,
      });
    }
  }, [profile, reset]);

  const createMutation = useMutation({
    mutationFn: profileApi.createProfile,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['profile'] });
      setIsEditing(false);
    },
  });

  const updateMutation = useMutation({
    mutationFn: profileApi.updateProfile,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['profile'] });
      setIsEditing(false);
    },
  });

  const mutation = hasProfile ? updateMutation : createMutation;

  const onSubmit = (data: NewProfileRequest | UpdateProfileRequest) => {
    if (hasProfile) {
      updateMutation.mutate(data);
    } else {
      createMutation.mutate(data as NewProfileRequest);
    }
  };

  if (isLoading) {
    return (
      <div className="text-center py-12">
        <div className="inline-block w-8 h-8 border-4 border-primary-600/30 border-t-primary-600 rounded-full animate-spin" />
      </div>
    );
  }

  if (!profile && !isEditing) {
    return (
      <div className="animate-fade-in">
        <div className="flex items-center gap-3 mb-6">
          <div className="p-2 rounded-xl bg-primary-600/20">
            <User className="h-6 w-6 text-primary-400" />
          </div>
          <h1 className="text-3xl font-bold text-white">Профиль</h1>
        </div>
        <Card>
          <p className="text-gray-400 mb-6">У вас еще нет профиля</p>
          <Button onClick={() => setIsEditing(true)} variant="primary">
            Создать профиль
          </Button>
        </Card>
      </div>
    );
  }

  return (
    <div className="animate-fade-in">
      <div className="flex items-center gap-3 mb-8">
        <div className="p-2 rounded-xl bg-primary-600/20">
          <User className="h-6 w-6 text-primary-400" />
        </div>
        <h1 className="text-3xl font-bold text-white">Профиль</h1>
      </div>

      {isEditing ? (
        <Card>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">
                Биография
              </label>
              <textarea
                {...register('bio', {
                  required: !hasProfile ? 'Биография обязательна' : false,
                })}
                className="input min-h-[100px]"
                rows={4}
              />
              {errors.bio && (
                <p className="mt-1 text-sm text-red-400">{errors.bio.message}</p>
              )}
            </div>

            <Input
              label="URL аватара"
              type="url"
              {...register('avatar_url', {
                required: !hasProfile ? 'URL аватара обязателен' : false,
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
              className="w-32 h-32 rounded-full object-cover mb-6 border-4 border-primary-500/30 shadow-medium"
            />
          )}
          <div className="mb-6">
            <h2 className="text-xl font-semibold mb-2 text-white">Имя пользователя</h2>
            <p className="text-gray-300 font-medium">{profile?.user_name || 'Не указано'}</p>
          </div>
          <div className="mb-6">
            <h2 className="text-xl font-semibold mb-2 text-white">Email</h2>
            <p className="text-gray-300">{profile?.user_email || 'Не указано'}</p>
          </div>
          {profile?.artist && (
            <div className="mb-6">
              <h2 className="text-xl font-semibold mb-2 text-white">Артист</h2>
              <p className="text-gray-300 font-medium">{profile.artist.name}</p>
            </div>
          )}
          <div className="mb-6">
            <h2 className="text-xl font-semibold mb-2 text-white">Биография</h2>
            <p className="text-gray-300 leading-relaxed">{profile?.bio}</p>
          </div>
          {profile?.created_at && (
            <p className="text-gray-500 text-sm mb-6">
              Создан: {formatDate(profile.created_at)}
            </p>
          )}
          
          {profile?.artist && (
            <div className="mt-6 p-6 bg-primary-600/10 rounded-xl border border-primary-500/30">
              <h3 className="text-lg font-semibold mb-3 text-white">Ваш артист</h3>
              <p className="text-gray-300 mb-2 font-medium">
                {profile.artist.name}
              </p>
              {profile.artist.description && (
                <p className="text-gray-400 text-sm mb-4">
                  {profile.artist.description}
                </p>
              )}
              <a
                href={`/artist/${profile.artist.id}`}
                className="text-primary-400 hover:text-primary-300 transition-colors text-sm font-medium"
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
            className="mt-6"
          >
            Редактировать
          </Button>
        </Card>
      )}
    </div>
  );
}

