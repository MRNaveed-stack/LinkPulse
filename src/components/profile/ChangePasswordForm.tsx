import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { PasswordChangeSchema } from '../../schemas/profile';
import { useChangePassword } from '../../hooks/useChangePassword';
import Button from '../common/Button';
import { Key, Lock, AlertCircle } from 'lucide-react';

export default function ChangePasswordForm() {
  const changePasswordMutation = useChangePassword();

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm({
    resolver: zodResolver(PasswordChangeSchema),
    defaultValues: {
      current_password: '',
      new_password: '',
      confirm_password: '',
    },
  });

  const onSubmit = (data: any) => {
    changePasswordMutation.mutate(
      {
        current_password: data.current_password,
        new_password: data.new_password,
      },
      {
        onSuccess: () => {
          reset();
        },
      }
    );
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      {/* Current Password */}
      <div>
        <label htmlFor="current_password" className="block text-sm font-medium text-gray-700">
          Current Password
        </label>
        <div className="mt-1 relative rounded-md shadow-sm">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <Lock className="h-4 w-4 text-gray-400" />
          </div>
          <input
            type="password"
            id="current_password"
            placeholder="••••••••"
            {...register('current_password')}
            className={`block w-full pl-10 pr-3 py-2 border rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm ${
              errors.current_password ? 'border-red-300' : 'border-gray-300'
            }`}
          />
        </div>
        {errors.current_password && (
          <p className="mt-1.5 text-xs text-red-600 flex items-center gap-1">
            <AlertCircle className="h-3 w-3" />
            {errors.current_password.message as string}
          </p>
        )}
      </div>

      {/* New Password */}
      <div>
        <label htmlFor="new_password" className="block text-sm font-medium text-gray-700">
          New Password
        </label>
        <div className="mt-1 relative rounded-md shadow-sm">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <Key className="h-4 w-4 text-gray-400" />
          </div>
          <input
            type="password"
            id="new_password"
            placeholder="••••••••"
            {...register('new_password')}
            className={`block w-full pl-10 pr-3 py-2 border rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm ${
              errors.new_password ? 'border-red-300' : 'border-gray-300'
            }`}
          />
        </div>
        <p className="mt-1.5 text-xs text-gray-400">
          Must be at least 8 characters, with 1 uppercase letter, 1 lowercase letter, and 1 number.
        </p>
        {errors.new_password && (
          <p className="mt-1.5 text-xs text-red-600 flex items-center gap-1">
            <AlertCircle className="h-3 w-3" />
            {errors.new_password.message as string}
          </p>
        )}
      </div>

      {/* Confirm Password */}
      <div>
        <label htmlFor="confirm_password" className="block text-sm font-medium text-gray-700">
          Confirm New Password
        </label>
        <div className="mt-1 relative rounded-md shadow-sm">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <Key className="h-4 w-4 text-gray-400" />
          </div>
          <input
            type="password"
            id="confirm_password"
            placeholder="••••••••"
            {...register('confirm_password')}
            className={`block w-full pl-10 pr-3 py-2 border rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm ${
              errors.confirm_password ? 'border-red-300' : 'border-gray-300'
            }`}
          />
        </div>
        {errors.confirm_password && (
          <p className="mt-1.5 text-xs text-red-600 flex items-center gap-1">
            <AlertCircle className="h-3 w-3" />
            {errors.confirm_password.message as string}
          </p>
        )}
      </div>

      {/* Action Buttons */}
      <div className="flex justify-end pt-4 border-t border-gray-100">
        <Button
          type="submit"
          loading={changePasswordMutation.isPending}
          className="px-6 py-2 shadow-md hover:shadow-lg transition duration-200"
        >
          Change Password
        </Button>
      </div>
    </form>
  );
}
