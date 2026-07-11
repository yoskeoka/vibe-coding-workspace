# Memory boundary and audit

## Summary

workspace memory は便利だが、何がどれだけ入っているかが会話から見えにくく、repo-owned な durable knowledge と一時的な rollout note が混ざりやすい。

今回の token 消費分析でも、memory の存在自体が主犯ではないにせよ、「どこまでを memory に置き、どこからを repo の docs/specs/skills に残すべきか」が曖昧だと、不要な検索や要約読みに繋がる懸念が見えた。

## Problem

- durable に残すべき知識が memory に寄りすぎると、repo だけ見ても運用意図が追えなくなる
- memory 参照時は agent から投入量が見えづらく、読まれた量や粒度を人間が把握しにくい
- rollout 断片や過去の細かい状況メモが残り続けると、再利用価値より探索コストが上回る可能性がある

## Desired boundary

- repo の設計・仕様・ workflow 契約は `docs/` と `skills/` に置く
- memory は「再発しやすい運用上の罠」「repo 横断の薄いポインタ」「当座の探索起点」に寄せる
- memory がなくても repo だけで主要な workflow 判断が追える状態を保つ

## Follow-up direction

- 現在の memory のうち、repo-owned knowledge を持ちすぎている領域を棚卸しする
- `AGENTS.md` / `AI_WORKFLOW.md` / relevant skills に、memory は補助であり source of truth ではないことをもう少し明示するか検討する
- repo に移すべき知識と、削ってよい rollout 運用メモの基準を決める

## Status

repo-owned workflow 契約との境界整理が必要。対処方針は見えているが、棚卸し対象と移管基準の明文化が先。
