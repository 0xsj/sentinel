export const THEMES = ['system', 'light', 'dark'] as const;
export type Theme = (typeof THEMES)[number];

const STORAGE_KEY = 'atelier.theme';

function readStoredTheme(): Theme {
	try {
		const value = globalThis.localStorage?.getItem(STORAGE_KEY);
		if (value && THEMES.includes(value as Theme)) return value as Theme;
	} catch {
		// Storage can be unavailable in a restricted preview context.
	}
	return 'system';
}

export function applyTheme(value: Theme): void {
	const root = globalThis.document?.documentElement;
	if (!root) return;
	if (value === 'system') delete root.dataset.theme;
	else root.dataset.theme = value;
}

export function setTheme(value: Theme): void {
	try {
		globalThis.localStorage?.setItem(STORAGE_KEY, value);
	} catch {
		// The visual preference still applies when storage is unavailable.
	}
	applyTheme(value);
}

export function hydrateTheme(): Theme {
	const value = readStoredTheme();
	applyTheme(value);
	return value;
}
