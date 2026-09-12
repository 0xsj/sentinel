<script lang="ts">
	import type { Snippet } from 'svelte';
	import './shell.css';
	interface Props { title?: string; activity?: Snippet; sidebar?: Snippet; titleActions?: Snippet; main: Snippet; status?: Snippet; class?: string; }
	let { title = 'Atelier', activity, sidebar, titleActions, main, status, class: className }: Props = $props();
</script>
<div class="atelier-workbench-shell {className ?? ''}" data-sidebar={sidebar ? 'visible' : 'hidden'}><header class="atelier-workbench__titlebar"><strong class="atelier-workbench__title">{title}</strong>{#if titleActions}<div class="atelier-workbench__title-actions">{@render titleActions()}</div>{/if}</header><div class="atelier-workbench__body">{#if activity}<aside class="atelier-workbench__activity" aria-label="Activity rail">{@render activity()}</aside>{/if}{#if sidebar}<aside class="atelier-workbench__sidebar" aria-label="Sidebar">{@render sidebar()}</aside>{/if}<main class="atelier-workbench__main">{@render main()}</main></div>{#if status}<footer class="atelier-workbench__status">{@render status()}</footer>{/if}</div>
