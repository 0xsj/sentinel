# Styles

Layered tokens and base styles are composed by index.css. The token stack covers
palette primitives, semantic surfaces/content/status, typography, spacing, shape,
motion, z-index and workbench geometry. Theme selection uses `data-theme="light"`
or `data-theme="dark"`; leaving it unset follows the system preference. Density
uses `data-density="compact"` or `data-density="comfortable"`.

Cascade order: reset → token → base → primitive → composition → screen → override.
Component styles must use their owning layer when introduced. Bits UI supplies
headless interaction behavior; Atelier owns the visual wrappers and CSS.
