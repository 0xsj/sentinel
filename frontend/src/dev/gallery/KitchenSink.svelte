<script lang="ts">
	import { onMount } from 'svelte';
	import { List, ScrollArea, Tree } from '../../lib/components/collections';
	import { Avatar, CodeBlock, CopyButton, InlineCode, Tag } from '../../lib/components/content';
	import { Card, KeyValue, Table, type TableColumn, type TableRow } from '../../lib/components/data-display';
	import { Checkbox, Input, RadioGroup, Select, Switch, Textarea } from '../../lib/components/forms';
	import { Alert, Badge, Banner, EmptyState, Progress, Skeleton, Spinner, Status, Toaster, type ToastItem } from '../../lib/components/feedback';
	import { Accordion, Breadcrumb, Collapsible, ContextMenu, DropdownMenu, MenuItem, MenuSeparator, TabPanel, Tabs, Toggle, ToggleGroup, type AccordionItem, type BreadcrumbItem } from '../../lib/components/navigation';
	import { Dialog, Popover, Tooltip } from '../../lib/components/overlays';
	import { Button, Heading, IconButton, Kbd, Label, Separator, Surface, Text } from '../../lib/components/primitives';
	import { Resizable } from '../../lib/components/layout';
	import { CommandPalette, type CommandItem } from '../../lib/workbench/ui/command-palette';
	import { WorkbenchPanel } from '../../lib/workbench/ui/panels';
	import { WorkbenchShell } from '../../lib/workbench/ui/shell';
	import { ViewHost, type WorkbenchView } from '../../lib/workbench/ui/view-host';
	import {
		DENSITIES,
		THEMES,
		hydrateDensity,
		hydrateTheme,
		setDensity,
		setTheme,
		type Density,
		type Theme
	} from '../../lib/runtime';
	import './kitchen-sink.css';

	let { host }: { host: string } = $props();
	let theme = $state<Theme>('system');
	let density = $state<Density>('comfortable');
	let formChecked = $state(true);
	let formPinned = $state(false);
	let formChoice = $state('canvas');
	let demoTab = $state('overview');
	let workbenchView = $state('overview');
	let paletteOpen = $state(false);
	let lastCommand = $state('No command selected');
	let menuRadio = $state('system');
	let toggleValue = $state('grid');
	let compactLabels = $state(false);
	let lastMenuAction = $state('No menu action selected');
	let listSelection = $state('tokens');
	let treeSelection = $state('atelier');
	let expandedTreeIds = $state<string[]>(['atelier']);
	let alertVisible = $state(true);
	let bannerVisible = $state(true);
	let feedbackToasts = $state<ToastItem[]>([{ id: 'indexing', title: 'Workspace indexed', description: '12 files are ready to inspect.', tone: 'success' }]);
	let toastSequence = $state(1);
	const inspectorColumns: readonly TableColumn[] = [{ key: 'name', label: 'View' }, { key: 'owner', label: 'Owner' }, { key: 'updated', label: 'Updated', align: 'end' }, { key: 'status', label: 'Status' }];
	const inspectorRows: readonly TableRow[] = [{ id: 'overview', name: 'Overview', owner: 'You', updated: '2m ago', status: 'Ready' }, { id: 'tokens', name: 'Token inspector', owner: 'You', updated: '18m ago', status: 'Draft' }, { id: 'preview', name: 'Preview', owner: 'Team', updated: 'Yesterday', status: 'Needs review' }];
	const workspaceProperties = [{ label: 'Workspace', value: 'atelier-core', description: 'Local project root' }, { label: 'Branch', value: 'main', description: 'No pending changes' }, { label: 'Views', value: 12, description: '3 pinned views' }, { label: 'Last indexed', value: '08:42:19', description: 'All files up to date' }] as const;
	const tokenSnippet = `:root {
	--accent: oklch(72% 0.17 68);
	--surface-panel: oklch(21% 0.02 260);
}`;
	let copyFeedback = $state('Nothing copied yet');
	let splitSize = $state(38);
	let inspectorOpen = $state(true);
	let accordionValue = $state('tokens');
	const breadcrumbItems: readonly BreadcrumbItem[] = [{ label: 'Atelier', href: '#overview' }, { label: 'Foundations', href: '#tokens' }, { label: 'Tokens' }];
	const accordionItems: readonly AccordionItem[] = [{ value: 'tokens', title: 'Token metadata', content: 'Semantic values are grouped by the surface, ink, and state meaning they carry.' }, { value: 'usage', title: 'Usage guidance', content: 'Prefer semantic tokens at the component boundary and keep feature state outside this layer.' }, { value: 'notes', title: 'Implementation notes', content: 'This section is intentionally controlled so a host can persist the open item.' }];
	const demoTabs = [{ value: 'overview', label: 'Overview' }, { value: 'tokens', label: 'Tokens' }, { value: 'preview', label: 'Preview', disabled: true }] as const;
	const workbenchViews: readonly WorkbenchView[] = [{ id: 'overview', label: 'Overview' }, { id: 'tokens', label: 'Tokens' }, { id: 'preview', label: 'Preview', disabled: true }];
	const paletteCommands: readonly CommandItem[] = [{ id: 'open-overview', label: 'Open overview', description: 'Show the overview view', shortcut: '⌘ 1' }, { id: 'inspect-tokens', label: 'Inspect tokens', description: 'Jump to the token inspector', shortcut: '⌘ 2' }, { id: 'toggle-sidebar', label: 'Toggle sidebar', description: 'Show or hide the Explorer panel', shortcut: '⌘ B' }];
	const formOptions = [{ value: 'canvas', label: 'Canvas view' }, { value: 'tokens', label: 'Token inspector' }, { value: 'preview', label: 'Preview', disabled: true }] as const;
	const radioOptions = [{ value: 'system', label: 'Use system theme' }, { value: 'light', label: 'Light theme' }, { value: 'dark', label: 'Dark theme' }] as const;
	const toggleOptions = [{ value: 'grid', label: 'Grid' }, { value: 'list', label: 'List' }, { value: 'split', label: 'Split' }] as const;
	const collectionItems = [{ id: 'overview', label: 'Overview', description: 'Workspace summary', meta: '⌘1' }, { id: 'tokens', label: 'Tokens', description: 'Semantic values', meta: '⌘2' }, { id: 'preferences', label: 'Preferences', description: 'App settings', meta: '⌘,' }, { id: 'archived', label: 'Archived', description: 'Read-only history', disabled: true }] as const;
	const treeNodes = [{ id: 'atelier', label: 'Atelier', children: [{ id: 'foundations', label: 'Foundations', children: [{ id: 'tokens-node', label: 'Tokens' }, { id: 'components-node', label: 'Components' }] }, { id: 'workbench-node', label: 'Workbench' }, { id: 'preferences-node', label: 'Preferences' }] }] as const;

	const surfaces = [
		'--surface-ground',
		'--surface-rail',
		'--surface-sunk',
		'--surface-panel',
		'--surface-panel-2',
		'--surface-raised'
	] as const;
	const ink = ['--ink', '--ink-muted', '--ink-subtle', '--ink-disabled'] as const;
	const spacing = [
		'--space-1',
		'--space-2',
		'--space-3',
		'--space-4',
		'--space-5',
		'--space-6',
		'--space-7',
		'--space-8',
		'--space-9',
		'--space-10'
	] as const;
	const status = [
		{ name: 'success', color: '--success', tint: '--success-tint' },
		{ name: 'warning', color: '--warning', tint: '--warning-tint' },
		{ name: 'danger', color: '--danger', tint: '--danger-tint' }
	] as const;

	onMount(() => {
		theme = hydrateTheme();
		density = hydrateDensity();
	});

	function chooseTheme(value: Theme): void {
		theme = value;
		setTheme(value);
	}

	function chooseDensity(value: Density): void {
		density = value;
		setDensity(value);
	}

	function runCommand(command: CommandItem): void {
		lastCommand = command.label;
	}

	function chooseMenuAction(action: string): void {
		lastMenuAction = action;
	}

	function showToast(): void {
		const id = `toast-${toastSequence}`;
		toastSequence += 1;
		feedbackToasts = [...feedbackToasts, { id, title: 'Preview queued', description: 'The preview will update when the view is ready.', tone: 'accent' }];
	}

	function dismissToast(id: string): void {
		feedbackToasts = feedbackToasts.filter((toast) => toast.id !== id);
	}

	function showCopyResult(value: string): void {
		copyFeedback = `Copied ${value.length} characters`;
	}
</script>

<svelte:head><title>Kitchen sink · Atelier</title></svelte:head>

<div class="sink-shell">
	<a class="skip-link" href="#gallery-main">Skip to gallery</a>
	<header class="sink-header">
		<div class="brand-lockup">
			<div class="brand-mark" aria-hidden="true">A</div>
			<div>
				<strong>Atelier</strong>
				<span>Design system kitchen sink</span>
			</div>
		</div>
		<div class="header-meta">
			<span>{host}</span>
			<div class="segmented" role="group" aria-label="Theme">
				{#each THEMES as option}
					<button type="button" class:active={theme === option} aria-pressed={theme === option} onclick={() => chooseTheme(option)}>{option}</button>
				{/each}
			</div>
			<div class="segmented" role="group" aria-label="Density">
				{#each DENSITIES as option}
					<button type="button" class:active={density === option} aria-pressed={density === option} onclick={() => chooseDensity(option)}>{option}</button>
				{/each}
			</div>
		</div>
	</header>

	<div class="sink-body">
		<nav class="sink-nav" aria-label="Kitchen sink sections">
			<p class="nav-label">Sections</p>
			<a href="#overview">Overview</a>
			<a href="#tokens">Tokens</a>
			<a href="#structure">Structure</a>
			<a href="#typography">Typography</a>
			<a href="#forms">Forms</a>
			<a href="#feedback">Feedback</a>
			<a href="#feedback-states">Feedback states</a>
			<a href="#data-display">Data display</a>
			<a href="#content-utilities">Content utilities</a>
			<a href="#workbench-interactions">Workbench interactions</a>
			<a href="#navigation">Navigation</a>
			<a href="#overlays">Overlays</a>
			<a href="#commands">Commands</a>
			<a href="#menus">Menus &amp; choices</a>
			<a href="#collections">Collections</a>
			<a href="#controls">Controls</a>
			<a href="#workbench">Workbench</a>
		</nav>

		<main id="gallery-main" class="sink-main" tabindex="-1">
			<section id="overview" class="intro" aria-labelledby="overview-title">
				<p class="eyebrow">Atelier / frontend foundation</p>
				<h1 id="overview-title">A place for every piece.</h1>
				<p class="intro-copy">A small, live view of the values the workbench will read. Theme and density are real token changes, not presentation-only examples.</p>
			</section>

			<section id="tokens" class="gallery-section" aria-labelledby="tokens-title">
				<div class="section-heading"><div><p class="eyebrow">Foundations</p><h2 id="tokens-title">Semantic tokens</h2></div><p>Components read these meanings instead of raw palette values.</p></div>
				<div class="token-grid">
					<div class="token-card token-card-wide"><h3>Surfaces</h3><div class="surface-swatches">{#each surfaces as token}<div class="surface-swatch" style={`--swatch: var(${token})`}><span></span><code>{token}</code></div>{/each}</div></div>
					<div class="token-card"><h3>Ink</h3><div class="ink-list">{#each ink as token}<p style={`color: var(${token})`}><code>{token}</code> Readable interface text</p>{/each}</div></div>
					<div class="token-card"><h3>Status</h3><div class="status-list">{#each status as item}<div class="status-chip" style={`--status-color: var(${item.color}); --status-tint: var(${item.tint})`}><span></span>{item.name}<code>{item.color}</code></div>{/each}</div></div>
				</div>
			</section>

			<section id="structure" class="gallery-section" aria-labelledby="structure-title">
				<div class="section-heading"><div><p class="eyebrow">Primitives</p><h2 id="structure-title">Structural primitives</h2></div><p>Surfaces establish depth; separators establish relationships without adding noise.</p></div>
				<div class="surface-examples"><Surface tone="sunk" padding="md" class="surface-example"><strong>Sunk</strong><span>Recessed workspace or input background.</span></Surface><Surface tone="panel" padding="md" class="surface-example"><strong>Panel</strong><span>Default surface for a workbench region.</span></Surface><Surface tone="raised" padding="md" class="surface-example"><strong>Raised</strong><span>Nearest surface for active controls.</span></Surface></div>
				<div class="separator-demo"><span>Sidebar</span><Separator orientation="vertical" /><span>Main content</span></div>
			</section>

			<section id="typography" class="gallery-section" aria-labelledby="typography-title">
				<div class="section-heading"><div><p class="eyebrow">Primitives</p><h2 id="typography-title">Typography</h2></div><p>Document structure and visual scale stay separate.</p></div>
				<Surface tone="panel" padding="lg" class="typography-demo"><p class="eyebrow">Section label</p><Heading level={3} size="lg">A heading with a deliberate level</Heading><Text size="lg" tone="muted" measure>Readable interface copy has a measure and a tone, while the heading still owns the document outline.</Text><Separator /><div class="typography-field"><Label for="typography-sample">Field label</Label><input id="typography-sample" class="field" value="Typography pairs with controls" /></div></Surface>
			</section>

			<section id="forms" class="gallery-section" aria-labelledby="forms-title">
				<div class="section-heading"><div><p class="eyebrow">Forms</p><h2 id="forms-title">Fields and choices</h2></div><p>Native text entry stays familiar; Bits UI owns the interaction-heavy choices.</p></div>
				<Surface tone="panel" padding="lg" class="forms-demo">
					<div class="form-grid"><div class="form-field"><Label for="form-name">Name</Label><Input id="form-name" placeholder="Untitled workspace" /></div><div class="form-field"><Label for="form-view">Default view</Label><Select id="form-view" bind:value={formChoice} options={formOptions} aria-label="Default view" /></div><div class="form-field form-field-wide"><Label for="form-notes">Notes</Label><Textarea id="form-notes" rows={3} placeholder="A quiet place for context..." /></div></div>
					<Separator />
					<div class="choice-list"><div class="choice-row"><Checkbox bind:checked={formChecked} aria-label="Remember this workspace" /><div class="choice-copy"><strong>Remember this workspace</strong><span>Restore the active view when the app opens.</span></div></div><div class="choice-row"><Switch bind:checked={formPinned} aria-label="Pin workspace" /><div class="choice-copy"><strong>Pin workspace</strong><span>Keep this workspace visible in the activity rail.</span></div></div></div>
				</Surface>
			</section>

			<section id="feedback" class="gallery-section" aria-labelledby="feedback-title">
				<div class="section-heading"><div><p class="eyebrow">Feedback</p><h2 id="feedback-title">State and progress</h2></div><p>Small signals carry status without taking over the workbench.</p></div>
				<Surface tone="panel" padding="lg" class="feedback-demo"><div class="feedback-row"><Status tone="success" pulse>Connected</Status><Status tone="warning">Needs attention</Status><Status tone="danger">Blocked</Status><Status tone="accent">Syncing</Status></div><div class="feedback-row"><Badge tone="neutral">Draft</Badge><Badge tone="accent">Active</Badge><Badge tone="success" variant="solid">Ready</Badge><Badge tone="warning" variant="outline">Review</Badge><Badge tone="danger" size="md">Error</Badge></div><Separator /><div class="progress-list"><div class="progress-item"><div class="progress-meta"><span>Workspace indexing</span><span>68%</span></div><Progress value={68} aria-label="Workspace indexing 68 percent" /></div><div class="progress-item"><div class="progress-meta"><span>Waiting for preview</span><Spinner size="sm" label="Waiting for preview" /></div><Progress aria-label="Waiting for preview" /></div></div></Surface>
			</section>

			<section id="feedback-states" class="gallery-section" aria-labelledby="feedback-states-title"><div class="section-heading"><div><p class="eyebrow">Feedback</p><h2 id="feedback-states-title">Loading, empty and alert states</h2></div><p>Longer-lived states explain what is happening and give the user a clear next step.</p></div><div class="feedback-state-stack">{#if alertVisible}<Alert tone="warning" title="Preview is stale" dismissible onDismiss={() => (alertVisible = false)}>Re-run the preview to see the latest workspace changes.</Alert>{/if}{#if bannerVisible}<Banner tone="accent" title="New workbench update" dismissible onDismiss={() => (bannerVisible = false)}>Panels now remember their last size.{#snippet actions()}<Button variant="quiet" size="sm">Review</Button>{/snippet}</Banner>{/if}<Surface tone="panel" padding="lg" class="feedback-state-grid"><div class="feedback-state-column"><p class="field-label">Loading</p><div class="skeleton-stack"><Skeleton width="78%" /><Skeleton width="48%" /><Skeleton variant="rect" width="100%" height="72px" /></div></div><EmptyState title="Nothing here yet" description="Create a view to start filling this workspace.">{#snippet action()}<Button variant="secondary" size="sm">Create view</Button>{/snippet}</EmptyState></Surface><div class="toast-demo"><Button variant="secondary" onclick={showToast}>Show toast</Button><Text size="sm" tone="muted">The host owns lifecycle; Toaster only renders the current queue.</Text></div><Toaster toasts={feedbackToasts} onDismiss={dismissToast} /></div></section>

			<section id="navigation" class="gallery-section" aria-labelledby="navigation-title"><div class="section-heading"><div><p class="eyebrow">Navigation</p><h2 id="navigation-title">Tabs</h2></div><p>Keyboard-oriented navigation for view-local state, with compact desktop geometry.</p></div><Surface tone="panel" padding="lg" class="navigation-demo"><Tabs items={demoTabs} bind:value={demoTab}><TabPanel value="overview"><div class="tab-panel-copy"><strong>Overview</strong><span>The active view can hold a workbench canvas without changing the surrounding shell.</span></div></TabPanel><TabPanel value="tokens"><div class="tab-panel-copy"><strong>Tokens</strong><span>Local navigation stays independent from application-level routing.</span></div></TabPanel></Tabs></Surface></section>

			<section id="overlays" class="gallery-section" aria-labelledby="overlays-title"><div class="section-heading"><div><p class="eyebrow">Overlays</p><h2 id="overlays-title">Tooltip, popover and dialog</h2></div><p>Short help, contextual actions and focused tasks each get an appropriate overlay.</p></div><Surface tone="panel" padding="lg" class="overlay-demo"><div class="overlay-actions"><Tooltip content="This is a keyboard-accessible hint" class="tooltip-demo-trigger"><span>Hover or focus me</span></Tooltip><Popover>{#snippet trigger()}<span class="popover-demo-trigger">Open popover</span>{/snippet}<div class="popover-demo-card"><strong>Quick view</strong><span>Contextual actions can sit close to the thing they affect.</span></div></Popover><Dialog title="Rename view" description="Give this view a short name for the workbench." size="sm">{#snippet trigger()}<span class="dialog-demo-trigger">Open dialog</span>{/snippet}<div class="dialog-demo-body"><Label for="dialog-name">View name</Label><Input id="dialog-name" value="Overview" /><div class="dialog-actions"><Button variant="quiet">Cancel</Button><Button variant="primary">Save view</Button></div></div></Dialog></div><Text size="sm" tone="muted">Tooltip, popover and dialog all own their keyboard and dismissal behavior through Bits UI.</Text></Surface></section>

			<section id="commands" class="gallery-section" aria-labelledby="commands-title"><div class="section-heading"><div><p class="eyebrow">Workbench</p><h2 id="commands-title">Command palette</h2></div><p>Search, keyboard navigation and execution stay separate from the command registry.</p></div><Surface tone="panel" padding="lg" class="command-demo"><div class="command-demo-copy"><strong>Quick actions</strong><Text size="sm" tone="muted">The palette is controlled by the host and reports the selected command through a callback.</Text></div><div class="command-demo-actions"><Button variant="secondary" onclick={() => (paletteOpen = true)}>Open command palette</Button><span class="kbd">⌘ K</span></div><Text size="sm" tone="quiet">Last command: {lastCommand}</Text></Surface><CommandPalette bind:open={paletteOpen} commands={paletteCommands} onExecute={runCommand} /></section>

			<section id="menus" class="gallery-section" aria-labelledby="menus-title"><div class="section-heading"><div><p class="eyebrow">Navigation + forms</p><h2 id="menus-title">Menus and choices</h2></div><p>Desktop actions, right-click context and preference choices share the same keyboard-first language.</p></div><Surface tone="panel" padding="lg" class="menus-demo"><div class="menu-row"><DropdownMenu>{#snippet trigger()}<span class="menu-demo-trigger">Actions</span>{/snippet}<MenuItem label="Open overview" shortcut="⌘ O" onSelect={() => chooseMenuAction('Open overview')} /><MenuItem label="Rename view" shortcut="F2" onSelect={() => chooseMenuAction('Rename view')} /><MenuSeparator /><MenuItem label="Delete view" disabled /></DropdownMenu><ContextMenu>{#snippet trigger()}<div class="context-demo-target">Right-click this canvas</div>{/snippet}<MenuItem label="Add panel" onSelect={() => chooseMenuAction('Add panel')} /><MenuItem label="Duplicate view" onSelect={() => chooseMenuAction('Duplicate view')} /><MenuSeparator /><MenuItem label="Paste" disabled /></ContextMenu><Text size="sm" tone="muted">Last action: {lastMenuAction}</Text></div><Separator /><div class="choices-grid"><div class="choice-demo"><Label>Display mode</Label><RadioGroup items={radioOptions} bind:value={menuRadio} ariaLabel="Display mode" /></div><div class="choice-demo"><Label>Layout</Label><ToggleGroup items={toggleOptions} bind:value={toggleValue} ariaLabel="Layout" /></div><div class="choice-demo"><Label>Compact labels</Label><Toggle bind:pressed={compactLabels}> {compactLabels ? 'On' : 'Off'} </Toggle></div></div><div class="shortcut-row"><Text size="sm" tone="muted">Shortcuts</Text><Kbd>⌘ K</Kbd><Kbd>⌘ B</Kbd><Kbd>F2</Kbd></div></Surface></section>

			<section id="collections" class="gallery-section" aria-labelledby="collections-title"><div class="section-heading"><div><p class="eyebrow">Collections</p><h2 id="collections-title">List and tree</h2></div><p>Selection is explicit, expansion is controlled, and scroll behavior belongs to the collection surface.</p></div><Surface tone="panel" padding="lg" class="collections-demo"><div class="collection-columns"><div class="collection-column"><Label>List</Label><List items={collectionItems} bind:value={listSelection} ariaLabel="Workspace views" /></div><div class="collection-column"><Label>Tree</Label><ScrollArea class="collection-scroll"><Tree items={treeNodes} bind:selectedId={treeSelection} bind:expandedIds={expandedTreeIds} ariaLabel="Workspace tree" /></ScrollArea></div></div><Text size="sm" tone="muted">Selected list item: {listSelection} · Selected tree item: {treeSelection}</Text></Surface></section>

			<section id="data-display" class="gallery-section" aria-labelledby="data-display-title"><div class="section-heading"><div><p class="eyebrow">Data display</p><h2 id="data-display-title">Inspection surfaces</h2></div><p>Cards frame related information while tables and property lists keep dense workbench data readable.</p></div><div class="data-display-grid"><Card title="Workspace metadata" description="A compact property surface for inspectors."><KeyValue items={workspaceProperties} columns={2} /></Card><Card title="Recent views" description="Rows stay presentational; selection and sorting belong to the feature."><Table columns={inspectorColumns} rows={inspectorRows} caption="Views in this workspace" compact /></Card></div></section>

			<section id="workbench-interactions" class="gallery-section" aria-labelledby="workbench-interactions-title"><div class="section-heading"><div><p class="eyebrow">Workbench interactions</p><h2 id="workbench-interactions-title">Adjustable structure</h2></div><p>Panel geometry and disclosure state stay keyboard-accessible and controlled by the consuming surface.</p></div><Surface tone="panel" padding="lg" class="interaction-demo"><Breadcrumb items={breadcrumbItems} /><div class="interaction-grid"><div class="interaction-column"><p class="field-label">Disclosure</p><Collapsible title="Inspector details" bind:open={inspectorOpen}><p>Keep secondary context available without adding another route or modal.</p></Collapsible><p class="field-label">Accordion</p><Accordion items={accordionItems} bind:value={accordionValue} ariaLabel="Token guidance" /></div><div class="interaction-column"><div class="interaction-label-row"><p class="field-label">Resizable split</p><Text size="sm" tone="quiet">{Math.round(splitSize)} / {Math.round(100 - splitSize)}</Text></div><div class="split-demo"><Resizable bind:value={splitSize}>{#snippet first()}<div class="split-pane"><strong>Explorer</strong><span>Drag the divider or focus it and use the arrow keys.</span><div class="split-placeholder"></div></div>{/snippet}{#snippet second()}<div class="split-pane"><strong>Canvas</strong><span>The second pane fills the remaining workbench space.</span><div class="split-placeholder split-placeholder-wide"></div></div>{/snippet}</Resizable></div></div></div></Surface></section>

			<section id="content-utilities" class="gallery-section" aria-labelledby="content-utilities-title"><div class="section-heading"><div><p class="eyebrow">Content utilities</p><h2 id="content-utilities-title">Identity and code</h2></div><p>Compact metadata, identity, and source snippets use the same quiet visual language as the workbench.</p></div><div class="content-demo-grid"><Card title="Workspace identity" description="Tags and avatars keep context close to the thing it describes."><div class="content-tag-row"><Avatar name="Atelier Core" status="online" size="lg" /><div class="content-identity-copy"><strong>Atelier Core</strong><span>Maintained by the platform team</span></div></div><div class="content-tag-row"><Tag tone="accent" dot>Design system</Tag><Tag tone="success" variant="outline">Synced</Tag><Tag tone="warning" dismissible>Review</Tag></div><p class="content-note">Tokens resolve through <InlineCode>var(--accent)</InlineCode><CopyButton value="var(--accent)" onCopy={showCopyResult} /></p><Text size="sm" tone="quiet">{copyFeedback}</Text></Card><Card title="Token source" description="Code stays readable and copyable without owning syntax highlighting."><CodeBlock code={tokenSnippet} language="tokens.css" onCopy={showCopyResult} /></Card></div></section>

			<section id="controls" class="gallery-section" aria-labelledby="controls-title">
				<div class="section-heading"><div><p class="eyebrow">Primitives</p><h2 id="controls-title">Control states</h2></div><p>Compact by default, with visible focus and honest disabled states.</p></div>
				<div class="token-card controls-card">
					<div class="control-row"><Button variant="primary">Primary action</Button><Button variant="secondary">Secondary</Button><Button variant="quiet">Quiet action</Button><Button variant="danger">Danger</Button><Button variant="primary" loading>Loading</Button><Button variant="secondary" disabled>Disabled</Button></div>
					<div class="control-row"><label class="field-label" for="sample-input">Input</label><input id="sample-input" class="field" value="Workbench value" /><select class="field" aria-label="Example select"><option>Example select</option><option>Another option</option></select><IconButton label="Open command palette" size="sm">⌘</IconButton><span class="kbd">⌘ K</span></div>
				</div>
			</section>

			<section id="workbench" class="gallery-section" aria-labelledby="workbench-title">
				<div class="section-heading"><div><p class="eyebrow">Composition</p><h2 id="workbench-title">Workbench shell</h2></div><p>Nested surfaces, persistent chrome, and panel-level scrolling.</p></div>
				<WorkbenchShell title="Atelier workspace" class="workbench-demo">
					{#snippet activity()}<div class="workbench-activity-buttons"><IconButton label="Overview" size="md" class="rail-icon active">✦</IconButton><IconButton label="Files" size="md" class="rail-icon">□</IconButton><IconButton label="Settings" size="md" class="rail-icon">◌</IconButton></div>{/snippet}
					{#snippet sidebar()}<WorkbenchPanel title="Explorer"><ScrollArea class="workbench-sidebar-scroll"><Tree items={treeNodes} bind:selectedId={treeSelection} bind:expandedIds={expandedTreeIds} ariaLabel="Explorer tree" /></ScrollArea></WorkbenchPanel>{/snippet}
					{#snippet main()}<div class="workbench-main-view"><ViewHost views={workbenchViews} bind:activeView={workbenchView}><TabPanel value="overview"><div class="workbench-view-content"><span class="preview-kicker">ACTIVE VIEW</span><strong>Workbench canvas</strong><p>This surface is ready for registered views.</p><div class="preview-blocks"><span></span><span></span><span></span></div></div></TabPanel><TabPanel value="tokens"><div class="workbench-view-content"><span class="preview-kicker">INSPECTOR</span><strong>Token inspector</strong><p>Registered views can own their content while the shell stays stable.</p></div></TabPanel></ViewHost></div>{/snippet}
					{#snippet status()}<Status tone="success">Ready</Status><span>main · {density}</span>{/snippet}
				</WorkbenchShell>
			</section>

			<section class="scale-section" aria-labelledby="scale-title"><div class="section-heading"><div><p class="eyebrow">Scale</p><h2 id="scale-title">Spacing rhythm</h2></div><p>Values stay small and predictable for dense desktop layouts.</p></div><div class="scale-list">{#each spacing as token}<div class="scale-row"><code>{token}</code><span class="scale-bar" style={`inline-size: var(${token})`}></span></div>{/each}</div></section>
		</main>
	</div>
</div>
