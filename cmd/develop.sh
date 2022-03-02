set -ex

if [ ! -d ./bin ]; then
    mkdir bin
fi

if [ ! -z "$1" ]
then
    go build -o bin github.com/mohanson/rvv-daily-benchmark/cmd/$1
fi
