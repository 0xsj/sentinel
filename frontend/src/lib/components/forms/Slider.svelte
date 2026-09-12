<script lang="ts">
	import { Slider as Primitive } from 'bits-ui';
	import './forms.css';

	interface Props { value?: number; min?: number; max?: number; step?: number; orientation?: 'horizontal' | 'vertical'; disabled?: boolean; ariaLabel?: string; class?: string; }
	let { value = $bindable(50), min = 0, max = 100, step = 1, orientation = 'horizontal', disabled = false, ariaLabel = 'Slider', class: className }: Props = $props();
</script>

<Primitive.Root type="single" bind:value {min} {max} {step} {orientation} {disabled} aria-label={ariaLabel} class="atelier-slider {className ?? ''}" data-orientation={orientation}>{#snippet children({ thumbItems })}<Primitive.Track class="atelier-slider__track"><Primitive.Range class="atelier-slider__range" /></Primitive.Track>{#each thumbItems as { index }}<Primitive.Thumb {index} class="atelier-slider__thumb" />{/each}{/snippet}</Primitive.Root>
