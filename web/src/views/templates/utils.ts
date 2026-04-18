export function formatDate(value?: string) {
  if (!value) return '--';
  return new Date(value).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  });
}

export function splitGuideSteps(value?: string) {
  if (!value) return [] as string[];

  const normalized = value
    .replace(/^[^:：]*[:：]\s*/, '')
    .replace(/\r?\n/g, ' ')
    .trim();

  return normalized
    .split(/\s*(?:\d+\.\s*|->|=>|，|,|；|;|。)\s*/)
    .map((item) => item.trim())
    .filter(Boolean);
}

export function uniqueValues(values: string[]) {
  return Array.from(new Set(values.map((item) => item.trim()).filter(Boolean)));
}
