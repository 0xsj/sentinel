export const DENSITIES = ['comfortable', 'compact'] as const;
export type Density = (typeof DENSITIES)[number];

const STORAGE_KEY = 'atelier.density';

function readStoredDensity(): Density {
	try {
		const value = globalThis.localStorage?.getItem(STORAGE_KEY);
		if (value && DENSITIES.includes(value as Density)) return value as Density;
	} catch {
		// Storage can be unavailable in a restricted preview context.
	}
	return 'comfortable';
}

export function applyDensity(value: Density): void {
	const root = globalThis.document?.documentElement;
	if (!root) return;
	if (value === 'comfortable') delete root.dataset.density;
	else root.dataset.density = value;
}

export function setDensity(value: Density): void {
	try {
		globalThis.localStorage?.setItem(STORAGE_KEY, value);
	} catch {
		// The visual preference still applies when storage is unavailable.
	}
	applyDensity(value);
}

export function hydrateDensity(): Density {
	const value = readStoredDensity();
	applyDensity(value);
	return value;
}
