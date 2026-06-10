#!/usr/bin/env bash
set -euo pipefail

REPO="${1:-yannisobert/Git-Health-Check}"
ROOT="$(cd "$(dirname "$0")/.." && pwd)"

if ! command -v gh >/dev/null 2>&1; then
  echo "Error: GitHub CLI (gh) is required. Install with: brew install gh"
  exit 1
fi

if ! gh auth status >/dev/null 2>&1; then
  echo "GitHub CLI is not authenticated."
  echo "Run: gh auth login"
  echo "  - Host: github.com"
  echo "  - Protocol: SSH (recommended with github-perso)"
  exit 1
fi

echo "→ Repository: $REPO"

existing_rulesets() {
  gh api "repos/$REPO/rulesets" --jq '.[].name' 2>/dev/null || true
}

apply_ruleset() {
  local file="$1"
  local name
  name="$(python3 -c "import json,sys; print(json.load(open(sys.argv[1]))['name'])" "$file")"

  if existing_rulesets | grep -qx "$name"; then
    echo "  ✓ Ruleset already exists: $name (skipped)"
    return
  fi

  echo "  → Creating ruleset: $name"
  gh api "repos/$REPO/rulesets" --method POST --input "$file"
  echo "  ✓ Created: $name"
}

echo "→ Applying branch rulesets..."
apply_ruleset "$ROOT/.github/rulesets/protect-main.json"
apply_ruleset "$ROOT/.github/rulesets/protect-dev.json"

if ! gh api "repos/$REPO/branches/dev" >/dev/null 2>&1; then
  echo "→ Branch dev does not exist on remote."
  echo "  Create it locally with:"
  echo "    git checkout -b dev && git push -u origin dev"
else
  echo "→ Branch dev exists on remote."
fi

echo ""
echo "Done. Verify at:"
echo "  https://github.com/$REPO/settings/rules"
