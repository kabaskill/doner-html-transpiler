interface DictionaryResponse {
  tags: Record<string, string>;
  attributes: Record<string, string>;
}

interface DictionarySectionProps {
  showDictionary: boolean;
  setShowDictionary: (show: boolean) => void;
  dictionary: DictionaryResponse | null;
}

export default function DictionarySection({ showDictionary, setShowDictionary, dictionary }: DictionarySectionProps) {
  return (
    <div className="bg-gray-50 p-6 rounded-lg">
      <div className="flex justify-between items-center mb-3">
        <h3 className="text-lg font-semibold text-gray-700">How it works</h3>
        <button
          onClick={() => setShowDictionary(!showDictionary)}
          className="px-3 py-1 text-sm bg-blue-100 text-blue-700 rounded hover:bg-blue-200 transition-colors"
        >
          {showDictionary ? "Hide Dictionary" : "Show Full Dictionary"}
        </button>
      </div>
      <p className="text-gray-600 mb-3">
        This tool converts German HTML syntax to standard HTML.
        {showDictionary ? "Here are all supported translations:" : "Some example translations:"}
      </p>

      {!showDictionary && (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
          <div>
            <strong>Basic Structure:</strong>
            <ul className="list-disc list-inside text-gray-600 mt-1">
              <li>
                <code>&lt;döner&gt;</code> → <code>&lt;html&gt;</code>
              </li>
              <li>
                <code>&lt;kopf&gt;</code> → <code>&lt;head&gt;</code>
              </li>
              <li>
                <code>&lt;körper&gt;</code> → <code>&lt;body&gt;</code>
              </li>
              <li>
                <code>&lt;titel&gt;</code> → <code>&lt;title&gt;</code>
              </li>
            </ul>
          </div>
          <div>
            <strong>Content Elements:</strong>
            <ul className="list-disc list-inside text-gray-600 mt-1">
              <li>
                <code>&lt;hauptüberschrift&gt;</code> → <code>&lt;h1&gt;</code>
              </li>
              <li>
                <code>&lt;absatz&gt;</code> → <code>&lt;p&gt;</code>
              </li>
              <li>
                <code>&lt;liste&gt;</code> → <code>&lt;ul&gt;</code>
              </li>
              <li>
                <code>&lt;listenelement&gt;</code> → <code>&lt;li&gt;</code>
              </li>
            </ul>
          </div>
        </div>
      )}

      {showDictionary && dictionary && (
        <div className="space-y-6">
          {/* Tags Dictionary */}
          <div>
            <h4 className="text-md font-semibold text-gray-700 mb-2">
              HTML Tags ({Object.keys(dictionary.tags).length} supported)
            </h4>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2 text-sm max-h-64 overflow-y-auto bg-white p-4 rounded border">
              {Object.entries(dictionary.tags)
                .sort(([a], [b]) => a.localeCompare(b))
                .map(([german, english]) => (
                  <div
                    key={german}
                    className="flex gap-4 items-center py-1 border-b border-gray-100 last:border-b-0"
                  >
                    <code className="text-blue-600">&lt;{german}&gt;</code>
                    <span className="text-gray-400">→</span>
                    <code className="text-green-600">&lt;{english}&gt;</code>
                  </div>
                ))}
            </div>
          </div>

          {/* Attributes Dictionary */}
          <div>
            <h4 className="text-md font-semibold text-gray-700 mb-2">
              HTML Attributes ({Object.keys(dictionary.attributes).length} supported)
            </h4>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2 text-sm max-h-48 overflow-y-auto bg-white p-4 rounded border">
              {Object.entries(dictionary.attributes)
                .sort(([a], [b]) => a.localeCompare(b))
                .map(([german, english]) => (
                  <div
                    key={german}
                    className="flex gap-4 items-center py-1 border-b border-gray-100 last:border-b-0"
                  >
                    <code className="text-blue-600">{german}</code>
                    <span className="text-gray-400">→</span>
                    <code className="text-green-600">{english}</code>
                  </div>
                ))}
            </div>
          </div>
        </div>
      )}

      {showDictionary && !dictionary && (
        <div className="text-center py-4">
          <p className="text-gray-500">Loading dictionary...</p>
        </div>
      )}
    </div>
  );
}
