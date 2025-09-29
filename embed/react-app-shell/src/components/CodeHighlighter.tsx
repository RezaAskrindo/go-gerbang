import { useEffect, useState } from "react";
import { codeToHtml } from 'shiki'
import { Button } from "@/components/ui/button";
import { Clipboard } from "lucide-react";
import { toast } from "sonner";

type CodeHighlighterProps = {
  code: string;
  language?: string;
};

export default function CodeHighlighter({ code, language = "json" }: CodeHighlighterProps) {
  const [html, setHtml] = useState<string>("");

  useEffect(() => {
    const highlight = async () => {
      const highlighted = await codeToHtml(code, { lang: language, theme: "github-dark-high-contrast" });
      setHtml(highlighted);
    };

    highlight();
  }, [code, language]);

  const handleCopyHtml = () => {
    navigator.clipboard.writeText(code)
      .then(() => toast.success("Copy to Clipboard"))
      .catch(err => console.error("Copy failed", err));
  };

  return <div className="relative min-h-10 max-h-96">
    <Button className="bg-code absolute top-3 right-2 z-10 size-7 hover:opacity-100 focus-visible:opacity-100" size="icon" onClick={handleCopyHtml}>
      <Clipboard />
    </Button>
    <div dangerouslySetInnerHTML={{ __html: html }} />
  </div>
}