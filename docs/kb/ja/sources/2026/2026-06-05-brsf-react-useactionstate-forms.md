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

## 要約

BRSF 氏の Qiita 記事は、React 19 の form Actions と `useActionState` を使い、各入力に `value` / `onChange` を書かず、`FormData` で uncontrolled form を送信する方法を紹介しています。通常の入力編集は browser に任せ、`name` 付きフィールドを送信し、返された action state を validation や後続 UI に使う、という動機は有用です。

## 確認時の注意

記事中のコードはそのまま使わないでください。React公式リファレンスでは `useActionState` の第3戻り値 `isPending` は action の実行中を示し、submit 完了を検知する値ではありません。action は次の state を return する必要があります。記事後半の例は `updated` を計算しているものの、wrapper action から返していないため、その値が hook state にはなりません。前半には前 state object を直接 mutate する例もあります。新しい state を return し、validation 結果は返却 state として扱い、form Action と pending state は公式仕様に沿って使います。

## ワークスペースでの意義

uncontrolled form workflow の日本語記事への入口として保存し、React公式ドキュメントと併記します。記事は動機と FormData の例に有用ですが、hook の正確な動作と実装は公式リファレンスを優先します。

## 出典

- Qiita記事: https://qiita.com/BRSF/items/206c141f06cd94b724c6
- React `useActionState`: https://react.dev/reference/react/useActionState
- React `<form>` Actions: https://react.dev/reference/react-dom/components/form
