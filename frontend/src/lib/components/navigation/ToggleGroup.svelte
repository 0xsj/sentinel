<script lang="ts">
	import { ToggleGroup as Primitive } from 'bits-ui';
	import './navigation.css';
	export type ToggleOption = { value: string; label: string; disabled?: boolean };
	type ToggleGroupType = 'single' | 'multiple';
	interface Props { items: readonly ToggleOption[]; type?: ToggleGroupType; value?: string | string[]; orientation?: 'horizontal' | 'vertical'; ariaLabel?: string; class?: string; }
	let { items, type = 'single', value = $bindable(''), orientation = 'horizontal', ariaLabel, class: className }: Props = $props();
</script>
{#if type === 'multiple'}
	<Primitive.Root type="multiple" value={Array.isArray(value) ? value : []} orientation={orientation} aria-label={ariaLabel} class="atelier-toggle-group {className ?? ''}" onValueChange={(next) => (value = next)}>{#each items as item (item.value)}<Primitive.Item value={item.value} disabled={item.disabled} class="atelier-toggle-group__item">{item.label}</Primitive.Item>{/each}</Primitive.Root>
{:else}
	<Primitive.Root type="single" value={typeof value === 'string' ? value : ''} orientation={orientation} aria-label={ariaLabel} class="atelier-toggle-group {className ?? ''}" onValueChange={(next) => (value = next)}>{#each items as item (item.value)}<Primitive.Item value={item.value} disabled={item.disabled} class="atelier-toggle-group__item">{item.label}</Primitive.Item>{/each}</Primitive.Root>
{/if}
