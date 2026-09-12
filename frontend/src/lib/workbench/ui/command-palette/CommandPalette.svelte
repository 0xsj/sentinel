<script lang="ts">
	import { Dialog as Primitive } from 'bits-ui';
	import { Input } from '../../../components/forms';
	import './command-palette.css';
	export type CommandItem = { id: string; label: string; description?: string; shortcut?: string; disabled?: boolean };
	interface Props { commands: readonly CommandItem[]; open?: boolean; onExecute?: (command: CommandItem) => void; class?: string; }
	let { commands, open = $bindable(false), onExecute, class: className }: Props = $props();
	let query = $state('');
	let selectedIndex = $state(0);
	const filteredCommands = $derived(commands.filter((command) => `${command.label} ${command.description ?? ''}`.toLowerCase().includes(query.trim().toLowerCase())));
	$effect(() => { if (selectedIndex >= filteredCommands.length) selectedIndex = Math.max(filteredCommands.length - 1, 0); });
	function execute(index: number): void { const command = filteredCommands[index]; if (!command || command.disabled) return; onExecute?.(command); open = false; }
	function handleKeydown(event: KeyboardEvent): void { if (event.key === 'ArrowDown') { event.preventDefault(); selectedIndex = Math.min(selectedIndex + 1, Math.max(filteredCommands.length - 1, 0)); } else if (event.key === 'ArrowUp') { event.preventDefault(); selectedIndex = Math.max(selectedIndex - 1, 0); } else if (event.key === 'Enter') { event.preventDefault(); execute(selectedIndex); } }
</script>
<Primitive.Root bind:open><Primitive.Portal><Primitive.Overlay class="atelier-command-palette__overlay" /><Primitive.Content class="atelier-command-palette {className ?? ''}" onkeydown={handleKeydown}><div class="atelier-command-palette__header"><div><Primitive.Title class="atelier-command-palette__title">Command palette</Primitive.Title><Primitive.Description class="atelier-command-palette__description">Search and invoke a registered workbench command.</Primitive.Description></div><Primitive.Close class="atelier-command-palette__close" aria-label="Close command palette">×</Primitive.Close></div><Input bind:value={query} class="atelier-command-palette__input" placeholder="Search commands" aria-label="Search commands" autofocus /><div class="atelier-command-palette__list" role="listbox" aria-label="Commands">{#each filteredCommands as command, index (command.id)}<button type="button" class="atelier-command-palette__item" class:is-selected={selectedIndex === index} disabled={command.disabled} role="option" aria-selected={selectedIndex === index} onclick={() => execute(index)}><span class="atelier-command-palette__item-copy"><strong>{command.label}</strong>{#if command.description}<small>{command.description}</small>{/if}</span>{#if command.shortcut}<kbd>{command.shortcut}</kbd>{/if}</button>{:else}<p class="atelier-command-palette__empty">No commands match “{query}”.</p>{/each}</div></Primitive.Content></Primitive.Portal></Primitive.Root>
