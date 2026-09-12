<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLButtonAttributes } from 'svelte/elements';
	import type { ButtonSize } from './button';
	import './button.css';

	interface Props extends Omit<HTMLButtonAttributes, 'children'> {
		size?: ButtonSize;
		label: string;
		children: Snippet;
		class?: string;
	}

	let { size = 'md', label, children, class: className, ...rest }: Props = $props();
</script>

<button
	{...rest}
	class="atelier-icon-button {className ?? ''}"
	data-size={size}
	aria-label={label}
	title={rest.title ?? label}
>
	{@render children()}
</button>
