import { useState, use, Suspense } from "react";
import Footer from "./components/Footer";
import InputSection from "./components/InputSection";
import OutputSection from "./components/OutputSection";
import DictionarySection from "./components/DictionarySection";
import { sanitizeHTML, hasSecurityRisk } from "./utils/htmlSecurity";

interface TranspileResponse {
  result?: string;
  error?: string;
}

interface DictionaryResponse {
  tags: Record<string, string>;
  attributes: Record<string, string>;
}

const API_BASE = "https://doner-html-transpiler.onrender.com/"

const createDictionaryPromise = (): Promise<DictionaryResponse | null> => {
  return fetch(`${API_BASE}/dictionary`)
    .then((response) => {
      if (!response.ok) {
        throw new Error("Failed to fetch dictionary");
      }
      return response.json() as Promise<DictionaryResponse>;
    })
    .catch((error) => {
      console.error("Dictionary fetch error:", error);
      return null;
    });
};

const dictionaryPromise = createDictionaryPromise();

function AppContent() {
  const [germanHtml, setGermanHtml] = useState("");
  const [standardHtml, setStandardHtml] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");
  const [showDictionary, setShowDictionary] = useState(false);

  const dictionary = use(dictionaryPromise);

  const transpileHtml = async () => {
    if (!germanHtml.trim()) {
      setError("Please enter some German HTML to transpile");
      return;
    }

    setIsLoading(true);
    setError("");
    setStandardHtml("");

    try {
      const response = await fetch(`${API_BASE}/transpile`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ content: germanHtml }),
      });

      const data: TranspileResponse = await response.json();

      if (!response.ok || data.error) {
        throw new Error(data.error || "Failed to transpile");
      }

      setStandardHtml(data.result || "");
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    } finally {
      setIsLoading(false);
    }
  };

  const clearAll = () => {
    setGermanHtml("");
    setStandardHtml("");
    setError("");
  };

  const exampleGermanHtml = `<döner>
  <kopf>
    <titel>Meine Deutsche Webseite</titel>
    <beschreibung>Eine Beispielseite mit deutschen HTML-Tags</beschreibung>
  </kopf>
  <körper>
    <hauptüberschrift>Willkommen!</hauptüberschrift>
    <absatz>Dies ist ein Beispiel für deutsche HTML-Tags.</absatz>
    <liste>
      <listenelement>Erstes Element</listenelement>
      <listenelement>Zweites Element</listenelement>
    </liste>
  </körper>
</döner>`;

  return (
    <main className="relative flex min-h-screen flex-col">
      <section className="flex flex-1 flex-col items-center justify-center p-4">
        <div className="w-full max-w-7xl space-y-8 rounded-lg bg-white p-6 shadow-lg">
          <div className="text-center">
            <h1 className="text-4xl font-bold text-gray-800 mb-2">{`</ 🥙 >`}</h1>
            <h2 className="text-3xl font-semibold text-gray-700 mb-1">D.Ö.N.E.R</h2>

            <p className="text-gray-600 text-xl">
              <span className="font-bold">D</span>eutsche <span className="font-bold">Ö</span>ffnung
              zur <span className="font-bold">N</span>ormalisierten
              <span className="font-bold"> ER</span>kennung von Webseiten
            </p>
          </div>
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Input Section */}
            <InputSection
              germanHtml={germanHtml}
              setGermanHtml={setGermanHtml}
              exampleGermanHtml={exampleGermanHtml}
            />

            {/* Output Section */}
            <OutputSection
              standardHtml={standardHtml}
              hasSecurityRisk={hasSecurityRisk}
              sanitizeHTML={sanitizeHTML}
            />
          </div>
          {/* Controls */}
          <div className="flex flex-col sm:flex-row gap-4 justify-center items-center">
            <button
              onClick={transpileHtml}
              disabled={isLoading || !germanHtml.trim()}
              className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors font-medium"
            >
              {isLoading ? "Transpiling..." : "Transpile HTML"}
            </button>
            <button
              onClick={clearAll}
              className="px-6 py-3 bg-gray-500 text-white rounded-lg hover:bg-gray-600 transition-colors"
            >
              Clear All
            </button>
          </div>
          {/* Error Display */}
          {error && (
            <div className="p-4 bg-red-100 border border-red-400 text-red-700 rounded-lg">
              <strong>Error:</strong> {error}
            </div>
          )} 
          {/* Info Section */}
          <DictionarySection
            showDictionary={showDictionary}
            setShowDictionary={setShowDictionary}
            dictionary={dictionary}
          />
        </div>
      </section>
      <Footer />
    </main>
  );
}

export default function App() {
  return (
    <Suspense
      fallback={
        <main className="flex min-h-screen flex-col items-center justify-center">
          <div className="text-center space-y-4">
            <h1 className="text-4xl font-bold text-white">{`</ 🥙 >`}</h1>
            <h2 className="text-3xl font-semibold text-white">D.Ö.N.E.R</h2>
            <p className="text-white">Loading...</p>
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
          </div>
        </main>
      }
    >
      <AppContent />
    </Suspense>
  );
}
