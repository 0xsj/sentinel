# Frontend composition root

App.svelte is the working Hello world composition. This layer owns app lifetime,
provider wiring and feature registration as those capabilities are introduced.
Nothing below root imports it.
