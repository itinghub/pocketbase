echo "start build"
cd "$(dirname "$0")"
GO_ENABLED=0 go build -v
echo "build done"
