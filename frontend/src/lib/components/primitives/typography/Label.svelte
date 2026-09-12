<script lang="ts">
	import type { Snippet } from 'svelte';
	import { Label as Primitive } from 'bits-ui';
	import type { HTMLLabelAttributes } from 'svelte/elements';
	import './typography.css';

	interface Props extends Omit<HTMLLabelAttributes, 'children' | 'id'> {
		id?: string;
		children?: Snippet;
		class?: string;
	}

	let { children, class: className, ...rest }: Props = $props();
</script>

<Primitive.Root {...rest} class="atelier-label {className ?? ''}">
	{@render children?.()}
</Primitive.Root>
