---
title: Build a 2D Space Shooter with Phaser and Claude Code
source_url: https://phaser.io/news/2026/02/phaser-claude-code-tutorial
source_type: video_backed_article
original_language: en
ingested_on: 2026-04-18
status: active
tags:
  - phaser
  - ai-game-dev
  - claude-code
  - browser-games
video_url: https://www.youtube.com/watch?v=247Z3jdw_hs
video_platform: youtube
channel: Peter Yang
published_on: 2026-02-04
duration: 20m00s
time_anchors:
  - "00:00: Peter Yang が子どもと Claude Code でゲームを作る理由"
  - "01:19: project setup"
  - "02:56: 無料 pixel-art asset 選定"
  - "05:33: Claude Code の AskUserQuestion flow で spec を作る"
  - "08:56: MVP build と反復"
  - "13:06: boss battle のデバッグ"
  - "15:24: GitHub と Vercel での公開"
selected_screenshots:
  - path: ../../assets/source-images/2026/phaser-claude-code-space-shooter/00-00-intro.jpg
    anchor: "00:00"
    note: "チュートリアル全体像と 5 ステップ構成"
  - path: ../../assets/source-images/2026/phaser-claude-code-space-shooter/01-19-project-setup.jpg
    anchor: "01:19"
    note: "Claude Code を使うための Cursor terminal 設定"
  - path: ../../assets/source-images/2026/phaser-claude-code-space-shooter/02-56-assets.jpg
    anchor: "02:56"
    note: "Legacy Pixel Collection を探して取り込む場面"
  - path: ../../assets/source-images/2026/phaser-claude-code-space-shooter/05-33-spec.jpg
    anchor: "05:33"
    note: "AskUserQuestion による spec 作成"
  - path: ../../assets/source-images/2026/phaser-claude-code-space-shooter/08-56-mvp.jpg
    anchor: "08:56"
    note: "playable milestone の反復"
  - path: ../../assets/source-images/2026/phaser-claude-code-space-shooter/13-06-boss-debug.jpg
    anchor: "13:06"
    note: "screenshot を使った boss デバッグ"
  - path: ../../assets/source-images/2026/phaser-claude-code-space-shooter/15-24-ship.jpg
    anchor: "15:24"
    note: "GitHub と Vercel での deploy flow"
named_entities:
  - Phaser
  - Claude Code
  - Cursor
  - AskUserQuestion
  - Animus
  - Legacy Pixel Collection
  - GitHub
  - Vercel
  - Peter Yang
related_pages:
  - ../../wiki/tools/phaser.md
  - ../../wiki/topics/ai-game-dev.md
---

# Build a 2D Space Shooter with Phaser and Claude Code

## ここで重要な理由

これは、人が読めるゲームアイデアから始め、質問を通じて spec にし、playable milestone を作り、最後に hosted demo まで出す AI 補助 browser-game prototype の簡潔な実例です。

## 要約

- Phaser の記事は、Peter Yang による 20 分の retro vertical space shooter チュートリアルを紹介している。
- ワークフローには Cursor、Claude Code、無料 pixel art asset、spec 作成のための `AskUserQuestion`、playable milestone の反復、GitHub と Vercel での deploy が含まれる。
- 出来上がるゲームには、縦スクロールシューティングの基礎、enemy wave、power-up、boss battle と health bar、screen shake、やや複雑な boss encounter のデバッグが入っている。
- 元の Peter Yang 投稿は、`space-pixel-shooter.vercel.app` の playable demo、Animus の `Legacy Pixel Collection`、setup から shipping までの 5 段階構成という具体アンカーも補っている。

## ワークスペースでの含意

- 再利用できるパターンは、「一発プロンプトでゲームを作らせる」ことではない。asset を集め、agent に design question をさせ、spec を書き、薄い MVP を作り、制御しやすい milestone で機能を足す流れである。
- Phaser は browser-first の高速ゲーム実験に引き続き向いている。Vercel に deploy するだけで共有可能な demo にできる。
- 子ども向けや音声主導の共同制作では、人間側にコードやタイピングを強いる代わりに、AI coding agent に design intent を渡して進める例として参考になる。
- ワークスペースのゲーム計画では、MVP first、その後 wave / power-up、最後に boss mechanic と polish という milestone discipline の資料として残す価値がある。

## ソースアクセスメモ

Phaser のページは ingest 時、直接ブラウザ取得を preview gate へ redirect していたが、検索 index から article metadata と本文要約は取得できた。関連する Peter Yang の投稿と Class Central の記述が、チュートリアル構成、動画長、著者、シラバス水準の流れを裏づけている。

## ソース

- source URL: https://phaser.io/news/2026/02/phaser-claude-code-tutorial
- video URL: https://www.youtube.com/watch?v=247Z3jdw_hs
- subtitle source: 2026-04-18 に `yt-dlp` で YouTube の `en-orig` 自動字幕を取得
- transcript policy: 字幕全文は保持せず、source link、time anchor、screenshot、詳細な segment note のみを durable な KB コピーとして残す

## アンカー別の動画メモとスクリーンショット

### 00:00: Peter Yang が子どもと Claude Code でゲームを作る理由

![チュートリアル全体像と 5 ステップ構成](../../assets/source-images/2026/phaser-claude-code-space-shooter/00-00-intro.jpg)

冒頭で、Claude Code は週末の共同制作ツールとして紹介される。Peter Yang は 7 歳の子どもと一緒にゲームを作っていると話し、それを beginner-friendly な流れの動機にしている。過去の例として underwater shooter や AI 生成画像を使った animal hospital game を見せたうえで、今回の対象を retro 2D space shooter として提示する。手順は project setup、asset 発見、spec 作成、MVP / milestone 実装、最後の仕上げと shipping の 5 ステップ。

### 01:19: project setup

![Claude Code を使うための Cursor terminal 設定](../../assets/source-images/2026/phaser-claude-code-space-shooter/01-19-project-setup.jpg)

setup は空の `space-shooter-game` フォルダを Cursor で開くところから始まる。Cursor を選ぶ理由は、file tree で Claude Code の編集が見え、同時に terminal も使えるからだと説明される。導入の振り返りでは Claude Code の quick-start command に触れ、そのまま Cursor terminal から Claude Code を起動する。Yang は、この種の disposable な趣味フォルダでは `claude --dangerously-skip-permissions` を好むとしており、permission prompt を減らして agent を自由に動かしている。ただし、fun prototype 向けの話であり、sensitive な本番作業向けではないと意図的に範囲を絞っている。

### 02:56: 無料 pixel-art asset 選定

![Legacy Pixel Collection を探して取り込む場面](../../assets/source-images/2026/phaser-claude-code-space-shooter/02-56-assets.jpg)

asset セクションでは、音声 dictation と Claude の web search を使って、無料の retro 2D space-shooter pixel art を探す様子が示される。Claude は OpenGameArt や itch.io を提案するが、Yang は具体性を保つため Animus の `Legacy Pixel Collection` を使う。pack をダウンロードして project folder に drag し、Claude に folder を調べさせて、space shooter に向いた asset を列挙させる。重要なのは、実装前に agent が実際の filename を見ていることだ。これにより、player ship、enemy、bullet、power-up、background、boss art を本当に存在する file に結び付けられる。

### 05:33: Claude Code の AskUserQuestion flow で spec を作る

![AskUserQuestion による spec 作成](../../assets/source-images/2026/phaser-claude-code-space-shooter/05-33-spec.jpg)

spec 用 prompt では、Claude Code に要件を書かせ、作業を 3 つの playable milestone に分け、使う pixel art asset を紐づけ、不明点があれば `AskUserQuestion` を使うよう求めている。Yang は 3 つの実践を強調する。build 前に spec を書くこと、作業を testable な milestone に分けること、asset 対応を明示すること。Claude は technology platform、scrolling direction、希望機能を質問し、Yang は Phaser、top-down な vertical shooter、power-up、boss battle、score、multiple wave を選ぶ。生成された spec は、vertical scrolling shooter、基本戦闘、progression / power-up、boss battle、polish を含む。加えて、生成 spec は盛り込みすぎることが多いので、人間が scope を切るべきだとも警告する。

### 08:56: MVP build と反復

![playable milestone の反復](../../assets/source-images/2026/phaser-claude-code-space-shooter/08-56-mvp.jpg)

最初の実装依頼では milestone 1 を作らせ、linked pixel art を使うよう再度念押しする。少し待つと、移動できる player ship、shooting、enemy、粗いながらも動く background を持つゲームがローカルで走る。milestone 2 では enemy variety、screen shake、health bar が追加されるが、一部 sprite はまだ不自然である。milestone 3 前に、Yang は boss battle にどの sprite を使うのかを問い、最初の答えが弱いと larger Legacy Pixel Collection を検索するよう促す。Claude は top-down boss folder を見つけ、そのまま進める。最終的なゲームには wave、power-up、boss trigger が入り、gameplay tuning と art 問題が見え始める程度の摩擦も残る。

### 13:06: boss battle のデバッグ

![screenshot を使った boss デバッグ](../../assets/source-images/2026/phaser-claude-code-space-shooter/13-06-boss-debug.jpg)

この区間は、visual feedback の重要さを示す。ゲームは boss 接近を警告するが、boss は現れず、UI は HP 0 を表示する。Yang は screenshot を撮って Claude Code に貼り、失敗内容を説明する。debug は初期 build より時間がかかり、さらに test を速くするため boss を直接呼び出す hotkey も追加する。根本原因は、boss HP を Phaser object 上に持たせたため engine 側の振る舞いと干渉したことで、HP を scene-level state に移すと解決した。大きな教訓は、難しい bug でも反復を続け、screenshot を証拠として使い、fix 後に root cause の説明も agent に求めること。

### 15:24: GitHub と Vercel での公開

![GitHub と Vercel での deploy flow](../../assets/source-images/2026/phaser-claude-code-space-shooter/15-24-ship.jpg)

公開フェーズでは、ローカル folder 名をそのまま使った GitHub repository を作り、remote URL をコピーして Claude Code に commit / push をさせる。Yang は、tutorial 版では asset pack を必要以上に push しているので、本来は実際に使う file だけに絞るべきだとも述べる。hosting は Vercel を使い、無料 account を作り、新規 project に Git URL を貼って deploy する。公開 URL が発行され、最後の playthrough で boss fight、power-up、level completion、background 変化が確認される。締めでは、setup、asset、`AskUserQuestion` 付き spec、3 milestone、bug fixing、deployment という 5 段階ループを繰り返している。

## フォローアップ

- ワークスペースで Phaser prototype を始めるとき、この flow が実装前の `docs/specs/` 質問票としてよいのか、coding tool 内の対話ループとしてよいのかを試す。
- 実際に適用するなら、使った prompt と milestone 境界も保存する。
