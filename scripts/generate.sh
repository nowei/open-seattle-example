#!/bin/bash

go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
if [[ ":$PATH:" != *":$HOME/go/bin/:"* ]]; then
    export PATH=$PATH:$HOME/go/bin/
fi
$(cd server/internal && go generate ./...)
