---
title: 分岐ナラティブ執筆
last_reviewed: 2026-04-25
status: seed
sources:
  - ../../sources/2026/2026-04-22-unity-interactive-writing-esoteric-ebb.md
---

# 分岐ナラティブ執筆

## 現在のシグナル

大規模な非線形 RPG の執筆は、プレイヤー状態へのフィードバックを保ちながら、通常の文章執筆に近い速度で書ける authoring loop を維持できるかで決まる。`Esoteric Ebb` の振り返りは、narrative ambition を `Ink`、Unity 連携、text search、命名規約、後半の bug-fix 圧力と結び付けている点が有用。

## このワークスペースで重要な理由

ワークスペースのゲーム project も、将来的に dialogue、choice、tutorial、branching scenario を持つかもしれない。CRPG でなくても、選択肢を局所的には簡単に増やせる一方、全体整合性の確認は難しくなるという失敗パターンをそのまま継承しうる。

## 再利用できるパターン

- content を書く前に interaction contract を定義する。choice density、dynamic-content の目標、設計が authored branching なのか systemic emergence なのかを先に決める。
- branch を安く扱える tool を選ぶ。plain-text review、速い recompile、project-wide search、Unity や engine との統合、custom command の作りやすさが重要。
- narrative state は明示的な data model として扱う。area、quest、quest-point、character、tutorial などの stable prefix を用い、variable を探しやすく監査しやすくする。
- 重要な choice には visible feedback を必須にする。後の text、quest state、UI、ending logic を何も変えない variable は narrative debt になりやすい。
- content 規模が膨らむ前に validation を入れる。syntax check、unused variable report、branch reachability check、高リスク経路の playtest script を用意する。

## ツールアンカー

- `Ink`: Unity との相性が強い plain-text interactive fiction scripting。
- `Yarn Spinner`: Unity project でよく使われる dialogue system。
- `Arcweave`: visual と text の両面を持つ interactive narrative planning tool。
- `articy:draft`: narrative design と game content planning の suite。

## Open questions

- choice や tutorial が分岐する場合、ワークスペースの game spec に標準の「narrative state and validation」節を入れるべきか。
- `Ink` や `Yarn Spinner` の project file から、どんな validation check を安く自動生成できるか。
- ここでは visual planning tool が有効か、それとも plain text と search のほうがワークスペースに合うか。
