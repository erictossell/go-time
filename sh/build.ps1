$env:CGO_ENABLED="1"
$env:GOARCH="amd64"
$env:CC="zig cc"
go build -o go-time.exe