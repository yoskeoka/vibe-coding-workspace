# Core Beliefs

These beliefs guide our technical trade-offs.

1.  **AI-First**: optimizing for context retrieval by LLMs is more important than optimizing for human navigation in deep folder structures.
2.  **Correctness over Speed**: Specs must be updated before code.
3.  **Don't fix what isn't broken**: Working, stable code should not be refactored purely for aesthetics (e.g., DRY). The risk of breaking something and the cost of re-testing outweigh marginal improvements in code elegance — especially when the duplicated code is unlikely to change.
4.  **Human Review Over Token Burn**: Prefer preparing clear artifacts for quick human verification (code, PR descriptions, quality-gate results, compact logs) over spending large agent context to automate checks a human can review with little effort.
5.  **Trim Tool Output at the Source**: When tools or APIs return large responses but the workflow needs only a few fields, add `jq`, scripts, or other structured filters so agents read the decision-relevant subset instead of raw payloads.
