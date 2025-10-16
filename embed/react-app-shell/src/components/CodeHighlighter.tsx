import { startTransition, useEffect, useState } from "react";
import hljs from 'highlight.js/lib/core';
import json from 'highlight.js/lib/languages/json';
import 'highlight.js/styles/github-dark.css';

hljs.registerLanguage('json', json);

export default function HighlightedJson({ code }: { code: string }) {
  const [html, setHtml] = useState("");

  useEffect(() => {
    const highlighted = hljs.highlight(code, { language: 'json' }).value;
    startTransition(() => {
      setHtml(highlighted);
    })
  }, [code]);

  return (
    <pre className="hljs p-4 rounded bg-neutral-800 overflow-auto">
      <code dangerouslySetInnerHTML={{ __html: html }} />
    </pre>
  );
}