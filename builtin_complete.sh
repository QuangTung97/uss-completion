_get_completion_func() {
  local defs=( $(complete -p $1) )
  local defs_len=${#defs[@]}
  local function_name=${defs[$defs_len - 2]}
  echo $function_name
}

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
  COMPLETION_USS_CD_FUNC=$(_get_completion_func cd)
fi

_completion_uss_cd() {
  local last_word="${COMP_WORDS[$COMP_CWORD]}"
  local last_word_trim="$(_completion_remove_quote "$last_word")"

  if [[ "$last_word_trim" == "uss"* ]]; then
    local IFS=$'\n'
    
    COMPREPLY=($(GO_FLAGS_COMPLETE_URI=1 uss "$last_word"))
    return 0
  else
    # call the underlying completion
    $COMPLETION_USS_CD_FUNC
  fi
}

complete -o nospace -F _completion_uss_cd cd

# -----------------------------------
# uss completion wrapper for cat
# -----------------------------------
if [[ -z "$COMPLETION_USS_CAT_FUNC" ]]; then
  COMPLETION_USS_CAT_FUNC=$(_get_completion_func cat)
fi

_completion_uss_cat() {
  local last_word="${COMP_WORDS[$COMP_CWORD]}"
  local last_word_trim="$(_completion_remove_quote "$last_word")"

  if [[ "$last_word_trim" == "uss"* ]]; then
    local IFS=$'\n'
    compopt -o nospace
    COMPREPLY=($(GO_FLAGS_COMPLETE_URI=1 uss "$last_word"))
    return 0
  else
    # call the underlying completion
    $COMPLETION_USS_CAT_FUNC
  fi
}

complete -F _completion_uss_cat cat

# -----------------------------------
# uss completion wrapper for ls
# -----------------------------------
if [[ -z "$COMPLETION_USS_LS_FUNC" ]]; then
  COMPLETION_USS_LS_FUNC=$(_get_completion_func ls)
fi

_completion_uss_ls() {
  local last_word="${COMP_WORDS[$COMP_CWORD]}"
  local last_word_trim="$(_completion_remove_quote "$last_word")"

  if [[ "$last_word_trim" == "uss"* ]]; then
    local IFS=$'\n'
    compopt -o nospace
    COMPREPLY=($(GO_FLAGS_COMPLETE_URI=1 uss "$last_word"))
    return 0
  else
    # call the underlying completion
    $COMPLETION_USS_LS_FUNC
  fi
}

complete -F _completion_uss_ls ls
