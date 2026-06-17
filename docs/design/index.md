---
title: Design documents
summary: Architecture and improvement designs for Conflow runtime and tooling.
keywords: [design, architecture, plans]
---

# Design documents

Design documents capture **problem analysis, decisions, and batch roadmaps** before implementation. Each initiative lives in its own subdirectory with an `index.md` entry point; drill down from there into focused leaf documents.

Implementation plans derived from these designs are written separately in [`docs/plans/`](../plans/).

## Initiatives

| Initiative | Status | Summary |
|------------|--------|---------|
| [Workflow engine improvements](./workflow-engine-improvements/index.md) | Draft | Audit of the runtime orchestration layer (`pkg/conflow/block`, `job`, pub/sub) with prioritized improvement areas and batch roadmap |

## Conventions

- **Design docs** (`docs/design/`) — what to build and why; decisions, alternatives, scope, batch breakdown.
- **Plan docs** (`docs/plans/`) — how to build one batch; step-by-step tasks with file paths and verification commands.
- **Product docs** (`docs/product/`) — user-facing behaviour and API reference.

When implementing a batch, read the relevant design doc first, then the batch plan in `docs/plans/`.
