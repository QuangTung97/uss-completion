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
  local last_word="${words[$CURRENT]}"
  local last_word_trim="$(_completion_remove_quote "$last_word")"

  if [[ "$last_word_trim" == "uss"* ]]; then
    local -a results
    results=("${(@f)$(GO_FLAGS_COMPLETE_URI=1 GO_FLAGS_SHELL=zsh uss "$last_word")}")
    _zsh_add_completion_item "${results[@]}"
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
  local last_word="${words[$CURRENT]}"
  local last_word_trim="$(_completion_remove_quote "$last_word")"

  if [[ "$last_word_trim" == "uss"* ]]; then
    local -a results
    results=("${(@f)$(GO_FLAGS_COMPLETE_URI=1 GO_FLAGS_SHELL=zsh uss "$last_word")}")
    _zsh_add_completion_item "${results[@]}"
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
  local last_word="${words[$CURRENT]}"
  local last_word_trim="$(_completion_remove_quote "$last_word")"

  if [[ "$last_word_trim" == "uss"* ]]; then
    local -a results
    results=("${(@f)$(GO_FLAGS_COMPLETE_URI=1 GO_FLAGS_SHELL=zsh uss "$last_word")}")
    _zsh_add_completion_item "${results[@]}"
  else
    # call the underlying completion
    $COMPLETION_USS_LS_FUNC
  fi
}

compdef _completion_uss_ls ls
