#!/usr/bin/env bash
#
# Merge a chain of branches, each into the next one, then push all of them:
#   branch-a -> branch-b -> branch-c -> ...
#
# Usage:
#   ./merge-branches.sh origin "branch-a
#   branch-b
#   branch-c"
#
#   printf 'branch-a\nbranch-b\nbranch-c\n' | ./merge-branches.sh origin
#
#   ./merge-branches.sh    # prompts for remote and branches
#
# Branches are newline separated; blank lines and lines starting with # are ignored.
# On merge conflicts the script pauses until the conflicts are resolved and the
# merge is committed. Once all merges are done, every branch is pushed to the remote.

set -euo pipefail

remote="${1:-}"

if [[ -z "$remote" ]]; then
  read -r -p "Remote (e.g. origin): " remote < /dev/tty
fi

if [[ -z "$remote" ]]; then
  echo "Usage: $0 <remote> [branches]" >&2
  exit 1
fi

if [[ $# -gt 1 ]]; then
  input="$2"
elif [[ ! -t 0 ]]; then
  input="$(cat)"
else
  echo "Branches in merge order, one per line (empty line to finish):"
  input=""
  while IFS= read -r line < /dev/tty && [[ -n "$line" ]]; do
    input+="$line"$'\n'
  done
fi

if ! git remote get-url "$remote" > /dev/null 2>&1; then
  echo "Error: remote '$remote' does not exist." >&2
  exit 1
fi

branches=()
while IFS= read -r line; do
  line="${line#"${line%%[![:space:]]*}"}"
  line="${line%"${line##*[![:space:]]}"}"
  [[ -z "$line" || "$line" == \#* ]] && continue
  branches+=("$line")
done <<< "$input"

if [[ ${#branches[@]} -lt 2 ]]; then
  echo "Error: at least 2 branches are required." >&2
  exit 1
fi

for branch in "${branches[@]}"; do
  if ! git rev-parse --verify --quiet "refs/heads/$branch" > /dev/null; then
    echo "Error: local branch '$branch' does not exist." >&2
    exit 1
  fi
done

if [[ -n "$(git status --porcelain --untracked-files=no)" ]]; then
  echo "Error: working tree has uncommitted changes. Commit or stash them first." >&2
  exit 1
fi

start_branch="$(git rev-parse --abbrev-ref HEAD)"

in_merge() {
  git rev-parse --verify --quiet MERGE_HEAD > /dev/null
}

for ((i = 1; i < ${#branches[@]}; i++)); do
  source_branch="${branches[i - 1]}"
  target_branch="${branches[i]}"

  echo
  echo "==> Merging '$source_branch' into '$target_branch'"
  git checkout --quiet "$target_branch"

  if git merge --no-edit "$source_branch"; then
    continue
  fi

  if ! in_merge; then
    echo "Error: merge of '$source_branch' into '$target_branch' failed." >&2
    exit 1
  fi

  while true; do
    echo
    echo "Merge conflicts in '$target_branch'. Resolve them and commit the merge."
    read -r -p "Press Enter to continue..." < /dev/tty

    if in_merge; then
      echo "Merge is still in progress (unresolved conflicts or not committed yet)."
    elif ! git merge-base --is-ancestor "$source_branch" HEAD; then
      echo "'$source_branch' is not merged into '$target_branch' (was the merge aborted?)."
      echo "Run 'git merge $source_branch', resolve and commit, or press Ctrl+C to exit."
    else
      break
    fi
  done
done

git checkout --quiet "$start_branch"

echo
echo "==> Pushing to '$remote': ${branches[*]}"
git push "$remote" "${branches[@]}"

echo
echo "==> Done. Merged and pushed: ${branches[*]}"
