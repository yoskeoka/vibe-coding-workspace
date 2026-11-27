---
title: "React forms with useActionState"
source_url: https://qiita.com/BRSF/items/206c141f06cd94b724c6
source_type: article
original_language: ja
ingested_on: 2026-11-28
status: active
tags:
  - react
  - forms
  - javascript
related_pages:
  - ../../wiki/topics/react-form-actions.md
---

# React forms with useActionState

## Summary

BRSF's Qiita article introduces React 19 form Actions and `useActionState` as a way to submit uncontrolled fields through `FormData`, without writing a `value` / `onChange` pair for every field. Its motivating pattern is useful: let the browser own ordinary field editing, submit named fields, and keep returned action state for validation or follow-up UI.

## Verification notes

The article's examples should not be copied verbatim. React's official reference defines the third `useActionState` result, `isPending`, as whether the action is currently in progress; it is not a post-submit completion signal. The official API expects the action to return the next state. One later example in the article computes `updated` but does not return it from the wrapper action, so the hook state would not become that result. The earlier examples also mutate the previous state object directly. Prefer returning a new state value, handle validation results as returned state, and use the documented form Action / pending-state behavior.

## Workspace relevance

Keep this as a Japanese-language entry point for the uncontrolled-form workflow, paired with React's primary documentation. The article is most useful for its motivation and form-data examples; use the official reference for exact hook semantics and production implementation.

## Sources

- Qiita article: https://qiita.com/BRSF/items/206c141f06cd94b724c6
- React `useActionState`: https://react.dev/reference/react/useActionState
- React `<form>` Actions: https://react.dev/reference/react-dom/components/form
