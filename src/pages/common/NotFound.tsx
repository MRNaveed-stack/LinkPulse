import { Link } from 'react-router-dom';
import { AlertCircle, ArrowLeft, Home } from 'lucide-react';
import Button from '../../components/common/Button';

export default function NotFound() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-indigo-50/50 via-slate-50 to-purple-50/50 p-6 font-sans">
      <div className="bg-white border border-gray-200/80 shadow-2xl rounded-2xl p-8 max-w-lg w-full text-center relative overflow-hidden">
        {/* Glow decoration */}
        <div className="absolute -top-24 -left-24 w-48 h-48 bg-indigo-500/5 rounded-full blur-3xl pointer-events-none"></div>
        <div className="absolute -bottom-24 -right-24 w-48 h-48 bg-purple-500/5 rounded-full blur-3xl pointer-events-none"></div>

        <div className="mx-auto flex items-center justify-center h-20 w-20 rounded-full bg-indigo-50 text-indigo-600 mb-6 border border-indigo-100">
          <AlertCircle className="h-10 w-10 animate-bounce" />
        </div>

        <h1 className="text-8xl font-black text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 via-purple-600 to-pink-600 mb-2">
          404
        </h1>
        <h2 className="text-2xl font-bold text-gray-900 mb-4">Page Not Found</h2>
        
        <p className="text-gray-500 text-sm leading-relaxed mb-8 max-w-md mx-auto">
          The page you are looking for doesn't exist, was removed, or had its name changed. Let's get you back on track.
        </p>

        <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
          <Button
            variant="secondary"
            onClick={() => window.history.back()}
            className="w-full sm:w-auto flex items-center justify-center gap-2 border-gray-300 text-gray-700 hover:bg-gray-50"
          >
            <ArrowLeft className="h-4 w-4" />
            Go Back
          </Button>
          <Link to="/dashboard" className="w-full sm:w-auto">
            <Button className="w-full sm:w-auto flex items-center justify-center gap-2 bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 shadow-lg shadow-indigo-500/15">
              <Home className="h-4 w-4" />
              Dashboard Home
            </Button>
          </Link>
        </div>
      </div>
    </div>
  );
}
