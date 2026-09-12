<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';
	import './typography.css';

	type HeadingSize = 'sm' | 'md' | 'lg' | 'xl';

	interface Props extends Omit<HTMLAttributes<HTMLHeadingElement>, 'children'> {
		level: 1 | 2 | 3 | 4 | 5 | 6;
		size?: HeadingSize;
		children: Snippet;
		class?: string;
	}

	let { level, size = 'md', children, class: className, ...rest }: Props = $props();
	const tag = $derived(`h${level}`);
</script>

<svelte:element this={tag} {...rest} class="atelier-heading {className ?? ''}" data-size={size}>
	{@render children()}
</svelte:element>
