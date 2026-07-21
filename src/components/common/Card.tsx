import React from 'react';

interface CardProps {
  children: React.ReactNode;
  className?: string;
  onClick?: () => void;
}

const Card = ({ children, className = '', onClick }: CardProps) => (
  <div 
    className={`bg-white shadow rounded-lg overflow-hidden ${className}`}
    onClick={onClick}
  >
    {children}
  </div>
);

export default Card;