import { useEffect, useRef, useState } from "react";
import SecurityWarning from "./SecurityWarning";
import { cn } from "../utils/cn";

interface OutputSectionProps {
  standardHtml: string;
  hasSecurityRisk: (html: string) => string | null;
  sanitizeHTML: (html: string) => string;
}

export default function OutputSection({
  standardHtml,
  hasSecurityRisk,
  sanitizeHTML,
}: OutputSectionProps) {
  const [expanded, setExpanded] = useState(false);
  const tooltipRef = useRef<HTMLDivElement | null>(null);

  function handleClickOutside(event: MouseEvent) {
    const target = event.target as HTMLElement;

    if (expanded && !target.closest("iframe")) {
      if (tooltipRef.current && !tooltipRef.current.contains(target)) {
        setExpanded(false);
      }
    }
  }

  useEffect(() => {
    if (expanded) {
      document.addEventListener("mousedown", handleClickOutside);
      document.addEventListener("keydown", (e) => {
        if (e.key === "Escape") {
          setExpanded(false);
        }
      });
    }
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
      document.removeEventListener("keydown", (e) => {
        if (e.key === "Escape") {
          setExpanded(false);
        }
      });
    };
  }, [expanded]);

  return (
    <div className="flex flex-col gap-4 h-full">
      <h2 className="text-xl font-semibold text-gray-700">Standard HTML Output</h2>

      {standardHtml ? (
        <div className="w-full h-full border border-gray-300 rounded-lg bg-white overflow-auto">
          {/* Rendered HTML */}
          <div className="p-4 border-b border-gray-200">
            <div className="flex justify-between items-center mb-2">
              <h3 className="text-sm font-semibold text-gray-600 mb-2">Preview:</h3>
              <button
                type="button"
                onClick={() => setExpanded(true)}
                className="hover:scale-110 transition-transform cursor-pointer"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="black"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path d="M15 3h6v6" />
                  <path d="M10 14 21 3" />
                  <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6" />
                </svg>
              </button>
            </div>

            {(() => {
              const securityRisk = hasSecurityRisk(standardHtml);
              if (securityRisk) {
                return <SecurityWarning securityRisk={securityRisk} />;
              } else {
                return (
                  <div
                    ref={tooltipRef}
                    className={cn(
                      expanded && "absolute inset-0 w-4/5 h-4/5 m-auto shadow-2xl",
                      expanded ? "rounded-xl" : "rounded",
                      "max-w-none border border-slate-400 p-4 bg-gray-50"
                    )}
                  >
                    <iframe
                      srcDoc={sanitizeHTML(standardHtml)}
                      className="w-full h-full border-0"
                      sandbox="allow-same-origin"
                      title="HTML Preview"
                    />
                  </div>
                );
              }
            })()}
          </div>

          {/* Raw HTML Code */}
          <div className="p-4">
            <h3 className="text-sm font-semibold text-gray-600 mb-2">HTML Code:</h3>
            <pre className="text-xs font-mono text-gray-700 whitespace-pre-wrap bg-gray-50 p-3 rounded">
              {standardHtml}
            </pre>
          </div>
        </div>
      ) : (
        <div className="w-full h-full border border-gray-300 rounded-lg bg-gray-50 flex items-center justify-center">
          <p className="text-gray-500 text-center">
            Transpiled HTML will appear here...
            <br />
            <span className="text-sm">You'll see both the rendered preview and the HTML code</span>
          </p>
        </div>
      )}
    </div>
  );
}
