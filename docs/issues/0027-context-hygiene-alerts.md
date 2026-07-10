# Context hygiene alerts

## Summary

context 掃除を「重くなったらその場の勘でやる」運用にせず、次回は alert を上げられるチェック項目として持ちたい。

理想は、repo-owned な長文文書や skill / template / helper surface を対象にした軽量な script を用意し、token を直接読めない場面でも「この session は重くなりやすい状態か」を事前に察知できること。

## Why this is needed

- session token usage は後から気づきやすく、事前の guard になりにくい
- hidden prompt や platform-owned context は直接測れないことがある
- それでも repo 側の長文 surface は測れるので、先にそこへ alert を出すだけでも実用価値がある

## Candidate checks

### 1. Long workflow document size

- `AGENTS.md`
- `AI_WORKFLOW.md`
- `.github/PULL_REQUEST_TEMPLATE.md`
- long-lived skill docs under `skills/*/SKILL.md`

Alert examples:

- warn if a single workflow doc exceeds a chosen line or byte budget
- warn if the combined size of the common startup set exceeds a session budget

### 2. Repeated heavy-read surfaces

- same PR template or long skill doc being required in ordinary flows
- helper output that often escalates to verbose mode

Alert examples:

- warn when a helper or workflow contract has both compact and verbose modes but the compact path is still too large
- warn when routine startup docs have overlapping or duplicated sections

### 3. Memory boundary pressure

- memory registry size
- number of rollout-derived notes that duplicate repo-owned workflow knowledge

Alert examples:

- warn when memory grows but the repo-visible source of truth does not reflect the same rule
- warn when a memory summary item points at workflow behavior that should already live in `docs/` or `skills/`

## Initial threshold proposal

These are starting points, not final rules:

- `AGENTS.md` or `AI_WORKFLOW.md`: warn at roughly 250-300 lines each
- PR template: warn at roughly 120 lines
- single skill doc: warn at roughly 200 lines
- combined common startup set (`AGENTS.md`, `AI_WORKFLOW.md`, PR template, most-used skill docs): warn when it crosses a budget such as 800-1000 lines
- explicit session cleanup recommendation when a medium task also needs repeated PR follow-up and at least two of the above are already over budget

## Script idea

Add a lightweight checker such as:

```sh
tools/context-hygiene-check.sh
```

The script should:

- print per-file line and byte counts for the common workflow surfaces
- show which items exceed configured thresholds
- print a simple overall status such as `ok`, `warn`, or `cleanup-recommended`
- optionally inspect configured memory files when available, but degrade gracefully when local memory is absent

## Status

plan-ready once the workspace decides the first threshold set and whether memory inspection belongs in the initial version or a follow-up.
