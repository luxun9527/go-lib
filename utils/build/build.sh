#!/bin/bash
builtOn="$USER@$(hostname)"
builtAt="$(date +'%F %T %z')"
goVersion=$(go version | sed 's/go version //')
gitAuthor=$(git show -s --format='format:%aN <%ae>' HEAD)
gitCommit=$(git rev-parse HEAD)

ldflags="\
-X 'go-lib/utils/build/inject.builtOn=$builtOn' \
-X 'go-lib/utils/build/inject.builtAt=$builtAt' \
-X 'go-lib/utils/build/inject.goVersion=$goVersion' \
-X 'go-lib/utils/build/inject.gitAuthor=$gitAuthor' \
-X 'go-lib/utils/build/inject.gitCommit=$gitCommit' \
"

go build -o inject.exe -ldflags "$ldflags" main.go