# Workflow context read budget

## Summary

今回の大きな token 消費は、実装コードそのものよりも、長い workflow 文書や skill 文書、PR template、follow-up helper 出力などを同一 session で何度も扱ったことによる部分が大きかった。

`gh-pr-followup` 自体の compact 化はすでに進んでいるが、routine な workflow 実行でもまだ「長い文書を丸ごと読む」「必要節だけで足りる場面でも全文を読む」余地が残っている。

## Problem

- `AGENTS.md`、`AI_WORKFLOW.md`、skill docs、PR template などが大きい
- 実作業では必要な節が限られていても、運用上は全文読みに寄りやすい
- compact helper を導入した後でも、周辺の文脈取得が大きいと総 input token は膨らむ

## Candidate improvements

- 長文 workflow / skill / template は「必要節だけ読む」契約を明文化する
- PR template は全文確認ではなく、埋める節の最小参照に寄せる
- `review-task` と related skills で、普通の landing loop では verbose 情報に上がらない判断基準をさらに明示する
- handoff / final status で file-by-file recap を抑え、未完了事項中心に寄せる

## Why this is repo-owned

これは platform の hidden prompt ではなく、workspace が管理している workflow docs / skills / helpers の使い方の問題なので、repo 側で改善可能。

## Status

対処方向はかなり明確。必要ならこの issue から exec-plan を起こして、workflow docs / skills / template の compact-read contract を揃える。
