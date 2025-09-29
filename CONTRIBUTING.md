# Versioning

This project follows [Semantic Versioning 2.0.0](https://semver.org/).

## Version Format
Releases use the format:


- **MAJOR** — Incremented for incompatible changes or breaking changes.  
- **MINOR** — Incremented for new functionality that is backward-compatible.  
- **PATCH** — Incremented for backward-compatible bug fixes.

Examples:
- `v0.1.0` – Initial pre-release with limited functionality.  
- `v0.1.1` – Patch release fixing a bug in the `v0.1.x` series.  
- `v1.0.0` – First stable release
- `v1.1.0` – Backward-compatible feature added to the `v1.x.x` series.  
- `v1.1.1` – Patch release fixing a bug introduced in `v1.1.0`.  

## Lifecycle

- **Pre-Release (`v0.x.x`)**  
  Versions below `1.0.0` are considered unstable. The cli may change at any time.

- **Stable Release (`v1.0.0` and above)**  
  Once the project reaches `v1.0.0`, the cli is considered stable. Breaking changes will only occur with a **MAJOR** version bump.

- **Patch Releases (`x.x.PATCH`)**  
  Safe upgrades that fix bugs without changing existing functionality.

- **Minor Releases (`x.MINOR.0`)**  
  Backward-compatible features, improvements, or enhancements.

- **Major Releases (`MAJOR.0.0`)**  
  May include breaking changes, removed features, or significant redesigns.