import { useState, useEffect } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { useLogin } from '../../hooks/useLogin';
import { loginSchema } from '../../schemas/auth';
import { ZodError } from 'zod';
import { GoogleLoginButton } from '../../components/common/GoogleLoginButton';
import { useAuthStore } from '../../store/authStore';
import toast from 'react-hot-toast';

export default function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [errors, setErrors] = useState<{ email?: string; password?: string }>({});
  const loginMutation = useLogin();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { login: storeLogin } = useAuthStore();

  useEffect(() => {
    const accessToken = searchParams.get('access_token');
    const refreshToken = searchParams.get('refresh_token');
    const oauthError = searchParams.get('error');

    if (oauthError) {
      toast.error(`Google Authentication failed: ${oauthError}`);
    } else if (accessToken && refreshToken) {
      try {
        const payload = JSON.parse(atob(accessToken.split('.')[1]));
        storeLogin(
          { id: payload.sub, username: payload.username, email: payload.email },
          { accessToken, refreshToken }
        );
        toast.success('Logged in successfully with Google');
        navigate('/dashboard', { replace: true });
      } catch (err) {
        console.error('Failed to parse token payload:', err);
        toast.error('Failed to resolve Google login session.');
      }
    }
  }, [searchParams, storeLogin, navigate]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setErrors({});
    try {
      loginSchema.parse({ email, password });
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

    loginMutation.mutate(
      { email, password },
      {
        onSuccess: () => navigate('/dashboard', { replace: true }),
      }
    );
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-indigo-50/50 via-slate-50 to-purple-50/50 p-4 font-sans">
      <div className="bg-white border border-gray-200/80 shadow-2xl rounded-2xl p-8 max-w-md w-full">
        <h2 className="text-3xl font-extrabold text-center text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-purple-600 mb-2">
          Welcome Back
        </h2>
        <p className="text-gray-500 text-sm text-center mb-8">
          Sign in to access your LinkPulse dashboard
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
            {errors.email && <p className="text-xs text-rose-500 mt-1.5 font-medium">{errors.email}</p>}
          </div>

          <div>
            <div className="flex justify-between items-center mb-2">
              <label className="block text-xs font-semibold text-gray-500 uppercase tracking-wider">
                Password
              </label>
              <Link
                to="/forgot-password"
                className="text-xs text-indigo-600 hover:text-indigo-700 font-semibold transition-colors"
              >
                Forgot Password?
              </Link>
            </div>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full bg-gray-50 border border-gray-300 rounded-xl px-4 py-3 text-gray-900 placeholder-gray-400 focus:outline-none focus:bg-white focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all duration-200"
              placeholder="••••••••"
            />
            {errors.password && <p className="text-xs text-rose-500 mt-1.5 font-medium">{errors.password}</p>}
          </div>

          <button
            type="submit"
            disabled={loginMutation.isPending}
            className="w-full py-3 px-4 rounded-xl text-white font-bold bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 active:scale-[0.98] transition-all duration-150 cursor-pointer shadow-lg shadow-indigo-500/20 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loginMutation.isPending ? 'Logging in...' : 'Sign In'}
          </button>
        </form>

        {loginMutation.isError && (
          <p className="text-sm text-center text-rose-500 font-medium mt-4">
            Error: {(loginMutation.error as any)?.response?.data?.message || 'Login failed'}
          </p>
        )}

        <div className="relative flex py-5 items-center">
          <div className="flex-grow border-t border-gray-200"></div>
          <span className="flex-shrink mx-4 text-gray-400 text-xs font-semibold uppercase tracking-wider">Or</span>
          <div className="flex-grow border-t border-gray-200"></div>
        </div>

        <GoogleLoginButton />

        <p className="text-center text-gray-500 text-sm mt-8">
          Don't have an account?{' '}
          <Link to="/register" className="text-indigo-600 hover:text-indigo-700 font-semibold transition-colors">
            Register
          </Link>
        </p>
      </div>
    </div>
  );
}