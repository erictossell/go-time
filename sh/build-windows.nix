{ pkgs ? import <nixpkgs> { } }:

pkgs.mkShell {
  buildInputs = with pkgs; [ go zig ];
  shellHook = ''
    export GOOS=windows
    echo "Building for $GOOS"
    export GOARCH=amd64
    echo "Building for $GOARCH"
    export CGO_ENABLED=1
    echo "Building with CGO_ENABLED=$CGO_ENABLED"
    export CC="zig cc -target x86_64-windows-gnu"
    echo "Building with CC=$CC"
    export CXX="zig c++ -target x86_64-windows-gnu"
    echo "Building with CXX=$CXX"
    env GOOS=windows GOARCH=amd64 go build -o go-time.exe .
  '';
}
