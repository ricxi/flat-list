#!/bin/bash
# If you don't have the complete command on your linux distro, 
# the auto-complete feature for this script isn't going to work.

# give access to functions in post_req
alias call_user="./dev_scripts/post_req.sh"

# Used to import private go packages for this project (make sure to set up .netrc or .gitconfig)
export GOPRIVATE=github.com/ricxi/flat-list

# Runs go mod tidy.
# Call this in a go module directory, or
# pass it the path to a go module directory.
tidy() {
    local go_module_dir="$1"
    # I should check if this is empty

    if [ "$#" -eq 0 ]; then
        [ -f 'go.mod' ] && {
            go mod tidy -v &&
            return 0 || return 1
        } ||
        echo 'Error: call this inside a go module or pass a path to a go module as an argument: tidy <path_to_go_module>' &&
        return 1
    fi

    pushd "$go_module_dir" 1> /dev/null && {
        go mod tidy -v 
        popd 1> /dev/null
    }
}

# autocomplete setup for 'tidy' function
_tidy_completions() {
    local go_dirs=()

    for go_dir in */; do
        [ -f "${go_dir}go.mod" ] && go_dirs+=("${go_dir%/}")
    done

    COMPREPLY=("$(compgen -W "${go_dirs[*]}" "${COMP_WORDS[1]}")")
}
complete -F _tidy_completions tidy
