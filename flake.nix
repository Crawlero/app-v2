{
  description = "Development environment for Golang";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShell = pkgs.mkShell {
          buildInputs = [
            pkgs.go
            pkgs.go-migrate
            pkgs.air
            pkgs.tailwindcss
          ];

          shellHook = ''
            echo "Hello from the Golang development environment!"
          '';
        };
      });
}
