{ pkgs ? import <nixpkgs> { } }:

pkgs.buildGoModule {
  pname = "go-time";
  version = "0.1.58"; # Replace with your desired version
  src = ./.;
  vendorHash = "sha256-aDCqZLQ1xoFGJFYbNa4TfXZJOIT0VT9CqK8nkf4+nPk=";
}

