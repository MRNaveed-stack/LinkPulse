import { useState } from 'react';
import { useSearchParams, useNavigate, Link } from 'react-router-dom';
import { useResetPassword } from '../../hooks/useResetPassword';
import { resetPasswordSchema } from '../../schemas/auth';
import { ZodError } from 'zod';
import toast from 'react-hot-toast';

export default function ResetPassword() {
  const [searchParams] = useSearchParams();
  const tokenFromUrl = searchParams.get('token') || '';
  const [token, setToken] = useState(tokenFromUrl);
  const [newPassword, setNewPassword] = useState('');
  const [errors, setErrors] = useState<any>({});
  const resetMutation = useResetPassword();
  const navigate = useNavigate();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setErrors({});
    try {
      resetPasswordSchema.parse({ token, new_password: newPassword });
    } catch (err) {
      if (err instanceof ZodError) {
        const fieldErrors: any = {};
        err.issues.forEach((e) => {
          fieldErrors[e.path[0]] = e.message;
        });
        setErrors(fieldErrors);
        return;
      }
    }
    resetMutation.mutate(
      { token, new_password: newPassword },
      {
        onSuccess: () => {
          toast.success('Password reset successful! You can now login.');
          navigate('/login');
        },
      }
    );
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-indigo-50/50 via-slate-50 to-purple-50/50 p-4 font-sans">
      <div className="bg-white border border-gray-200/80 shadow-2xl rounded-2xl p-8 max-w-md w-full">
        <h2 className="text-3xl font-extrabold text-center text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-purple-600 mb-2">
          New Password
        </h2>
        <p className="text-gray-500 text-sm text-center mb-8">
          Enter your reset token and new password to recover access.
        </p>

        <form onSubmit={handleSubmit} className="space-y-5">
          <div>
            <label className="block text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">
              Reset Token
            </label>
            <input
              type="text"
              value={token}
              onChange={(e) => setToken(e.target.value)}
              className="w-full bg-gray-50 border border-gray-300 rounded-xl px-4 py-3 text-gray-900 placeholder-gray-400 focus:outline-none focus:bg-white focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all duration-200"
              placeholder="Paste your reset token here"
            />
            {errors.token && <p className="text-xs text-rose-500 mt-1.5 font-medium">{errors.token}</p>}
          </div>

          <div>
            <label className="block text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">
              New Password
            </label>
            <input
              type="password"
              value={newPassword}
              onChange={(e) => setNewPassword(e.target.value)}
              className="w-full bg-gray-50 border border-gray-300 rounded-xl px-4 py-3 text-gray-900 placeholder-gray-400 focus:outline-none focus:bg-white focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all duration-200"
              placeholder="••••••••"
            />
            {errors.new_password && <p className="text-xs text-rose-500 mt-1.5 font-medium">{errors.new_password}</p>}
          </div>

          <button
            type="submit"
            disabled={resetMutation.isPending}
            className="w-full py-3 px-4 rounded-xl text-white font-bold bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 active:scale-[0.98] transition-all duration-150 cursor-pointer shadow-lg shadow-indigo-500/20 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {resetMutation.isPending ? 'Resetting...' : 'Reset Password'}
          </button>
        </form>

        {resetMutation.isError && (
          <p className="text-sm text-center text-rose-500 font-medium mt-4">
            Error: {(resetMutation.error as any)?.response?.data?.message || 'Reset failed'}
          </p>
        )}

        <p className="text-center text-gray-500 text-sm mt-8">
          Back to{' '}
          <Link to="/login" className="text-indigo-600 hover:text-indigo-700 font-semibold transition-colors">
            Login
          </Link>
        </p>
      </div>
    </div>
  );
}
