import { startTransition, useEffect, useState } from "react";
import hljs from 'highlight.js/lib/core';
import json from 'highlight.js/lib/languages/json';
import 'highlight.js/styles/github-dark.css';

import { copyToClipboard, prettyPrintJson } from "./helpers";
import { Button } from "./ui/button";
import { Check, Copy, Download } from "lucide-react";
import { toast } from "sonner";

hljs.registerLanguage('json', json);

export default function HighlightedJson({ code }: { code: string }) {
  const [html, setHtml] = useState("");
  const [isCopy, setIsCopy] = useState(false);

  useEffect(() => {
    console.log(prettyPrintJson(code))
    const highlighted = hljs.highlight(prettyPrintJson(code), { language: 'json' }).value;
    startTransition(() => {
      setHtml(highlighted);
    })
  }, [code]);

  return (
    <pre className="hljs p-4 relative rounded bg-neutral-800 overflow-auto">
      <div className="absolute right-2 top-2 flex gap-1">
        <Button variant="secondary" size="icon" onClick={() => {
          const blob = new Blob([prettyPrintJson(code)], { type: "application/json" });
          const url = URL.createObjectURL(blob);

          const a = document.createElement("a");
          a.href = url;
          a.download = (new Date()).toDateString()+".json";
          a.click();

          URL.revokeObjectURL(url);
        }}>
          <Download />
        </Button>
        <Button variant="secondary" size="icon" onClick={() => {
          toast.promise(copyToClipboard(prettyPrintJson(code)), {
            loading: "Copying...",
            success: (data) => {
              setIsCopy(true);
              setTimeout(() => {
                setIsCopy(false)
              }, 3000);
              return `${data.name} copy successful`;
            },
            error: "Copy failed",
          });
        }}>
          {isCopy ? <Check className="text-green-600" /> : <Copy />}
        </Button>
      </div>
      <code dangerouslySetInnerHTML={{ __html: html }} />
    </pre>
  );
}