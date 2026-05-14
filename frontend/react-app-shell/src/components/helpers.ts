export function getInitials(name: string) {
  if (!name || typeof name !== 'string') return '';

  const trimmed = name.trim();
  if (!trimmed) return '';

  const names = trimmed.split(/\s+/);
  if (names.length === 1) {
    return names[0].charAt(0).toUpperCase();
  } else {
    return (names[0].charAt(0) + names[1].charAt(0)).toUpperCase();
  }
}

export function prettyPrintJson(jsonString: string) {
  try {
    const obj = JSON.parse(jsonString);
    return JSON.stringify(obj, null, 2); // 2 = indentation spaces
  } catch (e) {
    return jsonString; // fallback if invalid JSON
  }
}

export async function copyToClipboard(text: string) {
  return navigator.clipboard.writeText(text).then(() => ({
    name: "Clipboard"
  }));
}

export function extractMessages(obj: any) {
  const messages: any[] = [];

  function traverse(obj: { [x: string]: { message: any; }; }) {
    for (const key in obj) {
      if (typeof obj[key] === "object") {
        traverse(obj[key]);
      }
      if (obj[key]?.message) {
        messages.push(obj[key].message);
      }
    }
  }

  traverse(obj);
  return messages;
}