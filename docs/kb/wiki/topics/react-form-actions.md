---
title: React Form Actions
last_reviewed: 2026-11-28
status: seed
sources:
  - ../../sources/2026/2026-06-05-brsf-react-useactionstate-forms.md
---

# React Form Actions

## Current signal

React 19 form Actions can submit named, uncontrolled fields as `FormData` without maintaining a `value` / `onChange` pair for every input. `useActionState` adds a returned state value and an action-pending flag for workflows that need validation messages or status UI.

## Practical guidance

- Use a function in the form's `action` prop when the standard Action behavior fits; React manages the submission transition and uncontrolled form reset.
- Use `useActionState` when the form needs to display the action's returned state. Its `isPending` value means the action is in progress, not that it has completed.
- Return a new state from the action. Do not mutate the previous state object or omit the return from a wrapper action.
- Keep the Qiita article as a Japanese introduction, but check code against the official React references linked in its source note.
