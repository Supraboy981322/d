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

              # server   TODO: rewrite in Zig
              brotli
              pkg-config

              # proto native desktop gui program
              mesa
              libXi
              libXcursor
              libXrandr
              libglvnd
              libXinerama
              wayland
              libxkbcommon

              #disgusting fat electron app (this is some odious bloat)
              #  HAVE JS DEVELOPERS NEVER HEARD OF STATICALLY LINKING?
              nss
              atk
              nspr
              dbus
              cups #WHY THE HELL DOES ELECTRON NEED CUPS? (who's printing from an electron app?)
              glib
              gtk3
              libGL
              cairo
              pango
              expat
              libgbm
              libx11
              libxcb
              libxext
              alsa-lib
              electron
              libxfixes
              libxrandr
              libxdamage
              at-spi2-atk
              nodejs-slim
              libxkbcommon
              libxcomposite

              # build tools
              bun
            ];
          in { 
          buildInputs = libs;
          packages = libs;
          shellHook = /* bash */ ''
            #who said arrays in Bash are bad? 
            #  NOTE: DAMN, electron SUCKS
            electron_libs=(
              "${pkgs.nspr.out}"
              "${pkgs.nss.out}"
              "${pkgs.glib.out}"
              "${pkgs.atk.out}"
              "${pkgs.at-spi2-atk.out}"
              "${pkgs.cups.lib}"
              "${pkgs.dbus.out}"
              "${pkgs.dbus.lib}"
              "${pkgs.cairo.out}"
              "${pkgs.gtk3.out}"
              "${pkgs.pango.out}"
              "${pkgs.libx11.out}"
              "${pkgs.libGL.out}"
              "${pkgs.alsa-lib.out}"
              "${pkgs.libxcomposite.out}"
              "${pkgs.libxdamage.out}"
              "${pkgs.libxext.out}"
              "${pkgs.libxfixes.out}"
              "${pkgs.libxrandr.out}"
              "${pkgs.libgbm.out}"
              "${pkgs.expat.out}"
              "${pkgs.libxcb.out}"
              "${pkgs.libxkbcommon.out}"
            )
            for lib in "''${electron_libs[@]}"; do
              export LD_LIBRARY_PATH+=":$lib/lib"
              export PATH+=":$lib/bin"
            done

            REPO_ROOT="$(git rev-parse --show-toplevel)"
            build() (
              set -eou pipefail
              old_dir="$PWD"
              cd "$REPO_ROOT/src/dServer/web"
              bun run build.ts 2>&1 | sed 's/^/\t/g' 
              cd ..
              printf "uncompressed size: %d\n" "$(cat web_built/*.html | wc -c)"
              printf "gzipped size: %d\n" "$(cat web_built/*.html | gzip -c9k | wc -c)"
              printf "brotli size: %d\n" "$(cat web_built/*.html | brotli -cq 11 -w 24 | wc -c)"
              printf "building server"
              go build -x -v 2>&1 | sed 's/^/\t/g'
              cd "$REPO_ROOT/src/disgusting_electron_app"
              printf "building the disgusting electron app"
              npm run dist 2>&1 | sed 's/^/\t/g'
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
