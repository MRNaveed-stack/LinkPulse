import { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { ProfileUpdateSchema } from '../../schemas/profile';
import { getProfile } from '../../api/profile';
import Button from '../common/Button';
import { User, Mail, FileText, Link, CheckCircle, XCircle, AlertCircle } from 'lucide-react';

interface ProfileFormProps {
  initialData: {
    display_name: string;
    username: string;
    email: string;
    bio: string;
    avatar_url: string;
  };
  onSubmit: (data: any) => void;
  isLoading: boolean;
  onFormChange: (data: any) => void;
}

export default function ProfileForm({ initialData, onSubmit, isLoading, onFormChange }: ProfileFormProps) {
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors, isDirty },
  } = useForm({
    resolver: zodResolver(ProfileUpdateSchema),
    defaultValues: initialData,
  });

  const usernameValue = watch('username');
  const displayNameValue = watch('display_name');
  const bioValue = watch('bio');
  const avatarUrlValue = watch('avatar_url');
  const emailValue = watch('email');

  const [isCheckingUsername, setIsCheckingUsername] = useState(false);
  const [usernameAvailable, setUsernameAvailable] = useState<boolean | null>(null);

  // Sync state for live preview
  useEffect(() => {
    onFormChange({
      display_name: displayNameValue,
      username: usernameValue,
      bio: bioValue,
      avatar_url: avatarUrlValue,
      email: emailValue,
    });
  }, [displayNameValue, usernameValue, bioValue, avatarUrlValue, emailValue, onFormChange]);

  // Debounced Username Availability Check
  useEffect(() => {
    if (!usernameValue || usernameValue === initialData.username) {
      setUsernameAvailable(null);
      return;
    }

    setIsCheckingUsername(true);
    const delayDebounce = setTimeout(async () => {
      // Validate local format first before sending backend request
      if (!/^[a-zA-Z0-9_]+$/.test(usernameValue) || usernameValue.length < 2) {
        setUsernameAvailable(false);
        setIsCheckingUsername(false);
        return;
      }

      try {
        await getProfile(usernameValue);
        setUsernameAvailable(false); // Taken if resolves successfully
      } catch (err: any) {
        if (err?.response?.status === 404) {
          setUsernameAvailable(true); // Available if 404 not found
        } else {
          setUsernameAvailable(null);
        }
      } finally {
        setIsCheckingUsername(false);
      }
    }, 500);

    return () => clearTimeout(delayDebounce);
  }, [usernameValue, initialData.username]);

  // Prevent leaving warning if dirty
  useEffect(() => {
    const handleBeforeUnload = (e: BeforeUnloadEvent) => {
      if (isDirty) {
        e.preventDefault();
        e.returnValue = 'You have unsaved changes. Leave anyway?';
      }
    };
    window.addEventListener('beforeunload', handleBeforeUnload);
    return () => window.removeEventListener('beforeunload', handleBeforeUnload);
  }, [isDirty]);

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      {/* Display Name */}
      <div>
        <label htmlFor="display_name" className="block text-sm font-medium text-gray-700">
          Display Name
        </label>
        <div className="mt-1 relative rounded-md shadow-sm">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <User className="h-4 w-4 text-gray-400" />
          </div>
          <input
            type="text"
            id="display_name"
            placeholder="e.g. Muhammad Naveed"
            {...register('display_name')}
            className={`block w-full pl-10 pr-3 py-2 border rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm ${
              errors.display_name ? 'border-red-300' : 'border-gray-300'
            }`}
          />
        </div>
        {errors.display_name && (
          <p className="mt-1.5 text-xs text-red-600 flex items-center gap-1">
            <AlertCircle className="h-3 w-3" />
            {errors.display_name.message as string}
          </p>
        )}
      </div>

      {/* Username */}
      <div>
        <label htmlFor="username" className="block text-sm font-medium text-gray-700">
          Username
        </label>
        <div className="mt-1 relative rounded-md shadow-sm flex">
          <span className="inline-flex items-center px-3 rounded-l-md border border-r-0 border-gray-300 bg-gray-50 text-gray-500 text-sm">
            linkpulse.com/u/
          </span>
          <div className="relative flex-1">
            <input
              type="text"
              id="username"
              placeholder="username"
              {...register('username')}
              className={`block w-full px-3 py-2 border rounded-r-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm ${
                errors.username ? 'border-red-300' : 'border-gray-300'
              }`}
            />
            <div className="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
              {isCheckingUsername && (
                <div className="h-4 w-4 border-2 border-indigo-500 border-t-transparent rounded-full animate-spin"></div>
              )}
              {!isCheckingUsername && usernameAvailable === true && (
                <CheckCircle className="h-4 w-4 text-emerald-500" />
              )}
              {!isCheckingUsername && usernameAvailable === false && (
                <XCircle className="h-4 w-4 text-red-500" />
              )}
            </div>
          </div>
        </div>
        {usernameAvailable === true && (
          <p className="mt-1.5 text-xs text-emerald-600">✓ Username is available</p>
        )}
        {usernameAvailable === false && (
          <p className="mt-1.5 text-xs text-red-600">✗ Username is already taken</p>
        )}
        {errors.username && (
          <p className="mt-1.5 text-xs text-red-600 flex items-center gap-1">
            <AlertCircle className="h-3 w-3" />
            {errors.username.message as string}
          </p>
        )}
      </div>

      {/* Email */}
      <div>
        <label htmlFor="email" className="block text-sm font-medium text-gray-700">
          Email Address
        </label>
        <div className="mt-1 relative rounded-md shadow-sm">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <Mail className="h-4 w-4 text-gray-400" />
          </div>
          <input
            type="email"
            id="email"
            placeholder="you@example.com"
            {...register('email')}
            className={`block w-full pl-10 pr-3 py-2 border rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm ${
              errors.email ? 'border-red-300' : 'border-gray-300'
            }`}
          />
        </div>
        {errors.email && (
          <p className="mt-1.5 text-xs text-red-600 flex items-center gap-1">
            <AlertCircle className="h-3 w-3" />
            {errors.email.message as string}
          </p>
        )}
      </div>

      {/* Bio */}
      <div>
        <label htmlFor="bio" className="block text-sm font-medium text-gray-700">
          Bio
        </label>
        <div className="mt-1 relative rounded-md shadow-sm">
          <div className="absolute top-3 left-3 flex items-start pointer-events-none">
            <FileText className="h-4 w-4 text-gray-400" />
          </div>
          <textarea
            id="bio"
            rows={3}
            placeholder="Brief profile introduction..."
            {...register('bio')}
            className={`block w-full pl-10 pr-3 py-2 border rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm ${
              errors.bio ? 'border-red-300' : 'border-gray-300'
            }`}
          />
        </div>
        <p className="mt-1.5 text-xs text-gray-400">Brief summary to introduce yourself to profile visitors.</p>
        {errors.bio && (
          <p className="mt-1.5 text-xs text-red-600 flex items-center gap-1">
            <AlertCircle className="h-3 w-3" />
            {errors.bio.message as string}
          </p>
        )}
      </div>

      {/* Avatar URL */}
      <div>
        <label htmlFor="avatar_url" className="block text-sm font-medium text-gray-700">
          Avatar Image URL
        </label>
        <div className="mt-1 relative rounded-md shadow-sm">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <Link className="h-4 w-4 text-gray-400" />
          </div>
          <input
            type="url"
            id="avatar_url"
            placeholder="https://example.com/avatar.jpg"
            {...register('avatar_url')}
            className={`block w-full pl-10 pr-3 py-2 border rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm ${
              errors.avatar_url ? 'border-red-300' : 'border-gray-300'
            }`}
          />
        </div>
        <p className="mt-1.5 text-xs text-gray-400">Direct link to your profile picture.</p>
        {errors.avatar_url && (
          <p className="mt-1.5 text-xs text-red-600 flex items-center gap-1">
            <AlertCircle className="h-3 w-3" />
            {errors.avatar_url.message as string}
          </p>
        )}
      </div>

      {/* Submit Button */}
      <div className="flex justify-end pt-4 border-t border-gray-100">
        <Button
          type="submit"
          loading={isLoading}
          disabled={usernameAvailable === false}
          className="px-6 py-2 shadow-md hover:shadow-lg transition duration-200"
        >
          Save Changes
        </Button>
      </div>
    </form>
  );
}
