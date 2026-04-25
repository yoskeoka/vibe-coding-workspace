---
title: OpenAI built a Tactical RPG with Phaser
source_url: https://phaser.io/news/2026/03/gpt-5-4-phaser-game-tactical-rpg
source_type: video_backed_article
original_language: en
ingested_on: 2026-04-11
status: active
tags:
  - phaser
  - ai-game-dev
  - tactical-rpg
video_url: https://www.youtube.com/watch?v=rvdUBieefR0
video_platform: youtube
channel: Matthew Berman
published_on: 2026-03-06
duration: 16m27s
time_anchors:
  - "00:00-02:45: 発表者が GPT-5.4 を知識作業、コーディング、ブラウザ利用、agentic task 向けモデルとして位置づける"
  - "06:57-08:03: ブログ記事の要約として reasoning、coding、agentic workflow、planning、tools、computer use、long context work を強調"
  - "08:03-09:24: OSWorld / computer-use benchmark と Gmail 自動化デモ"
  - "09:39-10:05: JSON 風入力を UI に流し込む一括 data-entry デモ"
  - "10:05-10:58: 軽い指定プロンプトから作られた theme-park simulation デモ"
  - "10:58-11:16: Phaser 記事が強調する 2D tactical RPG デモ区間"
  - "11:20-12:09: pricing 議論と frontier output のコスト警告"
  - "12:10-12:52: OpenClaw 向けの prompting 助言と model-specific prompt guide 推奨"
selected_screenshots:
  - path: ../../assets/source-images/2026/openai-gpt-5-4-phaser-tactical-rpg/07-25-plan-first.jpg
    anchor: "07:25"
    note: "planning、tools、computer use、long-context agents に関する blog post 抜粋"
  - path: ../../assets/source-images/2026/openai-gpt-5-4-phaser-tactical-rpg/08-45-osworld.jpg
    anchor: "08:45"
    note: "OSWorld の検証精度と tool yield 比較グラフ"
  - path: ../../assets/source-images/2026/openai-gpt-5-4-phaser-tactical-rpg/10-18-theme-park-demo.jpg
    anchor: "10:18"
    note: "生成された isometric asset と管理 UI を含む theme-park simulation デモ"
  - path: ../../assets/source-images/2026/openai-gpt-5-4-phaser-tactical-rpg/11-05-rpg-demo.jpg
    anchor: "11:05"
    note: "movement、attack、wait、cancel、end-turn 操作を備えた 2D tactical RPG デモ"
named_entities:
  - GPT-5.4
  - Codex
  - GPT-5.3 Codex
  - OSWorld Verified
  - Playwright Interactive
  - Phaser
  - OpenClaw
  - Matthew Berman
related_pages:
  - ../../wiki/tools/phaser.md
  - ../../wiki/topics/ai-game-dev.md
---

# OpenAI built a Tactical RPG with Phaser

## ここで重要な理由

これは、趣味のブラウザゲーム試作でも現実味のある AI 補助ゲーム開発ワークフローを短く引ける retrieval anchor です。

## 要約

- 記事自体は薄いが、埋め込み動画から `GPT-5.4 + Codex + Playwright Interactive + image generation + Phaser` という具体的なツール連鎖を保存している。
- 動画全体は Phaser 記事より広く、GPT-5.4 を知識作業、コード、browser/computer use、planning、tools、agentic workflow 向けの汎用モデルとして位置づけ、benchmark、作業デモ、ゲームデモ、価格、prompt migration 助言まで扱う。
- Phaser に関係する区間は `10:05-11:16` 付近で、light prompt からの visual/gameplay iteration による theme-park simulation と 2D tactical RPG を、より厚みのある browser-game prototype 例として見せている。
- したがって、これは再現手順付き実装ガイドではなく、Phaser 系 browser game が agent-driven prototype の有力対象になりつつあることを示す証拠として扱うべき。

## ワークスペースでの含意

- 技術実装ガイドではなく proof point として扱う。コード、アーキテクチャ、プロンプト、再現可能な build path は公開されていない。
- 耐久的な価値は、model-driven code generation、browser/computer use、plan-first、visual feedback、generated assets、game/UI iteration がひとつのループにあること。
- 今後 Phaser 実験をするときは、ゲームデモ部分を手動検証の観点として使える。controls、state panel、turn flow、movement/attack affordance、asset の噛み合い、見た目だけでなく simulation が機能しているかを見る。
- pricing と prompt-guide の話は、高性能モデル利用を意図的に絞るべきという警告でもある。再利用できる文脈は cache し、高価な出力は効果の大きい場面に限り、model family を切り替えるなら prompt も書き換えるべき。

## ソース

- source URL: https://phaser.io/news/2026/03/gpt-5-4-phaser-game-tactical-rpg
- video URL: https://www.youtube.com/watch?v=rvdUBieefR0
- subtitle source: 2026-04-22 時点で `yt-dlp --list-subs` により YouTube の `en-orig` 自動字幕が利用可能だった
- transcript policy: 字幕全文はこのリポジトリに複製せず、source link、time anchor、selected screenshot、segment note のみを durable な KB コピーとして保持する

## アンカー別の動画メモとスクリーンショット

### 00:00-02:45: モデルの位置づけ

Matthew Berman は `GPT-5.4` を、幅広い world knowledge、coding 力、browser use、agent work を併せ持つモデルとして紹介する。以前のように、会話や知識に強いモデルと coding に強いモデルを別々に選ぶ時代との対比が主眼で、この動画は GPT-5.4 を chat 専用でも code 専用でもない、現実の knowledge work 向けの収束モデルとして位置づけている。

### 06:57-08:03: planning、tools、computer use、long context

![planning、tools、computer use、long-context agents に関する blog post 抜粋](../../assets/source-images/2026/openai-gpt-5-4-phaser-tactical-rpg/07-25-plan-first.jpg)

発表者は、今回のリリースを reasoning、coding、agentic workflow をまとめて改善したものとして要約する。ワークスペース的に重要なのは、実行前の planning、tool/software environment、spreadsheet や document のような professional artifact、computer use、long-context agent task を前面に出している点である。これは、このワークスペース自身の plan-first 運用と直接かみ合う。

### 08:03-09:24: OSWorld と computer-use デモ

![OSWorld の検証精度と tool yield 比較グラフ](../../assets/source-images/2026/openai-gpt-5-4-phaser-tactical-rpg/08-45-osworld.jpg)

OSWorld 議論では、tool yield 回数に対する精度が主題になる。発表者は、GPT-5.4 が GPT-5.2 より少ない tool call で高い精度に到達しており、tool-driven work を理論上より安く効率的にできると読む。その後、sent mail の確認、label 付与、calendar invite 作成といった Gmail 風 automation を見せる。ここで残すべき anchor は benchmark の数値そのものではなく、agent 品質を task success、tool call 数、screenshot や UI state への反応で測るという評価の形である。

### 09:39-10:05: 一括データ入力

このデモでは、JSON 風の構造化入力を UI にほぼリアルタイムで流し込む。ワークスペース目線では、UI automation デモを評価するとき、データ忠実性、速度、対象サイトが anti-automation blocker なしで agentic use を許すかを見るべきだという示唆になる。

### 10:05-10:58: テーマパーク simulation デモ

![生成された isometric asset と管理 UI を含む theme-park simulation デモ](../../assets/source-images/2026/openai-gpt-5-4-phaser-tactical-rpg/10-18-theme-park-demo.jpg)

最初のゲームデモは `Miniature` という isometric な theme-park simulation で、time-speed control、build palette、生成された attraction、path、visitor dot、funds、guest count、happiness、cleanliness、park rating、ride status、guest reaction を含む。発表者は、軽い指定の prompt から作ったと述べ、見た目だけでなく logic と metrics も prototype に入っていると強調する。

### 10:58-11:16: 2D tactical RPG デモ

![movement、attack、wait、cancel、end-turn 操作を備えた 2D tactical RPG デモ](../../assets/source-images/2026/openai-gpt-5-4-phaser-tactical-rpg/11-05-rpg-demo.jpg)

ここが Phaser 記事の主対象区間で、tile map、character card、movement range、attack / wait / cancel / end-turn 操作、battle log、turn state、generated character art を備えた 2D tactical RPG が映る。retrieval の観点では、turn-based combat UI と tactical movement affordance が見えるため、theme-park 区間よりも Phaser / game-dev 向け anchor として強い。

### 11:20-12:09: pricing

pricing 区間では、frontier model の利用コストが依然高いこと、とくに output-heavy な利用や Pro 系が重いことを警告している。input cache は一部の緩和策になるが、本当に厳しいのは output cost だと述べる。ワークスペース実験では、小さな prompt、cache / 再利用文脈、手元の検証ループ、高価なモデル呼び出しを結果が変わる場面にだけ使うことが重要になる。

### 12:10-12:52: OpenClaw と prompt migration

発表者は `OpenClaw` で GPT-5.4 を主力モデルとして試す価値を勧めつつ、`Opus` や Claude 系とは prompting が異なると警告する。実務的な助言は、model-specific な prompting guide を使い、必要ならモデル族ごとに別 prompt set を持つこと。このワークスペースでも、prompt 振る舞いは versioned workflow knowledge として残し、frontier model 間で同じ書き方がそのまま通ると仮定しないほうがよい。

## フォローアップ

- OpenAI が RPG の repo、技術記事、詳細な showcase page を公開したら、より強い実装知見を追記する。
