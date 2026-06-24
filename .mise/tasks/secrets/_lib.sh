#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
APP_DIR="${SECRETS_APP_DIR:-$(cd -- "$SCRIPT_DIR/../../.." && pwd)}"
cd "$APP_DIR"

RECIPIENTS_PATH=".secrets/recipients.txt"
SOPS_CONFIG_PATH=".sops.yaml"
KEY_DIR=".secrets/keys"
TMP_FILES=()

cleanup() {
  if (( ${#TMP_FILES[@]} > 0 )); then
    rm -f "${TMP_FILES[@]}"
  fi
}
trap cleanup EXIT

die() {
  printf 'secrets: %s\n' "$*" >&2
  exit 1
}

usage() {
  cat >&2 <<'USAGE'
usage:
  secrets new-key <user>
  secrets add-key <user> <age-recipient>
  secrets remove-user <user>
  secrets users
USAGE
}

validate_user() {
  local user="$1"
  if [[ ! "$user" =~ ^[A-Za-z0-9][A-Za-z0-9._-]*$ ]]; then
    die "user must start with a letter or number and contain only letters, numbers, dots, underscores, or dashes"
  fi
}

validate_age_recipient() {
  local key="$1"
  if [[ ! "$key" =~ ^age1[0-9a-z]+$ ]]; then
    die "age recipient must look like age1..."
  fi
}

recipient_rows() {
  if [[ ! -f "$RECIPIENTS_PATH" ]]; then
    return 0
  fi

  awk '
    function fail(message) {
      printf "line %d: %s\n", NR, message > "/dev/stderr"
      exit 1
    }
    {
      text = $0
      sub(/^[[:space:]]+/, "", text)
      sub(/[[:space:]]+$/, "", text)
      if (text == "" || substr(text, 1, 1) == "#") {
        next
      }
      count = split(text, fields, /[[:space:]]+/)
      if (count != 2) {
        fail("expected <user> <age-recipient>")
      }
      user = fields[1]
      key = fields[2]
      if (user !~ /^[A-Za-z0-9][A-Za-z0-9._-]*$/) {
        fail("user must start with a letter or number and contain only letters, numbers, dots, underscores, or dashes")
      }
      if (key !~ /^age1[0-9a-z]+$/) {
        fail("age recipient must look like age1...")
      }
      if (user in seen) {
        fail("duplicate user \"" user "\"")
      }
      seen[user] = 1
      print user "\t" key
    }
  ' "$RECIPIENTS_PATH" | sort -k1,1
}

write_file() {
  local path="$1"
  local mode="$2"
  local tmp="${path}.tmp"

  cat > "$tmp"
  chmod "$mode" "$tmp"
  mv "$tmp" "$path"
}

make_temp() {
  local var_name="$1"
  local path
  path="$(mktemp)"
  TMP_FILES+=("$path")
  printf -v "$var_name" '%s' "$path"
}

save_rows() {
  local rows="$1"
  local age_keys

  mkdir -p "$(dirname -- "$RECIPIENTS_PATH")"

  {
    printf '# Public age recipients allowed to decrypt development SOPS secrets.\n'
    printf '# Format: <user> <age-recipient>\n'
    printf '#\n'
    printf '# Create and register a local identity:\n'
    printf '#   mise run secrets:new-key -- alice\n'
    printf '#\n'
    printf '# Register a public recipient shared by another user:\n'
    printf '#   mise run secrets:add-key -- alice age1...\n'
    while IFS=$'\t' read -r user key; do
      [[ -n "${user:-}" ]] || continue
      printf '%s %s\n' "$user" "$key"
    done < "$rows"
  } | write_file "$RECIPIENTS_PATH" 0644

  {
    printf '# Managed by mise secrets:* tasks. Public user mappings live in .secrets/recipients.txt.\n'
    if [[ ! -s "$rows" ]]; then
      printf 'creation_rules: []\n'
    else
      age_keys="$(awk 'BEGIN { sep = "" } { printf "%s%s", sep, $2; sep = "," }' "$rows")"
      printf 'creation_rules:\n'
      printf "  - path_regex: '^\\\\.secrets/.*\\\\.sops\\\\.(env|json|yaml|yml)$'\n"
      printf '    age: %s\n' "$age_keys"
    fi
  } | write_file "$SOPS_CONFIG_PATH" 0644
}

sorted_rows_with_recipient() {
  local user="$1"
  local key="$2"
  local current_rows="$3"

  awk -v user="$user" -v key="$key" '
    BEGIN { updated = 0 }
    $1 == user {
      print user "\t" key
      updated = 1
      next
    }
    { print }
    END {
      if (!updated) {
        print user "\t" key
      }
    }
  ' "$current_rows" | sort -k1,1
}

has_user() {
  local user="$1"
  local rows="$2"

  awk -v user="$user" '$1 == user { found = 1 } END { exit found ? 0 : 1 }' "$rows"
}

print_update_keys_hint() {
  local matches=()
  shopt -s nullglob
  matches=(.secrets/*.sops.*)
  shopt -u nullglob

  if (( ${#matches[@]} == 0 )); then
    return 0
  fi

  printf 'Apply recipient changes to existing encrypted files:\n'
  printf '%s\n' "${matches[@]}" | sort | while IFS= read -r path; do
    printf '  sops updatekeys -y %s\n' "$path"
  done
}

add_key() {
  local user="$1"
  local key="$2"
  local current_rows updated_rows

  validate_user "$user"
  validate_age_recipient "$key"

  make_temp current_rows
  make_temp updated_rows

  recipient_rows > "$current_rows"
  if has_user "$user" "$current_rows"; then
    sorted_rows_with_recipient "$user" "$key" "$current_rows" > "$updated_rows"
    save_rows "$updated_rows"
    printf 'updated recipient for %s\n' "$user"
  else
    sorted_rows_with_recipient "$user" "$key" "$current_rows" > "$updated_rows"
    save_rows "$updated_rows"
    printf 'added recipient for %s\n' "$user"
  fi
  print_update_keys_hint
}

new_key() {
  local user="$1"
  local key_path key current_rows updated_rows

  validate_user "$user"

  make_temp current_rows
  make_temp updated_rows

  recipient_rows > "$current_rows"
  if has_user "$user" "$current_rows"; then
    die "user \"$user\" already has a registered recipient; use secrets:add-key to rotate it"
  fi

  mkdir -p "$KEY_DIR"
  chmod 0700 "$KEY_DIR"
  key_path="$KEY_DIR/$user.txt"
  if [[ -e "$key_path" ]]; then
    die "$key_path already exists"
  fi

  age-keygen -o "$key_path" >&2
  chmod 0600 "$key_path"

  key="$(age-keygen -y "$key_path")"
  validate_age_recipient "$key"

  sorted_rows_with_recipient "$user" "$key" "$current_rows" > "$updated_rows"
  save_rows "$updated_rows"

  printf 'created private age identity: %s\n' "$key_path"
  printf 'registered recipient for %s: %s\n' "$user" "$key"
  print_update_keys_hint
}

remove_user() {
  local user="$1"
  local current_rows updated_rows

  validate_user "$user"

  make_temp current_rows
  make_temp updated_rows

  recipient_rows > "$current_rows"
  if ! has_user "$user" "$current_rows"; then
    die "user \"$user\" is not in .secrets/recipients.txt"
  fi

  awk -v user="$user" '$1 != user { print }' "$current_rows" > "$updated_rows"
  save_rows "$updated_rows"

  printf 'removed recipient for %s\n' "$user"
  print_update_keys_hint
}

list_users() {
  local rows
  make_temp rows

  recipient_rows > "$rows"
  if [[ ! -s "$rows" ]]; then
    printf 'No users have access yet.\n'
    return 0
  fi

  while IFS=$'\t' read -r user key; do
    printf '%-20s %s\n' "$user" "$key"
  done < "$rows"
}

main() {
  if (( $# == 0 )); then
    usage
    die "missing command"
  fi

  case "$1" in
    new-key)
      (( $# == 2 )) || { usage; die "usage: new-key <user>"; }
      new_key "$2"
      ;;
    add-key)
      (( $# == 3 )) || { usage; die "usage: add-key <user> <age-recipient>"; }
      add_key "$2" "$3"
      ;;
    remove-user)
      (( $# == 2 )) || { usage; die "usage: remove-user <user>"; }
      remove_user "$2"
      ;;
    users)
      (( $# == 1 )) || { usage; die "usage: users"; }
      list_users
      ;;
    *)
      usage
      die "unknown command \"$1\""
      ;;
  esac
}

main "$@"
