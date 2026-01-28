# Repository Guidelines

## Project Structure & Module Organization
This repository is currently empty (no source tree, tests, or assets checked in yet). As you add code, keep a predictable layout so contributors can orient quickly. Recommended baseline:

- `src/` for application or library code.
- `tests/` for automated tests.
- `assets/` for static files (images, fixtures, sample data).
- `docs/` for longer-form documentation.

If you choose a different structure, document it here with concrete paths and responsibilities.

## Build, Test, and Development Commands
There are no build or test tools configured yet. Once you introduce tooling, list the canonical commands and what they do. Example format:

- `npm run dev`: start the local development server.
- `npm test`: run the full test suite.
- `make build`: produce a production build.

Keep this section authoritative and minimal; avoid listing one-off scripts.

## Coding Style & Naming Conventions
No formatting or linting rules are defined yet. When you add them, document:

- Indentation and line length (e.g., 2 vs 4 spaces, 100–120 columns).
- Language-specific style guides (e.g., ESLint/Prettier, Black, gofmt).
- Naming patterns (e.g., `PascalCase` for types, `camelCase` for functions, `kebab-case` for file names).

## Testing Guidelines
No test framework is configured. When you add tests, specify:

- Frameworks (e.g., Jest, pytest, Go test).
- Naming patterns (e.g., `*.test.ts`, `test_*.py`).
- Coverage targets or required gates, if any.

## Commit & Pull Request Guidelines
There is no Git history or commit convention established yet. Once you have history, summarize the preferred pattern (e.g., Conventional Commits or short imperative subjects).

For pull requests, at minimum include:

- A clear description of the change and motivation.
- Links to relevant issues or tickets.
- Screenshots or recordings for UI changes.
- Notes on testing performed.

## Security & Configuration Tips
If the project requires secrets or environment variables, document them in `.env.example` (never commit real secrets). Note any setup steps contributors must run before development.
