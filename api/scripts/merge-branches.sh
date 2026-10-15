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

log() {
  echo "==> $*"
}

error() {
  echo "==> ERROR: $*" >&2
}

fail() {
  error "$@"
  exit 1
}

# Catch-all for any command that fails without its own error message.
trap 'error "Command failed (exit $?) at line $LINENO: $BASH_COMMAND"' ERR

remote="${1:-}"

if [[ -z "$remote" ]]; then
  read -r -p "Remote (e.g. origin): " remote < /dev/tty
fi

if [[ -z "$remote" ]]; then
  error "Remote is required."
  echo "Usage: $0 <remote> [branches]" >&2
  exit 1
fi

if [[ $# -gt 1 ]]; then
  input="$2"
elif [[ ! -t 0 ]]; then
  input="$(cat)"
else
  echo "Branches in merge order, one per line. Press Enter on an empty line to start."
  input=""
  n=1
  while IFS= read -r -p "  $n. " line < /dev/tty && [[ -n "$line" ]]; do
    input+="$line"$'\n'
    n=$((n + 1))
  done
fi

if ! git remote get-url "$remote" > /dev/null 2>&1; then
  fail "Remote '$remote' does not exist."
fi

branches=()
while IFS= read -r line; do
  line="${line#"${line%%[![:space:]]*}"}"
  line="${line%"${line##*[![:space:]]}"}"
  [[ -z "$line" || "$line" == \#* ]] && continue
  branches+=("$line")
done <<< "$input"

if [[ ${#branches[@]} -lt 2 ]]; then
  fail "At least 2 branches are required, got ${#branches[@]}."
fi

for branch in "${branches[@]}"; do
  if ! git rev-parse --verify --quiet "refs/heads/$branch" > /dev/null; then
    fail "Local branch '$branch' does not exist."
  fi
done

if [[ -n "$(git status --porcelain --untracked-files=no)" ]]; then
  fail "Working tree has uncommitted changes. Commit or stash them first."
fi

start_branch="$(git rev-parse --abbrev-ref HEAD)"
branch_count=${#branches[@]}
merge_count=$((branch_count - 1))

in_merge() {
  git rev-parse --verify --quiet MERGE_HEAD > /dev/null
}

log "Remote: $remote"
log "Merge order: $(printf '%s -> ' "${branches[@]}" | sed 's/ -> $//')"

for ((i = 1; i < branch_count; i++)); do
  source_branch="${branches[i - 1]}"
  target_branch="${branches[i]}"

  echo
  log "[$i/$merge_count] Checking out '$target_branch'"
  if ! git checkout --quiet "$target_branch"; then
    fail "Could not check out '$target_branch'."
  fi

  if git merge-base --is-ancestor "$source_branch" HEAD; then
    log "[$i/$merge_count] '$target_branch' already contains '$source_branch', nothing to merge"
    continue
  fi

  log "[$i/$merge_count] Merging '$source_branch' into '$target_branch'"
  if git merge --no-edit "$source_branch"; then
    log "[$i/$merge_count] Merged '$source_branch' into '$target_branch'"
    continue
  fi

  if ! in_merge; then
    fail "Merge of '$source_branch' into '$target_branch' failed (see git output above). You are on '$target_branch'."
  fi

  error "Merge conflicts while merging '$source_branch' into '$target_branch':"
  git diff --name-only --diff-filter=U | sed 's/^/      /' >&2
  while true; do
    echo
    echo "Resolve the conflicts in '$target_branch' and commit the merge."
    read -r -p "Press Enter to continue..." < /dev/tty

    if in_merge; then
      error "Merge is still in progress (unresolved conflicts or not committed yet)."
    elif ! git merge-base --is-ancestor "$source_branch" HEAD; then
      error "'$source_branch' is not merged into '$target_branch' (was the merge aborted?)."
      echo "Run 'git merge $source_branch', resolve and commit, or press Ctrl+C to exit." >&2
    else
      log "[$i/$merge_count] Merged '$source_branch' into '$target_branch' (conflicts resolved)"
      break
    fi
  done
done

echo
log "Switching back to '$start_branch'"
if ! git checkout --quiet "$start_branch"; then
  error "Could not switch back to '$start_branch', continuing with push."
fi

failed_pushes=()
for ((i = 0; i < branch_count; i++)); do
  branch="${branches[i]}"

  echo
  log "[$((i + 1))/$branch_count] Pushing '$branch' to '$remote'"
  if git push --quiet "$remote" "$branch"; then
    log "[$((i + 1))/$branch_count] Pushed '$branch' to '$remote'"
  else
    error "Failed to push '$branch' to '$remote'."
    failed_pushes+=("$branch")
  fi
done

echo
if [[ ${#failed_pushes[@]} -gt 0 ]]; then
  fail "Merged all branches, but failed to push: ${failed_pushes[*]}"
fi

log "Done. Merged and pushed: ${branches[*]}"
