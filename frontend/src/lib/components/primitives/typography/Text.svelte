<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';
	import './typography.css';

	type TextTag = 'p' | 'span' | 'div' | 'small';
	type TextSize = 'sm' | 'md' | 'lg';
	type TextTone = 'default' | 'muted' | 'quiet' | 'accent' | 'danger';
	type TextWeight = 'regular' | 'medium' | 'strong';

	interface Props extends Omit<HTMLAttributes<HTMLElement>, 'children'> {
		as?: TextTag;
		size?: TextSize;
		tone?: TextTone;
		weight?: TextWeight;
		measure?: boolean;
		truncate?: boolean;
		children?: Snippet;
		class?: string;
	}

	let { as: tag = 'p', size = 'md', tone = 'default', weight = 'regular', measure = false, truncate = false, children, class: className, ...rest }: Props = $props();
</script>

<svelte:element this={tag} {...rest} class="atelier-text {className ?? ''}" data-size={size} data-tone={tone} data-weight={weight} data-measure={measure || undefined} data-truncate={truncate || undefined}>
	{@render children?.()}
</svelte:element>
