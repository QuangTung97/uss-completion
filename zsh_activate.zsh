export PATH="$PATH:."

_zsh_add_completion_item() {
  local -a ns_items regular_items
  for item in "$@"; do
    if [[ "$item" == *"<NS>" ]]; then
      ns_items+=("${item%<NS>}")
    else
      regular_items+=("$item")
    fi
  done

  (( ${#ns_items} > 0 )) && compadd -S '' -Q -a ns_items
  (( ${#regular_items} > 0 )) && compadd -Q -a regular_items
}

_uss_completion() {
  local -a args=("${words[@]:1:$CURRENT-1}")
  local -a results=("${(@f)$(GO_FLAGS_COMPLETION=1 GO_FLAGS_SHELL=zsh "${words[1]}" "${args[@]}")}")
  _zsh_add_completion_item "${results[@]}"
}

compdef _uss_completion uss
