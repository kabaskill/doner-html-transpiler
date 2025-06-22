interface SecurityWarningProps {
  securityRisk: string;
}

export default function SecurityWarning({ securityRisk }: SecurityWarningProps) {
  return (
    <div className="border rounded p-4 bg-red-50 border-red-200">
      <div className="flex items-center gap-2 text-red-700 mb-2">
        <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
          <path
            fillRule="evenodd"
            d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z"
            clipRule="evenodd"
          />
        </svg>
        <span className="font-semibold">Security Warning</span>
      </div>
      <p className="text-sm text-red-600 mb-2">
        Preview disabled for security: {securityRisk}
      </p>
      <p className="text-xs text-red-500">
        The HTML code is shown below, but preview is blocked to prevent XSS attacks.
      </p>
    </div>
  );
}
