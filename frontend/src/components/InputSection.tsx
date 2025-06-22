import SlateEditor from "./SlateEditor";

interface InputSectionProps {
  germanHtml: string;
  setGermanHtml: (value: string) => void;
  exampleGermanHtml: string;
}

export default function InputSection({ germanHtml, setGermanHtml, exampleGermanHtml }: InputSectionProps) {
  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h2 className="text-xl font-semibold text-gray-700">Döner Input</h2>
        <button
          onClick={() => setGermanHtml(exampleGermanHtml)}
          className="px-3 py-1 text-sm bg-blue-100 text-blue-700 rounded hover:bg-blue-200 transition-colors"
        >
          Load Example
        </button>
      </div>
      <SlateEditor
        value={germanHtml}
        onChange={setGermanHtml}
        className="w-full text-gray-700"
      />
    </div>
  );
}
