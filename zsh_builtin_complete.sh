_completion_remove_quote() {
  local input="$1"
  if [[ "$input" == \"* ]]; then
    echo "${input#\"}"
  else
    echo "${input#\'}"
  fi
}

# -----------------------------------
# uss completion wrapper for cd
# -----------------------------------
if [[ -z "$COMPLETION_USS_CD_FUNC" ]]; then
  COMPLETION_USS_CD_FUNC=$_comps[cd]
fi

_completion_uss_cd() {
  local word_len=${#words[@]}
  local last_word="${words[$word_len]}"
  local last_word_trim="$(_completion_remove_quote "$last_word")"

  if [[ "$last_word_trim" == "uss"* ]]; then
    local -a results
    results=("${(@f)$(GO_FLAGS_COMPLETE_URI=1 GO_FLAGS_SHELL=zsh uss "$last_word")}")
    local -a ns_items regular_items
    for item in "${results[@]}"; do
        if [[ "$item" == *"<NS>" ]]; then
            ns_items+=("${item%<NS>}")
        else
            regular_items+=("$item")
        fi
    done
    (( ${#ns_items} > 0 )) && compadd -S '' -Q -a ns_items
    (( ${#regular_items} > 0 )) && compadd -Q -a regular_items
  else
    # call the underlying completion
    $COMPLETION_USS_CD_FUNC
  fi
}

compdef _completion_uss_cd cd

# -----------------------------------
# uss completion wrapper for cat
# -----------------------------------
if [[ -z "$COMPLETION_USS_CAT_FUNC" ]]; then
  COMPLETION_USS_CAT_FUNC=$_comps[cat]
fi

_completion_uss_cat() {
  local word_len=${#words[@]}
  local last_word="${words[$word_len]}"
  local last_word_trim="$(_completion_remove_quote "$last_word")"

  if [[ "$last_word_trim" == "uss"* ]]; then
    local -a results
    results=("${(@f)$(GO_FLAGS_COMPLETE_URI=1 GO_FLAGS_SHELL=zsh uss "$last_word")}")
    local -a ns_items regular_items
    for item in "${results[@]}"; do
        if [[ "$item" == *"<NS>" ]]; then
            ns_items+=("${item%<NS>}")
        else
            regular_items+=("$item")
        fi
    done
    (( ${#ns_items} > 0 )) && compadd -S '' -Q -a ns_items
    (( ${#regular_items} > 0 )) && compadd -Q -a regular_items
  else
    # call the underlying completion
    $COMPLETION_USS_CAT_FUNC
  fi
}

compdef _completion_uss_cat cat

# -----------------------------------
# uss completion wrapper for ls
# -----------------------------------
if [[ -z "$COMPLETION_USS_LS_FUNC" ]]; then
  COMPLETION_USS_LS_FUNC=$_comps[ls]
fi

_completion_uss_ls() {
  local word_len=${#words[@]}
  local last_word="${words[$word_len]}"
  local last_word_trim="$(_completion_remove_quote "$last_word")"

  if [[ "$last_word_trim" == "uss"* ]]; then
    local -a results
    results=("${(@f)$(GO_FLAGS_COMPLETE_URI=1 GO_FLAGS_SHELL=zsh uss "$last_word")}")
    local -a ns_items regular_items
    for item in "${results[@]}"; do
        if [[ "$item" == *"<NS>" ]]; then
            ns_items+=("${item%<NS>}")
        else
            regular_items+=("$item")
        fi
    done
    (( ${#ns_items} > 0 )) && compadd -S '' -Q -a ns_items
    (( ${#regular_items} > 0 )) && compadd -Q -a regular_items
  else
    # call the underlying completion
    $COMPLETION_USS_LS_FUNC
  fi
}

compdef _completion_uss_ls ls
