interface CardProps {
  children: React.ReactNode;
  className?: string;
  onClick?: () => void;
  hover?: boolean;
}

export default function Card({ children, className = '', onClick, hover = false }: CardProps) {
  const cardClass = hover || onClick ? 'card-hover' : 'card';
  
  return (
    <div
      className={`${cardClass} ${className}`}
      onClick={onClick}
    >
      {children}
    </div>
  );
}

