# Model-switch summary bloat

## Summary

長い実装セッションのあとに model を切り替えると、直前までの作業内容がかなり詳細な summary として新モデル側に引き継がれることがある。

今回の `ww` 実装でも、plan 通りの比較的直線的な変更に対して input token が不自然に大きく、その主要因の 1 つがこの summary だった。

## Why this matters

- 実装本体よりも「前セッションの詳細な作業メモ」を持ち越してしまう
- summary が file-by-file の変更一覧、検証履歴、PR 状態、未完了作業まで広く含むと、その後の軽い分析や follow-up でも高コストになる
- これは exec-plan や repo の spec とは別物で、platform / orchestrator 側の context carry-over として扱うべき問題

## Current understanding

- ここでいう summary は `docs/exec-plan` ではない
- 会話の model switch / compaction 時に内部的に生成される引き継ぎ要約に近い
- repo 側で完全には制御できないが、repo-owned workflow の handoff ルールを compact にすることで間接的に悪化を抑えられる可能性はある

## Follow-up direction

- platform-owned な summary と repo-owned な handoff / PR follow-up artifact を明確に分離して考える
- 今後同様の session で、summary が必要以上に詳細なときの典型パターンを集める
- repo 側で対処できる範囲としては、hand-off 時に必要最小限の未完了事項だけを残す運用へ寄せる

## Status

調査継続。repo ローカルだけでは完結しないため、まずは境界を明確にする issue として残す。
