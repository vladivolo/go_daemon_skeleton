BINARY = daemon
SOURCEDIR = .
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')
VERSION = 1.0.0
MAINTAINER = you@email.org

all: $(SOURCES) version
	go build -o $(BINARY) main.go server.go service.go version.go

version:
	utils/generate_version_info.sh "$(BINARY)" "$(VERSION)" "$(MAINTAINER)" "$(AUTO_BUILD_TAG)" "$(MAKE) $(MAKEFLAGS)" > version.go

clean:
	rm -f $(BINARY) version.go *~ */*~

tags:
	etags *.go

