---
title: React Form Actions
last_reviewed: 2026-11-28
status: seed
sources:
  - ../../sources/2026/2026-06-05-brsf-react-useactionstate-forms.md
---

# React Form Actions

## 現在の動向

React 19 の form Actions は、各 input に `value` / `onChange` を設定せず、`name` 付きの uncontrolled field を `FormData` として送信できます。validation message や状態表示が必要な場合、`useActionState` で action の返却 state と実行中 flag を扱えます。

## 実装の指針

- 標準 Action の挙動が合う場合は、form の `action` に関数を指定します。React が transition と uncontrolled form の reset を管理します。
- action の返却値を表示する場合に `useActionState` を使います。`isPending` は action 実行中を示し、完了を知らせる値ではありません。
- action は新しい state を返します。以前の state object を mutate したり、wrapper action で return を省略したりしないでください。
- Qiita 記事は日本語の導入として利用し、コードは source note にある React 公式リファレンスと照合します。
