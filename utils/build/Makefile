git_rev    = $(shell git rev-parse --short HEAD)
git_tag    = $(shell git describe --tags --abbrev=0)
git_branch = $(shell git rev-parse --abbrev-ref HEAD)
app_name   = "demo_app"

BuildVersion := $(git_branch)_$(git_rev)
BuildTime := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
BuildCommit := $(shell git rev-parse --short HEAD)
BuildGoVersion := $(shell go version)


# in detached HEAD state
ifeq ($(git_branch), HEAD)
   git_branch = $(shell git show-ref | grep $(shell git show HEAD | sed -n 1p | cut -d " " -f 2) | sed 's|.*/(.*)|\1|' | grep -v HEAD | sort | uniq | head -n 1)
   # when git checkout <<tag>>, branch may still be empty
   ifeq ($(git_branch), )
      git_branch := $(git_tag)
   endif
   BuildVersion := $(git_branch)_$(git_rev)
endif

ifeq ($(git_branch), develop)
   BuildVersion := develop_$(git_rev)
endif

ifeq ($(git_branch), master)
   BuildVersion := release_$(git_tag)_$(git_rev)
endif

# -ldflag 参数
GOLDFLAGS = -s -w -X 'main.BuildVersion=$(BuildVersion)'
GOLDFLAGS += -X 'main.BuildTime=$(BuildTime)'
GOLDFLAGS += -X 'main.BuildCommit=$(BuildCommit)'
GOLDFLAGS += -X 'main.BuildGoVersion=$(BuildGoVersion)'
