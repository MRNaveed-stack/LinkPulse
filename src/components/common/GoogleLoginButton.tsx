export const GoogleLoginButton = () => {
  const handleGoogleLogin = () => {
    // Redirect to backend endpoint that initiates Google OAuth flow
    window.location.href = '/api/auth/google/login';
  };

  return (
    <button
      onClick={handleGoogleLogin}
      className="w-full flex justify-center py-2 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 mt-4 cursor-pointer"
    >
      Continue with Google
    </button>
  );
};
