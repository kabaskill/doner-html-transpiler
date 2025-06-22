
/// <reference types="../vite-env.d.ts" />
import React, { useMemo, useCallback, useRef } from "react";
import {
  createEditor,
  Descendant,
  Text,
  Element as SlateElement,
} from "slate";
import {
  Slate,
  Editable,
  withReact,
  RenderElementProps,
  RenderLeafProps,
} from "slate-react";
import { withHistory } from "slate-history";

interface SlateEditorProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  className?: string;
}

// Convert plain text to Slate value
const textToSlate = (text: string): Descendant[] => {
  if (!text.trim()) {
    return [
      {
        type: "paragraph",
        children: [{ text: "" }],
      } as any,
    ];
  }

  // Split by lines and create paragraphs
  const lines = text.split('\n');
  return lines.map(line => ({
    type: "paragraph",
    children: [{ text: line }],
  } as any));
};

// Convert Slate value to plain text
const slateToText = (value: Descendant[]): string => {
  return value
    .map(node => {
      if (Text.isText(node)) {
        return node.text;
      }
      if (SlateElement.isElement(node)) {
        return node.children.map(child => 
          Text.isText(child) ? child.text : ''
        ).join('');
      }
      return '';
    })
    .join('\n');
};

// Simple element renderer - everything is just a paragraph
const Element = ({ attributes, children }: RenderElementProps) => {
  return (
    <p {...attributes} className="my-1">
      {children}
    </p>
  );
};

// Simple leaf renderer - no formatting
const Leaf = ({ attributes, children }: RenderLeafProps) => {
  return <span {...attributes}>{children}</span>;
};

// Main simplified Slate editor component
const SlateEditor: React.FC<SlateEditorProps> = ({ 
  value, 
  onChange, 
  placeholder = "Enter your text here...",
  className = ""
}) => {
  // Create a stable editor instance
  const editor = useMemo(() => withHistory(withReact(createEditor())), []);
  const valueRef = useRef(value);
  
  // Convert external value to Slate format
  const slateValue = useMemo(() => {
    valueRef.current = value;
    return textToSlate(value);
  }, [value]);
  
  const renderElement = useCallback(
    (props: RenderElementProps) => <Element {...props} />,
    [],
  );
  
  const renderLeaf = useCallback(
    (props: RenderLeafProps) => <Leaf {...props} />,
    [],
  );

  const handleChange = useCallback((newValue: Descendant[]) => {
    const textValue = slateToText(newValue);
    // Only call onChange if the value actually changed and it's different from external value
    if (textValue !== valueRef.current) {
      onChange(textValue);
    }
  }, [onChange]);

  return (
    <div className={`slate-editor ${className}`}>
      <Slate
        key={value} // Force re-render when value changes
        editor={editor}
        initialValue={slateValue}
        onChange={handleChange}
      >
        <Editable
          renderElement={renderElement}
          renderLeaf={renderLeaf}
          placeholder={placeholder}
          spellCheck
          className="min-h-[300px] w-full rounded border border-gray-300 p-4 font-mono text-sm leading-relaxed focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50"
          style={{ fontFamily: 'Monaco, Menlo, "Ubuntu Mono", monospace' }}
        />
      </Slate>
    </div>
  );
};

export default SlateEditor;