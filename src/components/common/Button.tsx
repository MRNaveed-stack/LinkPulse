import React from 'react';
import Spinner from './Spinner';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'danger';
  loading?: boolean;
  children: React.ReactNode;
}

const variantClasses = {
  primary: 'bg-indigo-600 hover:bg-indigo-700 text-white',
  secondary: 'bg-white hover:bg-gray-50 text-gray-700 border border-gray-300',
  danger: 'bg-red-600 hover:bg-red-700 text-white',
};

const Button = ({ 
  variant = 'primary', 
  loading = false, 
  children, 
  className = '', 
  disabled, 
  ...props 
}: ButtonProps) => (
  <button
    className={`inline-flex items-center justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed ${variantClasses[variant]} ${className}`}
    disabled={disabled || loading}
    {...props}
  >
    {loading && <Spinner size="sm" className="mr-2" />}
    {children}
  </button>
);

export default Button;