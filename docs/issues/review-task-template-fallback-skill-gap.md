# review-task template fallback skill gap

## Summary

`skills/review-task/SKILL.md` の PR template fallback 説明が、workspace 配下 child repo の実運用を十分に表していない。

今回 `ai-arena` で child repo 直下に `PULL_REQUEST_TEMPLATE.md` が見つからなかった際、本来は
workspace root の `vibe-coding-workspace/.github/PULL_REQUEST_TEMPLATE.md` を参照すべきだったが、
skill の記述だけではその fallback が弱く、repo 内 template 不在をそのまま手書き PR body 作成へ倒しやすかった。

影響範囲:

- `skills/review-task/SKILL.md`
- 必要なら `skills/review-task/scripts/gh-pr-followup` など周辺 helper の fallback 説明
- workspace 配下 child repo の PR 作成フロー全般

## Proposed Solution

- `skills/review-task/SKILL.md` の template fallback 順を明示する
  - child repo local template: `<child>/.github/PULL_REQUEST_TEMPLATE.md`
  - workspace root template: `vibe-coding-workspace/.github/PULL_REQUEST_TEMPLATE.md`
  - vendored workflow template / workflow repo template はその後段に置く
- `workspace-level` のような曖昧表現を避け、具体 path を書く
- helper script 側に同等の fallback 実装があるなら、skill 文言と実装順を一致させる

## Priority

中。product code の不具合ではないが、workspace 配下 repo で PR body 品質と workflow 一貫性を落としやすい。
