{ pkgs ? import <nixpkgs> { } }:

pkgs.mkShell {
  buildInputs = with pkgs; [ go gopls gotools go-tools golangci-lint sqlitebrowser zig ];
}

