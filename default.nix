{ pkgs ? import <nixpkgs> { } }:

pkgs.buildGoModule {
  pname = "go-time";
  version = "0.1.58"; # Replace with your desired version
  src = ./.;
  vendorHash = "sha256-RdEHYqMke/gOECqKntiuoxN1NajDgSN1MYIsjOvmW9c=";
}

