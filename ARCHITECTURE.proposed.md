# Proposed Atelier structure

Update: Svelte 5 is selected and the initial directory scaffold now exists.
[FOUNDATION.md](FOUNDATION.md) describes the current build order; reusable frontend
modules now live under `src/lib` (within `frontend/` for Wails). The original trees
below remain design history, not an exact inventory.

Draft for discussion · 2026-09-11. This is a target organization, not an
implementation inventory or a commitment to build every listed module. The
running app remains the Hello world scaffold. Names describe ownership; create
a package only when an actual slice exercises its behavior.

The starting shape is an independently usable desktop modular monolith. Carry
forward Overwatch's pure domains, consumer-owned ports, reusable leaves, and
explicit composition. Atelier defines its own workbench vocabulary and does not
import sibling projects. Go and Rust should support equivalent behavior through
idiomatic APIs; they need not mirror every file.

## Native tree

```text
atelier-wails/
├── main.go                         # thin Wails entry; embeds frontend assets
├── go.mod / go.sum
├── wails.json
├── Makefile
├── root/                           # process composition; nothing imports root
│   ├── build.go                    # construct services and concrete adapters
│   ├── lifecycle.go                # startup, admission, cancellation, shutdown
│   └── workflows/                  # explicit cross-context orchestration
├── internal/
│   ├── host/                       # Wails and OS boundary
│   │   ├── wails/                  # bound facades, runtime context, event bridge
│   │   ├── windows/                # window lifecycle and close coordination
│   │   ├── menus/                  # native menus mapped to application actions
│   │   ├── dialogs/                # native open/save/confirmation dialogs
│   │   ├── clipboard/              # platform clipboard adapter
│   │   ├── paths/                  # platform app-data/cache locations
│   │   └── lifecycle/              # OS quit/suspend/resume integration
│   ├── workspace/                  # candidate: workspace identity/open lifecycle
│   │   ├── domain/                 # values, invariants, transitions; no effects
│   │   ├── app/
│   │   │   ├── command/            # open, rename, close; consumer-owned ports
│   │   │   └── query/              # list, inspect; consumer-owned read ports
│   │   ├── infra/
│   │   │   └── persistence/        # records/codecs/migrations for this context
│   │   └── transport/
│   │       └── desktop/            # plain DTOs and mapping; no Wails imports
│   ├── preferences/                # candidate: scoped, validated user settings
│   │   └── ...                     # same layers, only where behavior requires
│   └── jobs/                       # later candidate: tracked background work
│       └── ...                     # state transitions, progress, cancellation
├── pkg/                            # domain-free; reusable outside Atelier
│   ├── errors/                     # stable failure meanings; no UI wording
│   ├── id/                         # identity values; generation injected
│   ├── clock/                      # time capability; deterministic substitutes
│   ├── secret/                     # explicit redaction at diagnostic boundaries
│   ├── logger/                     # diagnostic logging, owned sinks
│   ├── config/                     # process configuration, not user preferences
│   └── fileio/                     # reusable file primitives when consumed
├── frontend/                       # frontend tree below
│   ├── src/
│   └── wailsjs/                    # generated bindings; platform adapter only
├── build/                          # native packaging assets
├── tests/
│   ├── integration/                # real adapters and desktop boundary
│   └── desktop/                    # launch, interaction, shutdown scenarios
├── tools/
│   └── architecture/              # import/dependency checks when introduced
├── decisions/                      # accepted commitments and their tradeoffs
├── notes/                          # findings and verification evidence
├── ARCHITECTURE.proposed.md
└── README.md                       # actual inventory and run commands
```

Keep the existing root `main.go` for Wails and Go asset embedding. Pass embedded
assets into root; root configures the host and wires the facades. Host facades
receive operation interfaces, never the composition root. `context.Context`
belongs in effectful application APIs where cancellation is needed, not in pure
rule evaluation. Interfaces live with the package that consumes them.

Start with one Go module. `pkg/` means deliberately reusable code, not a promise
of separately versioned distribution. Go's `internal/` visibility does not
prevent peer contexts from importing each other; that rule needs a check.

## Frontend tree

The same conceptual organization lives in each project independently. In Wails
this is `frontend/src/`; in Tauri it is `src/`. File extensions and component
implementation await the UI framework decision.

```text
<frontend source>/
├── main.ts                         # existing frontend entry
├── root/                           # frontend composition; adapter selection
│   ├── bootstrap/                  # initialize/dispose app and subscriptions
│   ├── providers/                  # framework-level dependency wiring
│   └── contributions/              # explicitly register features with workbench
├── styles/
│   ├── tokens/                     # semantic color, space, type, motion, density
│   ├── themes/                     # light/dark/system; persistence via preferences
│   └── global/                     # reset, base styles, layer ordering
├── lib/                            # framework-independent reusable leaves
│   ├── result/                     # explicit expected outcomes and failures
│   ├── identity/                   # typed frontend identifiers
│   ├── disposal/                   # owned subscriptions and cleanup
│   ├── collections/                # selection/order helpers with real consumers
│   └── keybinding/                 # shortcut parsing and matching primitives
├── platform/                       # only frontend importer of native bindings
│   ├── contracts/                  # narrow capabilities needed by UI consumers
│   ├── desktop/                    # Wails OR Tauri implementation in this repo
│   ├── codecs/                     # validate/map native payloads and failures
│   └── preview/                    # explicit browser/demo adapters
├── components/                     # reusable UI; no native or feature imports
│   ├── primitives/                 # button, icon, text, surface, separator
│   ├── layout/                     # stack, inline, grid, scroll area, splitter
│   ├── forms/                      # field, input, textarea, select, checkbox,
│   │                               # radio, switch, combobox, validation messages
│   ├── navigation/                 # tabs, breadcrumbs, menus, toolbar
│   ├── overlays/                   # dialog, popover, tooltip, context menu
│   ├── collections/                # list, tree, table; keyboard/selection behavior
│   ├── feedback/                   # empty/loading/error, progress, banner, toast
│   └── patterns/                   # search field, settings row, confirmation,
│                                   # editable label, master-detail composition
├── workbench/                      # reusable desktop interaction foundation
│   ├── model/
│   │   ├── layout/                 # regions, splits, bounds; serializable values
│   │   ├── views/                  # view instances, activation, open/close state
│   │   ├── selection/              # active selection and scope
│   │   └── commands/               # action IDs, availability, execution contracts
│   ├── app/
│   │   ├── commands/               # dispatch actions from menus/palette/shortcuts
│   │   ├── keybindings/            # precedence, scopes, conflict handling
│   │   ├── navigation/            # active view and navigation history
│   │   └── restoration/           # snapshot, validate, restore through a port
│   └── ui/
│       ├── shell/                  # title area, sidebar, main area, status area
│       ├── panels/                 # tool regions and resizing
│       ├── view-host/              # tabs and rendering registered views
│       ├── command-palette/        # searches and invokes registered commands
│       └── status/                 # active work and notifications
├── features/                       # user-facing slices; own view models and UI
│   ├── workspace/                  # workspace picker/recent/open feedback
│   ├── preferences/                # settings screens and effective-value display
│   └── jobs/                       # later: progress/cancel/history surface
├── dev/
│   ├── gallery/                    # components, states, themes, density examples
│   └── fixtures/                   # deterministic scenarios through real ports
└── testing/                        # shared test helpers; tests beside owners
```

Dependencies point from composition and features toward workbench, components,
and leaves. `workbench/` does not import `features/`; root registers explicit
contributions. Components do not import the workbench. Platform implementations
satisfy consumer capabilities and are injected by root; models never call a
native SDK. A feature owns its use-case contracts; `platform/contracts/` is for
common desktop capabilities, not every operation in the application.

Workbench models describe interaction rules and remain free of rendering and
native APIs. Focus mechanics, DOM measurement and accessibility wiring belong
with UI behavior. Native domains own durable application rules; frontend models
own transient interaction. Avoid duplicating authoritative native rules in UI.

## Ownership that should stay explicit

| Concern | Owner and boundary |
| --- | --- |
| Workspace | Native candidate context: registry metadata and open/close rules. Revision 1 semantics are in [DOMAIN_CONTRACTS.md](DOMAIN_CONTRACTS.md); implementation adds only the required path capability. |
| Preferences | Native candidate context: definitions, validation, scope and persistence. Revision 1 semantics are in [DOMAIN_CONTRACTS.md](DOMAIN_CONTRACTS.md); frontend displays effective values. Process config stays separate. |
| Layout and views | Frontend workbench owns splits, active views and restoration format. Persist through a narrow snapshot port; this alone does not need a native layout domain. |
| Commands | Workbench actions unify palette, shortcut and menu invocation. Native application commands own durable changes. These are different APIs with explicit mapping. |
| Jobs | Introduce the specified native context when work needs tracked identity/lifecycle. A goroutine, future, or spinner alone is not a job domain. |
| Files | Host handles user selection; a file adapter performs IO; the consuming application owns interpretation and save/conflict policy. No universal file domain by default. |
| Undo/redo | Local text/selection history stays local. Durable undo requires an owning use case and conflict policy; never infer it from a generic command interface. |
| Windows | Native host owns real windows. Frontend owns view state within a window; root coordinates close with pending work. |
| Diagnostics | Logging/configuration are process foundations. Error display belongs to UI, preserving actionable failure meanings. |

Documents/editors, search indexes, extensions, agents, terminals, sync, accounts,
and update delivery are later candidates driven by a real product need. Tabs
host views and do not force a universal document model. A contribution registry
can begin as ordinary registration calls; a dynamic plugin loader is a separate
capability with its own lifecycle and access model.

## Dependency rules

1. Pure domains import only the standard library and selected pure value leaves.
   They do not read time, generate IDs, access files, or perform logging themselves.
2. Application operations depend on their own domain and small consumer-owned
   ports. Effects, cancellation and consistency are explicit at this boundary.
3. Infrastructure implements ports; desktop transport maps inputs and outcomes.
   Native framework imports are restricted to host and root. DTOs, database
   records and domain values have separate owners.
4. Contexts never import peer contexts. Root supplies translation adapters or
   explicit workflows through application APIs. A shared package cannot import
   a context, host, or root.
5. Shared is not synonymous with pure: logger and file IO are effectful leaves.
   They are available to adapters/root, not to pure domains. Keep value packages
   independent from their effectful implementations.
6. Frontend/native messages have explicit shapes and stable failure meanings.
   Subscription disposal, progress ordering, window lifetime, cancellation and
   resynchronization must be specified when events arrive. Cancellation does not
   prove that a write did not commit.
7. Root owns startup and shutdown. Persistence transactions belong to the use
   case promising a change. Multi-context workflows state partial-failure and
   recovery behavior; a root function alone does not provide atomicity.

## UI quality expected from the foundation

Components should support keyboard interaction, visible focus, accessible names,
focus restoration, appropriate announcements, disabled/read-only states, and
reduced motion. Collection controls need deliberate selection and navigation
semantics. Splits and panels need minimum sizes, overflow handling and a keyboard
resize path. Scoped shortcuts must respect text editing and composition input.

Keep loading, empty, failed, stale and successful states distinguishable where
they matter. Include these in the component gallery along with light/dark themes,
density, long labels and constrained window sizes. Native dialogs and custom
prompts should return explicit cancel/accept outcomes.

Use a maintained accessible primitive layer once the framework is selected;
wrap it only where Atelier adds a stable styling or behavior contract. Framework,
primitive library, styling implementation and persistence backend remain open.

## Build order: leaves into complete slices

1. **Running host — current baseline.** Hello world, native launch and shutdown.
2. **First UI slice.** Tokens, text, button, surface and layout; a small gallery
   that also runs in the desktop host. Establish focus and theme behavior here.
3. **First boundary slice.** Read/change one preference through UI → typed bridge
   → application → validated value → persistence adapter. Add only the leaves
   this needs; prove failure and restart behavior.
4. **Workbench slice.** Shell with two registered views, activation and one action
   reachable through a menu and shortcut. Grow layout and restoration from this.
5. **Workspace slice.** Implement the revision 1 registry contract, then add
   path inspection only to a concrete open workflow.
6. **Long-running slice.** Implement Jobs only with an actual operation requiring
   progress, cancellation and shutdown behavior.

Leaf-first describes dependency order within a useful slice. It does not require
building every generic package before showing a working feature.

Pure rules get deterministic tests beside their code. Adapters get integration
checks for real failure boundaries; the native bridge gets serialization and
outcome checks. UI checks cover keyboard/focus behavior in a real renderer and
native smoke scenarios cover launch, interaction and quit. Architecture checks
should enforce forbidden imports as the first real boundaries appear.

This draft was checked against the existing Atelier scaffolds, Overwatch's
`pkg/`, `internal/<context>/{domain,app,infra,transport}` and `root/` organization,
and the local n2f Go/Rust and Flover UI conventions. No code was reorganized and
no runtime checks were needed for this documentation-only proposal.
