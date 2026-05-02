{
  description = "d_gui";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

    flake-utils.url = "github:numtide/flake-utils";

    # import Zig overlay
    zig_overlay = {
      url = "github:mitchellh/zig-overlay";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };
  outputs = { self, nixpkgs, flake-utils, zig_overlay }:
    (flake-utils.lib.eachDefaultSystem (system:
      let
        zigVersion = "0.16.0";

        # selected Zig package
        zig = zig_overlay.packages.${system}.${zigVersion};

        # add the Zig overlay pkgs
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ zig_overlay.overlays.default ];
        };

        server = import ./server.nix;
      in {
        devShells.default = pkgs.mkShell ({
          shellHook = ''
            export LD_LIBRARY_PATH=${
              pkgs.lib.makeLibraryPath (with pkgs; [
                mesa
                libXi
                libXcursor
                libXrandr
                libglvnd
                libXinerama
                wayland
                libxkbcommon
            ])}:$LD_LIBRARY_PATH
          '';
          packages = (with pkgs; [
            # raylib deps
            mesa
            libXi
            libXcursor
            libXrandr
            libglvnd
            libXinerama
            wayland
            libxkbcommon
          ]) ++ [ zig ];
        });
      })
    );
}
