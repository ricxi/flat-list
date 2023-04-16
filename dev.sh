#!/bin/bash
# used to import private go packages for this project (make sure to set up .netrc or .gitconfig)
export GOPRIVATE=github.com/ricxi/flat-list

# api calls to test user
alias post_user='./dev_scripts/post_req.sh'

# Runs go mod tidy.
# Call this in a go module directory, or
# pass it the path to a go module directory.
tidy() {
    local go_module_dir="$1"

    if [ "$#" -eq 0 ]; then
        [ -f 'go.mod' ] && {
            go mod tidy -v &&
            return 0 || return 1
        } ||
        echo 'Error: call inside a go module or pass a path to go module as an argument: tidy <path_to_go_module>' &&
        return 1
    fi

    pushd "$go_module_dir" 1> /dev/null && {
        go mod tidy -v 
        popd 1> /dev/null
    }
}