---
title: AI 補助ゲーム開発
last_reviewed: 2026-04-25
status: seed
sources:
  - ../../sources/2026/2026-03-10-munimaru-component-oriented-gameplay.md
  - ../../sources/2026/2026-03-18-classmethod-godot-codex-three-themes.md
  - ../../sources/2026/2026-04-11-openai-gpt-5-4-phaser-tactical-rpg.md
  - ../../sources/2026/2026-04-11-shi3z-vibe-coding-vertical-shooter.md
  - ../../sources/2026/2026-04-18-phaser-claude-code-space-shooter.md
  - ../../sources/2026/2026-04-22-unity-interactive-writing-esoteric-ebb.md
  - ../../sources/2026/2026-04-10-sonicmoov-next2d-origin-story.md
  - ../../sources/2026/2026-04-13-gamespark-luminary-family-coop-arpg.md
---

# AI 補助ゲーム開発

## 現在のシグナル

最近の事例は、AI agent が gameplay system、browser playtest、asset generation、engine-specific tooling をひとつの loop で補助する流れを示している。

`Esoteric Ebb` の振り返りは、agent 的ではないが重要な production signal を足している。野心的なゲーム content は、authoring tool、state naming、search、recompile、validation が反復可能な速度を保てるときにだけ規模化できる。

## このワークスペースで重要な理由

- ゲーム project はワークスペース構成の大きな部分を占める。
- agent-driven な prototyping は、人間の可処分時間が限られる hobby project を前に進めやすくする。
- Phaser tactical RPG の参照は内部実装が薄いが、`code generation -> browser playtest -> UI/combat iteration -> asset generation` という loop の proof point としては十分役立つ。
- Peter Yang の space-shooter tutorial は、より手順的なパターンを足している。agent に design question をさせ、具体 spec を作り、MVP を作ってから feature milestone を積み、最後に deploy する。
- Godot / Codex 比較は、より厳密な評価軸を与える。初手の scaffold は強くても、movement feel、sync integrity、editor usability は人間の確認が必要である。
- shi3z の shooting prototype は、AI が開発支援だけでなく game loop そのものに入る具体例で、speech recognition、LLM interpretation、generated art、generated attack behavior が live play に影響する。
- `Next2D` の起源記事は、rendering architecture、production tooling、engine-specific `MCP` server が一体の stack である点で relevant。
- `Luminary` interview はよい対照になる。AI が開発に入っても、よい game-design の参照は別途必要であり、shared solo / co-op progression や punishment の削減は残すべき具体アンカーである。
- branching narrative system は固有の spec surface を要する。dialogue-heavy な game では、authoring tool、state-variable convention、feedback point、validation check を、content が膨らむ前に記録しておくべき。

## 再利用できるパターン

- まず vertical shooter のように mechanic が明快な小さな genre から始める。
- 実装前に asset を先に選び、agent が player、enemy、projectile、power-up、background の file 名に直接対応づけられるようにする。
- spec を書く前に、agent に clarifying question をさせる。
- 最初の build は playable MVP とみなし、その後 enemy wave、power-up、boss health bar、screen shake、deployment のような細い pass で機能を足す。
- agent 性能比較では、theme ごとに評価 session を分ける。action feel、multiplayer sync、editor tooling は別種の failure mode を露出する。
- deployment は早めに loop に入れる。Vercel に載せた static Phaser build 程度で共有と review には十分である。
- system-heavy な game では、設計が fixed script chain ではなく emergent interaction に依存するなら、architectural safety net と intent protocol を先に決める。
- accessibility-first な RPG では、どの friction system を意図的に外すかを早めに決める。「accessible」は tuning だけでなく system choice である。
- narrative-heavy な prototype では、authoring loop を plain-text で検索可能にし、launch 前に automated check を入れる。`Esoteric Ebb` の具体参照は `Ink` で、比較候補として `Yarn Spinner`、`Arcweave`、`articy:draft` を残しておく。

## Open questions

- どの部分を最初に自動化すべきか。prototype generation、playtesting、art support のどれが最初か。
- 最初の本格的 AI 補助 game-production testbed はどの project にするべきか。
- こうした demo から保存に値する最小 artifact は何か。code、prompt、playtest trace、curated screenshot のどれか。
- workspace の game plan に「spec 前に agent が design question をする」checkpoint を標準化すべきか。
- どの game experiment では AI を player-facing loop に入れ、どれでは development tooling のみに留めるべきか。
- branching dialogue、tutorial、quest logic を持つ spec には、標準の「narrative state and validation」節を入れるべきか。

## 関連ページ

- [phaser](../tools/phaser.md)
- [godot](../tools/godot.md)
- [next2d](../tools/next2d.md)
- [reversi-adventure](../projects/reversi-adventure.md)
- [branching-narrative-authoring](../patterns/branching-narrative-authoring.md)
- [component-oriented-gameplay](../patterns/component-oriented-gameplay.md)
