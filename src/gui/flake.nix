{
  description = "d_gui";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };
  outputs = { self, nixpkgs, flake-utils }:
    (flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in {
        devShells.default = pkgs.mkShell ({
          packages = with pkgs; [
            # raylib deps
            mesa
            libXi
            libXcursor
            libXrandr
            libglvnd
            libXinerama
            wayland
            libxkbcommon
          ];
        });
      })
    );
}
