import { InputHTMLAttributes, TextareaHTMLAttributes, forwardRef } from 'react';

interface BaseInputProps {
  label?: string;
  error?: string;
}

type InputProps = BaseInputProps &
  (
    | (InputHTMLAttributes<HTMLInputElement> & { multiline?: false; rows?: never })
    | (TextareaHTMLAttributes<HTMLTextAreaElement> & { multiline: true; rows?: number })
  );

const Input = forwardRef<HTMLInputElement | HTMLTextAreaElement, InputProps>(
  ({ label, error, className = '', ...props }, ref) => {
    if (props.multiline) {
      const { multiline, rows = 4, ...textareaProps } = props as any;
      return (
        <div className="w-full">
          {label && (
            <label className="block text-sm font-medium text-gray-700 mb-1">
              {label}
            </label>
          )}
          <textarea
            ref={ref as any}
            rows={rows}
            className={`input min-h-[100px] ${error ? 'border-red-500' : ''} ${className}`}
            {...textareaProps}
          />
          {error && <p className="mt-1 text-sm text-red-600">{error}</p>}
        </div>
      );
    }

    return (
      <div className="w-full">
        {label && (
          <label className="block text-sm font-medium text-gray-700 mb-1">
            {label}
          </label>
        )}
        <input
          ref={ref as any}
          className={`input ${error ? 'border-red-500' : ''} ${className}`}
          {...(props as InputHTMLAttributes<HTMLInputElement>)}
        />
        {error && <p className="mt-1 text-sm text-red-600">{error}</p>}
      </div>
    );
  }
);

Input.displayName = 'Input';

export default Input;
