import { useState, useEffect } from 'react';
import { useNavigate, useLocation, Link } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { useMutation } from '@tanstack/react-query';
import { authApi } from '../api/auth';
import { useAuth } from '../contexts/AuthContext';
import { VerifyRequest } from '../types';
import Button from '../components/Button';
import Input from '../components/Input';
import { CheckCircle } from 'lucide-react';

export default function VerifyPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const { setToken } = useAuth();
  const [error, setError] = useState<string>('');
  const [sessionId, setSessionId] = useState<string>('');

  useEffect(() => {
    const stateSessionId = location.state?.sessionId;
    if (stateSessionId) {
      setSessionId(stateSessionId);
    } else {
      navigate('/register');
    }
  }, [location, navigate]);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<VerifyRequest>({
    defaultValues: {
      session_id: sessionId,
    },
  });

  const mutation = useMutation({
    mutationFn: authApi.verify,
    onSuccess: (data) => {
      setToken(data.jwt_token);
      navigate('/');
    },
    onError: (err: any) => {
      setError(err.response?.data?.error || 'Ошибка верификации');
    },
  });

  const onSubmit = (data: VerifyRequest) => {
    setError('');
    mutation.mutate({ ...data, session_id: sessionId });
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary-50 to-primary-100">
      <div className="max-w-md w-full">
        <div className="bg-white rounded-2xl shadow-xl p-8">
          <div className="text-center mb-8">
            <div className="flex justify-center mb-4">
              <CheckCircle className="h-12 w-12 text-primary-600" />
            </div>
            <h1 className="text-3xl font-bold text-gray-900">Подтверждение</h1>
            <p className="text-gray-600 mt-2">
              Введите код подтверждения, отправленный на ваш email
            </p>
          </div>

          {error && (
            <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg text-red-700">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <Input
              label="Код подтверждения"
              {...register('code', {
                required: 'Код обязателен',
                pattern: {
                  value: /^\d{4}$/,
                  message: 'Код должен состоять из 4 цифр',
                },
              })}
              error={errors.code?.message}
              placeholder="1234"
              maxLength={4}
            />

            <Button
              type="submit"
              variant="primary"
              isLoading={mutation.isPending}
              className="w-full"
            >
              Подтвердить
            </Button>
          </form>

          <div className="mt-6 text-center">
            <Link to="/register" className="text-primary-600 hover:underline text-sm">
              Вернуться к регистрации
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}

