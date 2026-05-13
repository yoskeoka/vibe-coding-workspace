# MCP CLI Token Reducer

## 背景

MCP は AI Agent との接続を意図したプロトコルだが、実際には JSON-RPC ベースであり、MCP server が返す JSON は多くの client 実装でほぼそのまま LLM に流れ込む。そのため、server 側が AI 向けに過剰な構造や冗長なフィールドを返すと、そのままトークン消費増につながる。

加えて、MCP の接続性や使い勝手は server だけでは決まらず、MCP client を実装する AI Agent ツールの準拠度に大きく依存する。DCR、Client ID Metadata Documents、scope selection strategy などを client が十分に実装していない場合、server 側が仕様に沿っていても接続できない、あるいは扱いづらいケースが起こる。

## 課題感

### 1. JSON がそのまま LLM に渡ることによるトークン浪費

- MCP server が巨大な JSON を返すと、そのまま LLM のコンテキストを圧迫しやすい
- 不要フィールドや説明過多の tool schema も token cost になる
- `structuredContent` と `content` の両立は後方互換上は自然だが、client 実装次第では重複情報がそのままモデルに渡りうる

### 2. 非準拠 client による接続性・運用性の悪化

- 仕様上は auth discovery, CIMD, DCR fallback, scope challenge handling などが整理されている
- しかし現実には AI Agent ツールごとに MCP client 実装の成熟度が違う
- その結果、server 側で頑張っても「特定 client では繋がらない」「最小権限でうまく authorize できない」「期待した scope 交渉ができない」などの問題が残る

## 調査結果サマリ

### 課題 1 について

- 仕様には軽減策はある
  - `tools/list` などの pagination
  - `resource_link` による本体分離
  - `annotations.audience` / `priority` による client 側の取捨選択
- ただし、これらは「client が賢く使えば軽減できる」類の仕組みであり、LLM に渡す payload を protocol として強制的に最小化するものではない
- 現行 draft では `structuredContent` を返す tool は後方互換のため `content` にも serialized JSON を返すことが推奨されており、client によっては重複転送が起こりうる
- `structuredContent` と `content` の扱いの不統一は issue 化されているが、現時点では整理・ガイダンス強化が中心で、抜本的な wire-level 変更は見えていない
- model size や capability に応じて tool schema を出し分ける提案もあるが、まだ proposal 段階

結論:

現行 MCP は token 効率を protocol レベルで十分に解決しているとは言いづらい。近い将来に自然解消する見込みは弱い。

### 課題 2 について

- auth まわりの仕様自体はかなり整理されている
  - Client ID Metadata Documents を推奨
  - DCR は fallback
  - `WWW-Authenticate` / protected resource metadata / scope challenge / step-up authorization も明文化
- ただし、これは「準拠 client はどう作るべきか」が定義されているという話であって、非準拠 client を server 側で救済できるという話ではない
- DCR の強化 proposal はあるが、client 実装不足の互換レイヤーになるものではない

結論:

この問題は仕様不足というより ecosystem 側の client 実装差分の問題であり、spec の更新だけで近い将来に一気に解消する可能性は低い。

## ソリューション仮説

`mcp` を AI Agent に直接食わせるのではなく、`CLI` を 1 段挟む。

### アイデア

MCP server との通信は専用 CLI が担い、AI Agent からはその CLI を通常の shell command / local tool として呼び出す。

CLI は MCP Tools にいったん対象を限定し、次を担当する:

- server 接続と初期化
- tool discovery
- description / schema をもとに CLI 引数化
- tool 実行
- 実行結果の整形
- 不要フィールド削減
- 必要に応じた text-first / compact-first な出力変換
- auth まわりの吸収

### この方式で狙えること

#### 課題 1 への効果

- LLM に渡す前に JSON を圧縮できる
- `structuredContent` と `content` のどちらを使うかを CLI 側で統制できる
- schema や result を human/model-oriented なテキストに再整形できる
- 大きい payload をそのまま毎回 model に渡さず、必要最小限の表示にできる

#### 課題 2 への効果

- MCP client として必要な仕様対応を CLI 側に寄せられる
- AI Agent ツールが不完全な MCP client を内蔵していても、CLI を呼べる限り MCP server を利用できる
- DCR/CIMD/scope challenge などの面倒を agent 製品ごとに期待しなくてよくなる

## 前提と割り切り

- CLI 化により、MCP server との往復や schema 解決のオーバーヘッドは発生する
- ただし、AI Agent が比較的長時間自律実行するユースケースなら、数秒から数十秒の追加コストは許容できる前提
- 初期スコープは MCP Tools に限定するのが妥当
- Resources / Prompts / Sampling まで最初から広げると複雑になりやすい

## 現時点の判断

- 課題は実在する
- しかも spec の近い将来の自然進化だけでは解消しきらない可能性が高い
- したがって、`MCP CLI wrapper / bridge / distiller` には十分に検討価値がある

## 既存実装の追加調査

### 大きく 2 系統ある

1. generic MCP CLI client
2. transport bridge / gateway

前者は MCP server に CLI で直接つなぎ、tools/resources/prompts を対話的あるいは非対話的に扱うもの。後者は `stdio` / `SSE` / `Streamable HTTP` など transport の差分を吸収するもの。

### 1. generic MCP CLI client

#### `apify/mcpc`

- 2026-05-02 時点で最も近い
- 自らを universal MCP CLI client と位置づけている
- `stdio` / `Streamable HTTP`、OAuth 2.1、永続 session、JSON 出力、`grep`、AI sandbox 向け proxy などを持つ
- README 上でも interactive shell、scripting、AI agents 向け利用を前面に出している
- `--max-chars` による出力 truncation や progressive tool discovery もある

評価:

- 「MCP を CLI として使う」という方向性は既に存在する
- ただし、主眼は汎用 client 化であり、LLM 投入前の payload distillation や `structuredContent` / `content` の戦略的統制までは前面に出ていない

#### `wong2/mcp-cli`

- CLI inspector 色が強い
- 非対話モードで `call-tool` / `read-resource` / `get-prompt` を呼べる
- scripting や inspection には使いやすい

評価:

- 発想としては近いが、主用途は inspection / manual invocation
- token 最適化や AI Agent 向け output shaping は主題ではない

#### `tilesprivacy/mcp-cli`

- remote MCP server 向け interactive CLI client
- OAuth と transport fallback を持つ
- ただし 2026-01-11 時点で archive 済み

評価:

- 参考にはなるが、現行の有力候補とは言いづらい

### 2. transport bridge / gateway

#### `geelen/mcp-remote`

- remote MCP server を受けて local stdio MCP server として見せる
- client 側の HTTP/OAuth 対応不足を吸収する、という意味で課題 2 に近い

評価:

- かなり近いが、あくまで MCP client <-> MCP server の transport / auth 互換レイヤー
- AI Agent から shell command として自然に使う CLI interface や output distillation は主眼ではない

#### `supercorp-ai/supergateway`

- `stdio <-> SSE / Streamable HTTP / WS` の変換を行う gateway
- remote server を local stdio に見せる用途も、local stdio server を network transport に出す用途もある

#### `sparfenyuk/mcp-proxy`

- `stdio <-> SSE / StreamableHTTP` の proxy
- transport の切替に特化している

#### `unarii/vikstra-bridge`

- stdio client から remote HTTP/HTTPS MCP server へつなぐ bridge

#### `EvalsOne/MCP-connect`

- local stdio MCP server を HTTP bridge / Streamable HTTP で公開する
- classic request/response bridge も持つ

#### `acehoss/mcp-gateway`

- MCP server を REST API / OpenAPI として見せる gateway
- MCP 非対応 client に使わせる方向の変換

評価:

- この系統は transport 差分や client 非対応を埋める点では強い
- ただし、MCP を一度受けたうえで AI Agent 向けに compact な CLI UX / text-first output / token-aware shaping を行うところまではあまり踏み込んでいない

### 現時点での差分認識

- 既存実装は存在する
- 特に `mcpc` は「MCP を CLI として使う」という意味でかなり近い
- 一方で、こちらの仮説の中核は単なる CLI client 化ではなく、
  - LLM に渡す前の token distillation
  - `structuredContent` と `content` の統制
  - tool schema / tool result の compact-first な再整形
  - protocol 実装不備の吸収を agent-facing CLI UX と一体で提供すること
- 少なくとも公開 README レベルでは、この組み合わせを主眼にした実装はまだ薄い

暫定結論:

- 類似実装はある
- ただし、`generic MCP CLI client` と `transport bridge` が中心であり、`AI Agent 向け distiller / wrapper` という観点では未充足の余地がある

## 次の調査候補

1. `mcpc` を基準点として、既存実装が解いていないギャップは何か
2. `generic MCP CLI client` と `AI Agent 向け distiller/wrapper` の境界はどこに引くべきか
3. 最小プロダクトの境界はどこか
4. auth 対応を初版でどこまで持つべきか
5. output の圧縮戦略をどの粒度で持つべきか

## 追加調査: `mcpc` はこの課題に使えるか

2026-05-02 時点で、`mcp-cli-handoff.md` の仮説に最も近い既存実装として `apify/mcpc` を追加で精査した。
比較対象として `wong2/mcp-cli` も確認した。

## 結論

- `mcpc` は、AI Agent 側の不完全な MCP client 実装を CLI 側で吸収する用途にはかなり使える
- 特に remote MCP の OAuth 2.1、CIMD/DCR、persistent session、dynamic discovery を agent から shell 経由で使わせる役として有力
- 一方で、本メモの中核にある
  - LLM に渡す前の payload distillation
  - `structuredContent` / `content` の戦略的統制
  - compact-first な agent-facing output
 までは `mcpc` 単体ではまだ不足がある
- したがって `mcpc` は「有力な土台」ではあるが、「そのまま完成品」ではない

## `mcpc` が使える点

### 1. client 実装不足の吸収

README と実装の両方で、`mcpc` は単なる inspector ではなく universal MCP CLI client として設計されている。

- `stdio` / Streamable HTTP をサポート
- persistent session を持つ
- OAuth 2.1 をサポート
- AI agents / code mode を前面に出している
- dynamic discovery を通じて tool schema の先読みを減らす思想がある

特に auth まわりは、本メモの「AI Agent 製品ごとの client 実装差分を CLI に寄せる」という仮説にかなり沿っている。

- OAuth provider 実装あり
- CIMD を client identity として扱う
- server が CIMD 非対応なら DCR fallback
- pre-registration も扱う

つまり agent 本体が不完全な MCP client でも、shell を叩けるなら `mcpc` を sidecar client として利用できる。

### 2. MCP spec 準拠志向

`mcpc` は README 上で full OAuth 2.1 support with CIMD and DCR を明示し、feature support table も比較的整理されている。

さらに package.json には conformance test があり、少なくとも spec 準拠を意識した開発姿勢はある。

ただし後述の通り、未実装機能も残っているため「かなり強い準拠志向」は確認できるが、「全面的な完成」はまだ言い切れない。

### 3. token 問題への課題感は明確にある

ここは今回の追加調査で一番重要だった点だが、`mcpc` は token 問題を少なくとも明確に認識している。

README / blog / 実装から確認できる点:

- Progressive tool discovery を save tokens and increase accuracy と説明している
- AI agents / code mode の説明で fewer tokens と明言している
- 公開 blog でも static tool loading wastes tokens を問題設定に置いている
- `--max-chars` がある
- human output で冗長な result を多少削る実装がある

したがって「token 削減をゴールに置いているか、少なくとも課題感があるか」という問いには、
少なくとも後者は明確に yes と言える。

## `mcpc` のギャップ

### 1. token 最適化の中心は dynamic discovery であって payload distillation ではない

`mcpc` は token 問題を認識しているが、その主な打ち手は:

- 必要な tool schema だけ後から読む
- shell / code mode に逃がす
- human-readable output を少し整える

であり、本メモで狙っているような

- tool result の意味圧縮
- field pruning
- semantic truncation
- agent-facing compact mode
- result payload の再設計

までは前面に出ていない。

つまり token 課題感はあるが、解き方の中心が違う。

### 2. `structuredContent` / `content` の統制は部分対応に留まる

実装を読むと、human mode では次の工夫がある。

- `content` が text-only なら、`structuredContent` を表示せず text を canonical view とみなす
- `structuredContent` と同値な JSON text block は重複として content 側から除外する

これは良い方向だが、まだ「軽い整形」であり、「distiller」と呼べるほどではない。

特に `--json` は raw MCP payload をそのまま返す方針で、AI code mode の例も `.content[0].text | fromjson` を前提にしている。
つまり spec JSON を崩さないのが優先で、agent-facing compact JSON はまだ提供していない。

### 3. truncation は blunt で、意味単位の圧縮ではない

`--max-chars` はあるが、これは文字数での単純打ち切りである。

- 重要フィールド優先
- summary + details 分離
- list head/tail の intelligent reduction
- `resource_link` への自動退避

のような意味単位の圧縮ではない。

そのため「大きい JSON をそのまま LLM に食わせる問題」を本質的に解くものではない。

### 4. resources 系の output shaping はまだ弱い

`resources-read` には `--raw` はあるが、

- `--output` は未実装
- `maxSize` も実質未整備
- resource 本文を agent 向けに compact 化する仕組みも弱い

という状態で、本メモの「大きい payload をそのまま毎回 model に渡さない」という要求とはまだ距離がある。

### 5. spec 準拠は強いが未実装領域も残る

README の feature table では:

- Roots: planned
- Elicitation: planned
- Completion: planned
- Sampling: not applicable

となっている。

また conformance test も package.json 上は initialize scenario だけなので、
spec 準拠性を重視していることは確かだが、広い範囲を conformance で網羅しているわけではない。

## `wong2/mcp-cli` との比較

`wong2/mcp-cli` は比較対象として有用だが、本メモの用途には本命ではない。

理由:

- 自らを CLI inspector と位置づけている
- non-interactive `call-tool` / `read-resource` / `get-prompt` はある
- OAuth つき remote transport にもある程度対応している
- ただし output shaping や token-aware UX はほぼ主題ではない
- non-interactive mode は raw JSON をそのまま返すに近い

発想としては近いが、「inspection / manual invocation」寄りであり、
「AI Agent 向け distiller / wrapper」としては `mcpc` よりさらに遠い。

## 保守状況

### `apify/mcpc`

2026-05-02 時点では十分アクティブ。

- GitHub API 上の `pushed_at` は 2026-04-29
- repo `updated_at` は 2026-05-01
- README 上の latest release は 2026-04-15
- open issue / open PR はあるが、2026-04-28 の新しい issue / PR が複数ある

したがって「止まっている repo」ではなく、改善提案や小さな PR を出す現実味はある。

ただし triage が非常に綺麗というほどではない。

- open PR に 2026-03-07 の古いものが残っている
- WIP 的な PR も残っている

なので「方向が合えば取り込まれる可能性はある」が、「必ず迅速に反応される」とは言い切れない。
出すなら小さく切った issue / PR の方がよい。

### `wong2/mcp-cli`

こちらは保守継続はされているものの、`mcpc` よりかなり弱い。

- GitHub API 上の `pushed_at` は 2025-08-13
- open issue がいくつか残っている
- open PR も 2025-09 由来のものが滞留している

比較対象としては参考になるが、こちらに寄せて課題解決を進める優先度は低い。

## 判断

現時点の実務判断としては次の通り。

- `mcpc` は「AI Agent の不完全な MCP client 実装を shell-sidecar で補う」第一候補として十分有力
- ただし「token-aware distiller / wrapper」としては未完成
- したがって `mcpc` をベースにして
  - issue を出す
  - 小さい PR を出す
 という方向は現実的

## `mcpc` に対して切り出しやすそうな差分

比較的通しやすそうな論点:

1. `structuredContent` / `content` の選択ポリシー
2. `tools-call` / `resources-read` の compact mode
3. field pruning
4. semantic truncation
5. agent-facing output mode と spec-preserving raw JSON mode の分離

特に最初の一歩としては、

- raw spec JSON は壊さない
- human / agent mode だけ改善する
- `tools-call` と `resources-read` に `--compact` を追加する

のような PR が最も現実的に見える。

## 現時点の最終判断

- `mcpc` はこの課題に「使える」
- ただし「そのままで十分」ではない
- 「MCP client 実装不足の吸収」は既にかなり担える
- 「token-aware distillation」は未充足で、ここに issue / PR の余地がある
- ただし issue 化や PR 提案の Go は、実際に `mcpc` を触って課題解決の素材になるかを見てから判断するのが妥当
- したがって現段階の主目的は「調査結果を、次回そのまま実機評価に接続できる粒度で残すこと」にある

## この handoff の位置づけ

この文書は現時点では issue 草案ではなく、次の実機評価に入る前の調査台帳である。

ここで確定させたいのは:

- 既存実装の中でどれが基準点として最有力か
- その基準点が何をすでに解いていて、何をまだ解いていないか
- 後で issue / PR 化する場合に、どの論点に分割すべきか

ここでまだ確定させないもの:

- `mcpc` に本当に issue を出すか
- こちらが直接 PR を作るか
- どの論点を upstream に切り出すかの優先順位
- こちらの課題解決に `mcpc` をそのまま採用するか、薄い wrapper を別に足すか

## handoff だけで現時点で言えること

実機評価前でも、次はかなり強く言える。

### 1. 基準点としては `mcpc` が最有力

- `generic MCP CLI client` 群の中で、AI Agent との併用を最も明示している
- OAuth / CIMD / DCR / persistent session / tool discovery の実装が比較的揃っている
- 保守が生きている

したがって「まず何を掘るか」という問いに対する答えは、現時点では `mcpc` でよい。

### 2. `mcpc` は token 問題を認識している

これは重要な確認ポイントで、少なくとも課題感はある。

ただしその中心は:

- static tool loading の回避
- dynamic discovery
- code mode への移行

であり、こちらの関心である:

- result payload 自体の distillation
- compact-first output
- `structuredContent` / `content` の agent-facing policy

とはズレがある。

### 3. 今の段階では「使えそう」止まりで、「使える」とまではまだ言い切らない方が正確

ここまでの追加調査で言えるのは:

- 設計思想はかなり近い
- 実装も一部近い
- 保守状況も悪くない

だが、実際にこちらの課題解決の素材になるかはまだ未確認である。

特に未確認なのは:

- 実運用でどの程度 token を減らせるか
- tool result / resource result を wrapper なしでそのまま agent に食わせられるか
- 小さな追加変更で upstream に載せやすい形にできるか

## 実機評価で確認すべき論点

issue 化の前に、少なくとも次を実際に触って確認する必要がある。

### 1. 接続面

- remote MCP server に対する login / connect / reconnect の体験
- OAuth あり server での CIMD / DCR fallback の実効性
- こちらが問題にしている AI Agent 側 client 不足を、本当に shell sidecar として吸収できるか

### 2. discovery 面

- `tools-list`
- `tools-get`
- `grep`

で、schema の先読み量をどれだけ減らせるか。

### 3. result 面

- `tools-call` の human output が agent-facing に十分 compact か
- `--json` の raw payload が downstream wrapper なしで扱えるか
- `structuredContent` と `content` が両方ある result で、どの程度冗長になるか
- `resource_link` / embedded resource / large text resource で実用上の不都合がどこに出るか

### 4. truncation / size control 面

- `--max-chars` が実運用で有効か
- 単純 truncation で壊れるケースがどれくらいあるか
- `resources-read` の巨大 payload をどう扱うのが現実的か

### 5. upstream 適合性

- こちらが欲しい改善が `mcpc` の product direction に自然に乗るか
- その改善が「spec-preserving raw mode を壊さず追加できる」か
- issue にしたときに upstream から理解されやすい言い方は何か

## 実機評価の前に handoff に残しておくべき材料

次回の評価を早くするため、この文書には少なくとも以下を残しておくとよい。

### 1. 調査対象の確定

- 第一候補: `apify/mcpc`
- 比較対象: `wong2/mcp-cli`
- bridge / gateway 群は補助比較であり、本命ではない

### 2. いまの仮説

- `mcpc` は client 実装不足の吸収にはかなり使える
- token-aware distillation には不足がある
- したがって、本命シナリオは
  - そのまま採用
  - 小さな upstream 改善
  - 必要なら薄い wrapper を上に足す
  のいずれか

### 3. 未確定事項

- wrapper が本当に必要か
- upstream 改善だけで十分か
- issue を出すべき論点が 1 本か複数本か
- こちらが直接 PR を作るべきか

### 4. issue 化の前提条件

次の 3 つが揃うまでは issue 作成に進まない方がよい。

1. 実際の server に接続して、こちらの課題に近い運用を一度通す
2. 冗長 payload や扱いづらい output の具体例を採取する
3. 「何を変えると改善するか」を CLI UX レベルで一文で言えるようにする

## 将来の issue 化に向けた論点の分割案

これは今すぐ issue を出すためのものではなく、後で混線しないための下書きである。

### 論点 A: `structuredContent` / `content` policy

- 何を canonical にするか
- human mode と agent mode で分けるべきか
- duplicate suppression をどこまでやるか

### 論点 B: compact output mode

- `tools-call` の compact mode
- `resources-read` の compact mode
- summary / details の二段表示

### 論点 C: size control

- `--max-chars` ではなく意味単位の圧縮が必要か
- head/tail や field-priority の導入余地
- large payload を `resource_link` 的な参照に逃がす余地

### 論点 D: raw spec mode と agent-facing mode の分離

- `--json` は raw spec preserving のまま維持
- それとは別に compact JSON / agent JSON を持つべきか

## 次回やること

次のターンでは、issue 作成ではなく実機評価に進む。

優先順:

1. `mcpc` を実際に使って接続・discovery・tool call の外形を確認する
2. token / payload / output shaping の観点で困る具体例を採取する
3. その結果をこの handoff に追記する
4. その後に初めて issue 化の是非を判断する
- 「token-aware distillation」は未充足で、ここに issue / PR の余地がある
- したがって、今後の基準点としては `mcpc` が最有力
