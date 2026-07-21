import { Frown } from 'lucide-react';
import { Link } from 'react-router-dom';

interface ErrorPageProps {
  status?: number;
  message?: string;
}

const ErrorPage = ({ status = 404, message = 'Profile not found' }: ErrorPageProps) => (
  <div className="min-h-screen flex flex-col items-center justify-center px-4 bg-gray-50">
    <Frown className="h-16 w-16 text-gray-300 mb-4" />
    <h1 className="text-4xl font-bold text-gray-900 mb-2">{status}</h1>
    <p className="text-gray-500 text-lg mb-6">{message}</p>
    <Link
      to="/"
      className="text-indigo-600 hover:text-indigo-500 font-medium text-sm"
    >
      Go to LinkPulse
    </Link>
  </div>
);

export default ErrorPage;