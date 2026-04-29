{
  description = "Rust development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = [
            pkgs.gcc
            pkgs.libgcc
            pkgs.cargo
            pkgs.rustc
            pkgs.rustfmt
            pkgs.rustPlatform.rustLibSrc
            pkgs.rust-analyzer
            pkgs.pkg-config
          ];
        };
      }
    );
}