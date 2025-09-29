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