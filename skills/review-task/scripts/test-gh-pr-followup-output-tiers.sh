#!/usr/bin/env bash
set -euo pipefail

ROOT=$(git rev-parse --show-toplevel)
SCRIPT="$ROOT/skills/review-task/scripts/gh-pr-followup"

workdir=$(mktemp -d)
trap 'rm -rf "$workdir"' EXIT

state_dir="$workdir/state"
fake_gh="$workdir/fake-gh"

cat >"$fake_gh" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail

command_name="${1:-}"
shift || true

case "$command_name" in
    pr)
        subcommand="${1:-}"
        shift || true
        [ "$subcommand" = "view" ] || exit 1
        cat <<JSON
{
  "headRefOid": "${FAKE_HEAD_SHA}",
  "reviewDecision": "${FAKE_REVIEW_DECISION:-APPROVED}",
  "statusCheckRollup": [
    {
      "name": "test",
      "workflowName": "CI",
      "status": "COMPLETED",
      "conclusion": "SUCCESS",
      "detailsUrl": "https://example.com/checks/test"
    }
  ],
  "reviews": [
    {
      "id": "R1",
      "state": "COMMENTED",
      "author": {"login": "copilot-pull-request-reviewer[bot]"},
      "body": "Please tighten this.",
      "submittedAt": "2026-07-04T00:00:00Z",
      "commit": {"oid": "${FAKE_HEAD_SHA}"},
      "url": "https://example.com/review/R1"
    }
  ]
}
JSON
        ;;
    api)
        endpoint="${1:-}"
        shift || true
        case "$endpoint" in
            repos/*/issues/*/timeline)
                printf '%s\n' "${FAKE_TIMELINE_LINES}"
                ;;
            repos/*/pulls/*/comments)
                printf '%s\n' "${FAKE_COMMENT_LINES}"
                ;;
            *)
                exit 1
                ;;
        esac
        ;;
    *)
        exit 1
        ;;
esac
EOF
chmod +x "$fake_gh"

export GH_BIN="$fake_gh"
export GH_PR_FOLLOWUP_STATE_DIR="$state_dir"
export FAKE_REVIEW_DECISION="APPROVED"

timeline_lines_initial=$(cat <<'EOF'
{"id":101,"event":"review_requested","created_at":"2026-07-04T00:00:00Z","actor":{"login":"yoske"},"requested_reviewer":{"login":"copilot-pull-request-reviewer[bot]","type":"Bot"},"performed_via_github_app":{"slug":"copilot","name":"Copilot"},"commit_id":"sha-1","state":"PENDING"}
{"id":102,"event":"reviewed","created_at":"2026-07-04T00:01:00Z","actor":{"login":"copilot-pull-request-reviewer[bot]","type":"Bot"},"review":{"state":"COMMENTED","body":"Needs work","submitted_at":"2026-07-04T00:01:00Z"},"commit_id":"sha-1","url":"https://example.com/timeline/102"}
EOF
)
comment_lines_initial=$(cat <<'EOF'
{"id":201,"path":"skills/review-task/scripts/gh-pr-followup","line":42,"user":{"login":"copilot-pull-request-reviewer[bot]"},"body":"Nit","created_at":"2026-07-04T00:02:00Z","updated_at":"2026-07-04T00:02:00Z","commit_id":"sha-1","html_url":"https://example.com/comment/201"}
EOF
)

export FAKE_HEAD_SHA="sha-1"
export FAKE_TIMELINE_LINES="$timeline_lines_initial"
export FAKE_COMMENT_LINES="$comment_lines_initial"

compact_output=$("$SCRIPT" poll yoskeoka/vibe-coding-workspace 123)
jq -e '.checks == [{"name":"test","status":"COMPLETED","conclusion":"SUCCESS"}]' <<<"$compact_output" >/dev/null
jq -e 'has("reviews") | not' <<<"$compact_output" >/dev/null
jq -e '.state == {"changed_head":true,"last_checked_at":.state.last_checked_at}' <<<"$compact_output" >/dev/null
jq -e '.timeline_events | length == 2' <<<"$compact_output" >/dev/null
jq -e '.inline_comments | length == 1' <<<"$compact_output" >/dev/null

verbose_output=$("$SCRIPT" poll --verbose yoskeoka/vibe-coding-workspace 123)
jq -e '.reviews | length == 1' <<<"$verbose_output" >/dev/null
jq -e '.checks == [{"name":"test","workflow":"CI","status":"COMPLETED","conclusion":"SUCCESS","details_url":"https://example.com/checks/test"}]' <<<"$verbose_output" >/dev/null
jq -e '.state.file and .state.last_timeline_event_id == 102 and .state.last_review_comment_id == 201' <<<"$verbose_output" >/dev/null
jq -e '.timeline_events == [] and .inline_comments == []' <<<"$verbose_output" >/dev/null

timeline_lines_second=$(cat <<'EOF'
{"id":101,"event":"review_requested","created_at":"2026-07-04T00:00:00Z","actor":{"login":"yoske"},"requested_reviewer":{"login":"copilot-pull-request-reviewer[bot]","type":"Bot"},"performed_via_github_app":{"slug":"copilot","name":"Copilot"},"commit_id":"sha-1","state":"PENDING"}
{"id":102,"event":"reviewed","created_at":"2026-07-04T00:01:00Z","actor":{"login":"copilot-pull-request-reviewer[bot]","type":"Bot"},"review":{"state":"COMMENTED","body":"Needs work","submitted_at":"2026-07-04T00:01:00Z"},"commit_id":"sha-1","url":"https://example.com/timeline/102"}
{"id":103,"event":"reviewed","created_at":"2026-07-04T00:03:00Z","actor":{"login":"copilot-pull-request-reviewer[bot]","type":"Bot"},"review":{"state":"APPROVED","body":"Looks good","submitted_at":"2026-07-04T00:03:00Z"},"commit_id":"sha-1","url":"https://example.com/timeline/103"}
EOF
)
comment_lines_second=$(cat <<'EOF'
{"id":201,"path":"skills/review-task/scripts/gh-pr-followup","line":42,"user":{"login":"copilot-pull-request-reviewer[bot]"},"body":"Nit","created_at":"2026-07-04T00:02:00Z","updated_at":"2026-07-04T00:02:00Z","commit_id":"sha-1","html_url":"https://example.com/comment/201"}
{"id":202,"path":"skills/review-task/SKILL.md","line":188,"user":{"login":"copilot-pull-request-reviewer[bot]"},"body":"Please clarify verbose usage","created_at":"2026-07-04T00:04:00Z","updated_at":"2026-07-04T00:04:00Z","commit_id":"sha-1","html_url":"https://example.com/comment/202"}
EOF
)

export FAKE_TIMELINE_LINES="$timeline_lines_second"
export FAKE_COMMENT_LINES="$comment_lines_second"

incremental_output=$("$SCRIPT" poll yoskeoka/vibe-coding-workspace 123)
jq -e '.state.changed_head == false' <<<"$incremental_output" >/dev/null
jq -e '.timeline_events | map(.id) == [103]' <<<"$incremental_output" >/dev/null
jq -e '.inline_comments | map(.id) == [202]' <<<"$incremental_output" >/dev/null

export FAKE_HEAD_SHA="sha-2"
head_reset_output=$("$SCRIPT" poll yoskeoka/vibe-coding-workspace 123)
jq -e '.state.changed_head == true' <<<"$head_reset_output" >/dev/null
jq -e '.timeline_events | map(.id) == [101,102,103]' <<<"$head_reset_output" >/dev/null
jq -e '.inline_comments | map(.id) == [201,202]' <<<"$head_reset_output" >/dev/null

echo "gh-pr-followup output tier smoke test: PASS"
