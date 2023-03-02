#!/bin/bash

CUR_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd "$CUR_DIR"

OPT="paths=import,module=example"

# gen for HTTP service
for f in $( find . -type f | grep '.proto' ); do
    echo compiling: $f
    protoc -I ./protospec \
        --go_out=. \
        --go_opt=$OPT \
        --go-echo-http_out=. \
        --go-echo-http_opt=$OPT \
        $f || exit
done

echo "ok"
