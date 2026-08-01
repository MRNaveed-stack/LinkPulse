import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useForgotPassword } from '../../hooks/useForgotPassword';
import { forgotPasswordSchema } from '../../schemas/auth';
import { ZodError } from 'zod';

export default function ForgotPassword() {
  const [email, setEmail] = useState('');
  const [error, setError] = useState('');
  const forgotMutation = useForgotPassword();
  const [tokenSent, setTokenSent] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setTokenSent('');
    try {
      forgotPasswordSchema.parse({ email });
    } catch (err) {
      if (err instanceof ZodError) {
        setError(err.issues[0].message);
        return;
      }
    }

    forgotMutation.mutate(
      { email },
      {
        onSuccess: (response) => {
          const devToken = response.data.token;
          setTokenSent(
            devToken
              ? `Dev token: ${devToken} (use this to reset password)`
              : `If this email exists, a reset link has been sent.`
          );
        },
      }
    );
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-indigo-50/50 via-slate-50 to-purple-50/50 p-4 font-sans">
      <div className="bg-white border border-gray-200/80 shadow-2xl rounded-2xl p-8 max-w-md w-full">
        <h2 className="text-3xl font-extrabold text-center text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-purple-600 mb-2">
          Reset Password
        </h2>
        <p className="text-gray-500 text-sm text-center mb-8">
          Enter your email address and we'll send you a link to reset your password.
        </p>

        <form onSubmit={handleSubmit} className="space-y-5">
          <div>
            <label className="block text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">
              Email Address
            </label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full bg-gray-50 border border-gray-300 rounded-xl px-4 py-3 text-gray-900 placeholder-gray-400 focus:outline-none focus:bg-white focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all duration-200"
              placeholder="name@example.com"
            />
            {error && <p className="text-xs text-rose-500 mt-1.5 font-medium">{error}</p>}
          </div>

          <button
            type="submit"
            disabled={forgotMutation.isPending}
            className="w-full py-3 px-4 rounded-xl text-white font-bold bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 active:scale-[0.98] transition-all duration-150 cursor-pointer shadow-lg shadow-indigo-500/20 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {forgotMutation.isPending ? 'Sending...' : 'Send Reset Link'}
          </button>
        </form>

        {forgotMutation.isError && (
          <p className="text-sm text-center text-rose-500 font-medium mt-4">
            Error: {(forgotMutation.error as any)?.response?.data?.message || 'Request failed'}
          </p>
        )}

        {tokenSent && (
          <div className="mt-6 p-4 bg-indigo-50 border border-indigo-100 rounded-xl">
            <p className="text-sm text-indigo-800 text-center font-medium leading-relaxed">
              {tokenSent}
            </p>
            {tokenSent.includes('Dev token') && (
              <div className="mt-3 text-center">
                <Link
                  to={`/reset-password?token=${tokenSent.split('Dev token: ')[1].split(' ')[0]}`}
                  className="text-xs text-indigo-600 hover:text-indigo-700 font-bold underline transition-colors"
                >
                  Go directly to Reset Page
                </Link>
              </div>
            )}
          </div>
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