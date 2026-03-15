{
  description = "a... server thing (not sure what to call it)";
  
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    (flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        server = pkgs.buildGoModule {
          pname = "dServer";
          version = "0.0.0";
          src = ./src/dServer;
          vendorHash = pkgs.lib.fakeHash;
        };

        client = pkgs.buildGoModule {
          pname = "d";
          version = "0.0.0";
          src = ./src/dServer;
          vendorHash = pkgs.lib.fakeHash;

          nativeBuildInputs = with pkgs; [
            pkg-config
          ];
          buildInputs = with pkgs; [
            brotli
          ];
        };
      in {
        packages = {
          inherit server client;
          default = client;
        };

        devShells.default = pkgs.mkShell (
          let
            libs = with pkgs; [
              go

              # server
              brotli
              pkg-config

              # proto desktop gui program
              mesa
              libXi
              libXcursor
              libXrandr
              libglvnd
              libXinerama
              wayland
              libxkbcommon

              # build tools
              bun
            ];
          in { 
          buildInputs = libs;
          packages = libs;
          shellHook = ''
            REPO_ROOT="$(git rev-parse --show-toplevel)"
            build() (
              set -eou pipefail
              old_dir="$PWD"
              cd "$REPO_ROOT/src/dServer/web"
              bun run build.ts || exit
              cd ..
              printf "uncompressed size: %d\n" "$(cat web_built/*.html | wc -c)"
              printf "gzipped size: %d\n" "$(cat web_built/*.html | gzip -c9k | wc -c)"
              printf "brotli size: %d\n" "$(cat web_built/*.html | brotli -cq 11 -w 24 | wc -c)"
              go build || true
              cd $old_dir
            )
            run() (
              set -eou pipefail
              old_dir="$PWD"
              cd "$REPO_ROOT/src/dServer/web"
              bun run build.ts || exit
              cd ..
              printf "uncompressed size: %d\n" "$(cat web_built/*.html | wc -c)"
              printf "gzipped size: %d\n" "$(cat web_built/*.html | gzip -c9k | wc -c)"
              printf "brotli size: %d\n" "$(cat web_built/*.html | brotli -cq 11 -w 24 | wc -c)"
              go run . || true
              cd $old_dir
            )
          '';
        });
      })
    ) // {
      nixosModules.default = { config, lib, ... }: {
        options.services.dServer = {
          enable = lib.mkEnableOption "dServer";
          package = lib.mkOption {
            type = lib.types.package;
            description = "dServer server package";
          };
        };
        config = lib.mkIf config.services.dServer.enable { 
          services.dServer.package = lib.mkDefault self.packages.${config.nixpkgs.system}.server;
          systemd.services.dServer = {
            description = "dServer instance";
            wantedBy = [ "multi-user.target" ];
            serviceConfig = {
              ExecStart = "${config.services.dServer.package}/bin/dServer";
              Restart = "always";
            };
          };
        };
      };
    };
}
