#!/bin/bash
PRECOMMIT_LINT="golangci-lint run --new-from-rev HEAD --whole-files && (go mod tidy; git diff --exit-code --quiet -- go.*)"
GENERAL_LINT="golangci-lint run && (go mod tidy; git diff-index --quiet HEAD -- go.*)"

lint_dirs=$(find ~+ -name "go.mod")
result=0
for go_mod_file in ${lint_dirs[@]}; do
    directory=$(dirname $go_mod_file)
    cd $directory
    if [[ $1 = "--pre-commit" ]]
    then
        git diff --quiet HEAD $REF -- . || (echo "linting" $directory && eval "$PRECOMMIT_LINT")
    else
        (echo "linting" $directory && eval "$GENERAL_LINT")
    fi
    ((result=$?))
done

exit $result
