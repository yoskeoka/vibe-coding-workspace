# Specifications

Put detailed specs here.
Start with a high-level file such as `system-overview.md`.

- `codex-hooks.md`: workspace root から child repo hook へ dispatch する Codex hook 契約
- `codex-stop-goal-reminder.md`: workflow branch の完了境界を確認する workspace-only Stop hook 契約
- `github-actions-pinning.md`: GitHub Actions の `uses:` 参照を `pinact` で管理する運用契約
- `japanese-textlint-ci.md`: changed Japanese Markdown に `textlint` を流し、warning と PR comment で結果を出す契約
- `slopless-ci.md`: changed Markdown に `slopless` を流し、warning と PR comment で結果を出す契約
- `workflow-context-contract.md`: workflow 文書の compact-read と所有境界の契約
