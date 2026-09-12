<script lang="ts">
	import type { Snippet } from 'svelte';
	import { Toggle as Primitive } from 'bits-ui';
	import './navigation.css';
	type ToggleSize = 'sm' | 'md';
	interface Props { pressed?: boolean; size?: ToggleSize; disabled?: boolean; children?: Snippet; class?: string; }
	let { pressed = $bindable(false), size = 'md', disabled = false, children, class: className }: Props = $props();
</script>
<Primitive.Root bind:pressed {disabled} class="atelier-toggle {className ?? ''}" data-size={size}>{@render children?.()}</Primitive.Root>
