_get_completion_func() {
  local defs=( $(complete -p $1) )
  local defs_len=${#defs[@]}
  local function_name=${defs[$defs_len - 2]}
  echo $function_name
}

# -----------------------------------
# uss completion wrapper for cd
# -----------------------------------

COMPLETION_USS_CD_FUNC=$(_get_completion_func cd)

_completion_uss_cd() {
  local word_len=${#COMP_WORDS[@]}
  local last_word="${COMP_WORDS[$word_len - 1]}"
  local last_word_trim=${last_word#"\""}

  if [[ "$last_word_trim" == "uss"* ]]; then
    local IFS=$'\n'
    COMPREPLY=($(GO_FLAGS_COMPLETE_URI=1 uss "$last_word"))
  else
    # call the underlying completion
    $COMPLETION_USS_CD_FUNC
  fi
}

if [[ "$COMPLETION_USS_CD_FUNC" != "_completion_uss_cd" ]]; then
  complete -o nospace -F _completion_uss_cd cd
fi

# -----------------------------------
# uss completion wrapper for cat
# -----------------------------------

COMPLETION_USS_CAT_FUNC=$(_get_completion_func cat)

_completion_uss_cat() {
  local word_len=${#COMP_WORDS[@]}
  local last_word="${COMP_WORDS[$word_len - 1]}"
  local last_word_trim=${last_word#"\""}

  if [[ "$last_word_trim" == "uss"* ]]; then
    local IFS=$'\n'
    compopt -o nospace
    COMPREPLY=($(GO_FLAGS_COMPLETE_URI=1 uss "$last_word"))
  else
    # call the underlying completion
    $COMPLETION_USS_CAT_FUNC
  fi
}

if [[ "$COMPLETION_USS_CAT_FUNC" != "_completion_uss_cat" ]]; then
  complete -F _completion_uss_cat cat
fi
