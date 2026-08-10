{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-26.05";
    nixpkgs-unstable.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs = { self, nixpkgs, nixpkgs-unstable }:
  let
    system = "x86_64-linux";
    pkgs = import nixpkgs { inherit system; };
    pkgs' = import nixpkgs-unstable { inherit system; };
  in
  {
    devShells.${system}.default = pkgs.mkShell {
      shellHook = ''
        set -a
        source env.sh
        set +a

        export LD_LIBRARY_PATH="$LD_LIBRARY_PATH:${pkgs.lib.makeLibraryPath [ pkgs'.sqlite ]}"
      '';

      CGO_ENABLED = 1;

      buildInputs = with pkgs; [
        go
        gopls
        watchexec
        gcc
        just
        pkgs'.sqlite
        pkgs'.sqitchSqlite
      ];
    };
  };
}
