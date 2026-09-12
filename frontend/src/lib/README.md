# Frontend library boundaries

Dependencies point inward:

```text
root → features → workbench → components → kernel
              ↘ services → consumer ports + kernel
root → platform adapters → consumer ports + kernel
UI → runtime bindings → portable models → kernel
styles → tokens and ordered cascade layers
```

Root constructs adapters and injects capabilities. Features do not construct them.
Workbench receives contributions from root and never imports features. Components
receive values/callbacks. Pure models do not import Svelte. Native SDKs stay in
platform/desktop; future codecs accept unknown rather than trust a cast.

See FOUNDATION.md at the project root for the library build order. README-only
folders reserve ownership; they are not exports, runtime modules or functionality.
No automated import enforcement is installed yet; these boundaries are reviewed.
