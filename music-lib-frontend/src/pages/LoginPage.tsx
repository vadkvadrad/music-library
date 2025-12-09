import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { useMutation } from '@tanstack/react-query';
import { authApi } from '../api/auth';
import { useAuth } from '../contexts/AuthContext';
import { LoginRequest } from '../types';
import Button from '../components/Button';
import Input from '../components/Input';
import Card from '../components/Card';
import { Music } from 'lucide-react';

export default function LoginPage() {
  const navigate = useNavigate();
  const { setToken } = useAuth();
  const [error, setError] = useState<string>('');

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginRequest>();

  const mutation = useMutation({
    mutationFn: authApi.login,
    onSuccess: (data) => {
      setToken(data.jwt_token);
      navigate('/');
    },
    onError: (err: any) => {
      setError(err.response?.data?.error || 'Ошибка входа');
    },
  });

  const onSubmit = (data: LoginRequest) => {
    setError('');
    mutation.mutate(data);
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-dark p-4">
      <div className="max-w-md w-full animate-fade-in">
        <Card className="p-8">
          <div className="text-center mb-8">
            <div className="flex justify-center mb-4">
              <div className="p-3 rounded-xl bg-gradient-primary shadow-glow">
                <Music className="h-10 w-10 text-white" />
              </div>
            </div>
            <h1 className="text-3xl font-bold text-white mb-2">Вход</h1>
            <p className="text-gray-400">Войдите в свою учетную запись</p>
          </div>

          {error && (
            <div className="mb-4 p-4 bg-red-500/20 border border-red-500/30 rounded-xl text-red-400">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
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
              })}
              error={errors.password?.message}
            />

            <Button
              type="submit"
              variant="primary"
              isLoading={mutation.isPending}
              className="w-full"
            >
              Войти
            </Button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-gray-400">
              Нет аккаунта?{' '}
              <Link to="/register" className="text-primary-400 hover:text-primary-300 transition-colors font-medium">
                Зарегистрироваться
              </Link>
            </p>
          </div>
        </Card>
      </div>
    </div>
  );
}

