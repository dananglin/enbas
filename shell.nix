# Creates a development environment for the Enbas project.
let
  commit_ref = "55d1f923c480dadce40f5231feb472e81b0bab48";
  nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/tarball/${commit_ref}";
  pkgs = import nixpkgs {
    config = { };
    overlays = [ ];
  };
in

pkgs.mkShellNoCC {
  packages = with pkgs; [
    delve
    git
    go
    go-grip
    golangci-lint
    gopls
    mage
    tmux
  ];

  shellHook = ''
    export GOROOT=$( which go | xargs dirname | xargs dirname )/share/go
    exec tmux new-session -s "Enbas Development"
  '';
}
