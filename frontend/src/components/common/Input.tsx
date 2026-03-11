import React, { type InputHTMLAttributes, forwardRef } from 'react';
import { clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';
import { EyeIcon, EyeSlashIcon } from '@heroicons/react/24/outline';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  hint?: string;
  leftIcon?: React.ReactNode;
  rightIcon?: React.ReactNode;
  showPasswordToggle?: boolean;
}

const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, hint, leftIcon, rightIcon, showPasswordToggle, className, type = 'text', ...props }, ref) => {
    const [showPassword, setShowPassword] = React.useState(false);
    const inputType = showPasswordToggle ? (showPassword ? 'text' : 'password') : type;

    const inputClasses = twMerge(
      clsx(
        'w-full px-4 py-2 border rounded-lg transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-offset-0',
        error
          ? 'border-red-300 focus:border-red-500 focus:ring-red-200'
          : 'border-gray-300 focus:border-primary-500 focus:ring-primary-200',
        leftIcon && 'pl-10',
        (rightIcon || showPasswordToggle) && 'pr-10',
        className
      )
    );

    return (
      <div className="w-full">
        {label && (
          <label className="block text-sm font-medium text-gray-700 mb-1">
            {label}
          </label>
        )}

        <div className="relative">
          {leftIcon && (
            <div className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400">
              {leftIcon}
            </div>
          )}

          <input
            ref={ref}
            type={inputType}
            className={inputClasses}
            {...props}
          />

          {showPasswordToggle && (
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 transition-colors"
            >
              {showPassword ? (
                <EyeSlashIcon className="w-5 h-5" />
              ) : (
                <EyeIcon className="w-5 h-5" />
              )}
            </button>
          )}

          {rightIcon && !showPasswordToggle && (
            <div className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400">
              {rightIcon}
            </div>
          )}
        </div>

        {error && (
          <p className="mt-1 text-sm text-red-600">{error}</p>
        )}

        {hint && !error && (
          <p className="mt-1 text-sm text-gray-500">{hint}</p>
        )}
      </div>
    );
  }
);

Input.displayName = 'Input';

export default Input;
