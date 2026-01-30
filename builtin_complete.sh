_get_completion_func() {
  local defs=( $(complete -p $1) )
  local defs_len=${#defs[@]}
  local function_name=${defs[$defs_len - 2]}
  echo $function_name
}

COMPLETION_CD_FUNC=$(_get_completion_func cd)

_completion_uss_cd() {
  local word_len=${#COMP_WORDS[@]}
  local last_word="${COMP_WORDS[$word_len - 1]}"
  local last_word_trim=${last_word#"\""}

  echo "WORD:" $last_word_trim >> log.txt

  if [[ "$last_word_trim" == "uss"* ]]; then
    GO_FLAGS_COMPLETE_URI=1 uss "$last_word"
  else
    # call the underlying completion
    $COMPLETION_CD_FUNC
  fi
}

if [[ "$COMPLETION_CD_FUNC" != "_completion_uss_cd" ]]; then
  complete -o nospace -F _completion_uss_cd cd
fi
