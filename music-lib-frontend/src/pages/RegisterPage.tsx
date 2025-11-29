import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { useMutation } from '@tanstack/react-query';
import { authApi } from '../api/auth';
import { RegisterRequest } from '../types';
import Button from '../components/Button';
import Input from '../components/Input';
import { Music } from 'lucide-react';

export default function RegisterPage() {
  const navigate = useNavigate();
  const [error, setError] = useState<string>('');
  const [sessionId, setSessionId] = useState<string>('');

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterRequest>();

  const mutation = useMutation({
    mutationFn: authApi.register,
    onSuccess: (data) => {
      setSessionId(data.session_id);
      navigate('/verify', { state: { sessionId: data.session_id } });
    },
    onError: (err: any) => {
      setError(err.response?.data?.error || 'Ошибка регистрации');
    },
  });

  const onSubmit = (data: RegisterRequest) => {
    setError('');
    mutation.mutate(data);
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary-50 to-primary-100">
      <div className="max-w-md w-full">
        <div className="bg-white rounded-2xl shadow-xl p-8">
          <div className="text-center mb-8">
            <div className="flex justify-center mb-4">
              <Music className="h-12 w-12 text-primary-600" />
            </div>
            <h1 className="text-3xl font-bold text-gray-900">Регистрация</h1>
            <p className="text-gray-600 mt-2">Создайте новую учетную запись</p>
          </div>

          {error && (
            <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg text-red-700">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <Input
              label="Имя"
              {...register('name', {
                required: 'Имя обязательно',
                minLength: {
                  value: 2,
                  message: 'Имя должно содержать минимум 2 символа',
                },
              })}
              error={errors.name?.message}
            />

            <Input
              label="Email"
              type="email"
              {...register('email', {
                required: 'Email обязателен',
                pattern: {
                  value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
                  message: 'Неверный формат email',
                },
              })}
              error={errors.email?.message}
            />

            <Input
              label="Пароль"
              type="password"
              {...register('password', {
                required: 'Пароль обязателен',
                minLength: {
                  value: 6,
                  message: 'Пароль должен содержать минимум 6 символов',
                },
              })}
              error={errors.password?.message}
            />

            <Button
              type="submit"
              variant="primary"
              isLoading={mutation.isPending}
              className="w-full"
            >
              Зарегистрироваться
            </Button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-gray-600">
              Уже есть аккаунт?{' '}
              <Link to="/login" className="text-primary-600 hover:underline">
                Войти
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

