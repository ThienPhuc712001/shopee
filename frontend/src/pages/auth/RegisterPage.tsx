import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useAppDispatch, useAppSelector } from '../../hooks';
import { register, clearError } from '../../store/slices/authSlice';
import { ROUTES } from '../../constants';
import Input from '../../components/common/Input';
import Button from '../../components/common/Button';

const registerSchema = z.object({
  first_name: z.string().min(1, 'First name is required'),
  last_name: z.string().min(1, 'Last name is required'),
  email: z.string().email('Invalid email address'),
  phone: z.string().min(10, 'Phone number must be at least 10 digits'),
  password: z.string().min(8, 'Password must be at least 8 characters'),
  confirm_password: z.string(),
}).refine((data) => data.password === data.confirm_password, {
  message: "Passwords don't match",
  path: ['confirm_password'],
});

type RegisterForm = z.infer<typeof registerSchema>;

const RegisterPage: React.FC = () => {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const { isAuthenticated, isLoading, error } = useAppSelector((state) => state.auth);

  const {
    register: registerForm,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterForm>({
    resolver: zodResolver(registerSchema),
  });

  React.useEffect(() => {
    if (isAuthenticated) {
      navigate(ROUTES.HOME, { replace: true });
    }
  }, [isAuthenticated, navigate]);

  React.useEffect(() => {
    return () => {
      dispatch(clearError());
    };
  }, [dispatch]);

  const onSubmit = async (data: RegisterForm) => {
    try {
      await dispatch(register({
        email: data.email,
        password: data.password,
        first_name: data.first_name,
        last_name: data.last_name,
        phone: data.phone,
      })).unwrap();
      navigate(ROUTES.HOME, { replace: true });
    } catch (err) {
      // Error is handled by Redux slice
    }
  };

  return (
    <div className="min-h-screen bg-gray-100 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
      <div className="sm:mx-auto sm:w-full sm:max-w-md">
        <Link to={ROUTES.HOME} className="flex justify-center">
          <div className="w-12 h-12 bg-primary-600 rounded-lg flex items-center justify-center">
            <span className="text-white font-bold text-2xl">E</span>
          </div>
        </Link>
        <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
          Create your account
        </h2>
        <p className="mt-2 text-center text-sm text-gray-600">
          Already have an account?{' '}
          <Link to={ROUTES.LOGIN} className="font-medium text-primary-600 hover:text-primary-500">
            Sign in
          </Link>
        </p>
      </div>

      <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
        <div className="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
          <form className="space-y-4" onSubmit={handleSubmit(onSubmit)}>
            {error && (
              <div className="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg">
                {error}
              </div>
            )}

            <div className="grid grid-cols-2 gap-4">
              <Input
                label="First name"
                autoComplete="given-name"
                error={errors.first_name?.message}
                {...registerForm('first_name')}
              />
              <Input
                label="Last name"
                autoComplete="family-name"
                error={errors.last_name?.message}
                {...registerForm('last_name')}
              />
            </div>

            <Input
              label="Email address"
              type="email"
              autoComplete="email"
              error={errors.email?.message}
              {...registerForm('email')}
            />

            <Input
              label="Phone number"
              type="tel"
              autoComplete="tel"
              error={errors.phone?.message}
              {...registerForm('phone')}
            />

            <Input
              label="Password"
              type="password"
              autoComplete="new-password"
              error={errors.password?.message}
              showPasswordToggle
              {...registerForm('password')}
            />

            <Input
              label="Confirm password"
              type="password"
              autoComplete="new-password"
              error={errors.confirm_password?.message}
              showPasswordToggle
              {...registerForm('confirm_password')}
            />

            <div className="flex items-center">
              <input
                id="terms"
                name="terms"
                type="checkbox"
                required
                className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
              />
              <label htmlFor="terms" className="ml-2 block text-sm text-gray-900">
                I agree to the{' '}
                <a href="/terms" className="text-primary-600 hover:text-primary-500">
                  Terms and Conditions
                </a>
              </label>
            </div>

            <Button
              type="submit"
              fullWidth
              isLoading={isLoading}
            >
              Create account
            </Button>
          </form>
        </div>
      </div>
    </div>
  );
};

export default RegisterPage;
