{
  description = "Go Development Environment";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in
    {
      devShells.${system}.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          gotools
          golangci-lint
          gopls
        ];

        shellHook = ''
          export GOPATH="$PWD/.nix-go"
        '';
      };
    };
}
