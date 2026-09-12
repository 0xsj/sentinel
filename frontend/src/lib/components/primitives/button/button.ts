import type { Snippet } from 'svelte';
import type { HTMLButtonAttributes } from 'svelte/elements';

export const BUTTON_VARIANTS = ['primary', 'secondary', 'quiet', 'danger'] as const;
export type ButtonVariant = (typeof BUTTON_VARIANTS)[number];

export const BUTTON_SIZES = ['sm', 'md', 'lg'] as const;
export type ButtonSize = (typeof BUTTON_SIZES)[number];

export interface ButtonProps extends Omit<HTMLButtonAttributes, 'children'> {
	variant?: ButtonVariant;
	size?: ButtonSize;
	loading?: boolean;
	children?: Snippet;
	class?: string;
}
