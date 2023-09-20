#!/bin/bash
folder=$(pwd)
builtOn="$(USER)@$(hostname)"
builtAt="$(date +'%F %T %z')"
goVersion=$(go version | sed 's/go version //')
gitAuthor=$(git show -s --format='format:%aN <%ae>' HEAD)
gitCommit=$(git rev-parse HEAD)
gitTag=$( git -C "$folder" describe --tags HEAD|| git -C "$folder" rev-parse --abbrev-ref HEAD | grep -v HEAD || git -C "$folder" rev-parse HEAD )
ldflags="\
-X 'go-lib/utils/build/inject.builtOn=$builtOn' \
-X 'go-lib/utils/build/inject.builtAt=$builtAt' \
-X 'go-lib/utils/build/inject.goVersion=$goVersion' \
-X 'go-lib/utils/build/inject.gitAuthor=$gitAuthor' \
-X 'go-lib/utils/build/inject.gitCommit=$gitCommit' \
-X 'go-lib/utils/build/inject.gitTag=$gitTag' \
"

go build -o inject.exe -ldflags "$ldflags" main.go && ./inject.exe -version