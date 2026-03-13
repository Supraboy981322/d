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
          vendorSha256 = pkgs.lib.fakeHash;
        };

        client = pkgs.buildGoModule {
          pname = "d";
          version = "0.0.0";
          src = ./src/dServer;
          vendorSha256 = pkgs.lib.fakeHash;

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

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            brotli
            pkg-config
          ];
        };
      })
    ) // {
      nixosModules.default = { config, lib, pkgs, ... }: {
        options.services.dServer = {
          enable = lib.mkEnableOption "dServer";
          package = lib.mkOption {
            type = lib.types.package;
            default = self.packages.${pkgs.system}.server;
          };
        };
        config = lib.mkIf config.services.dServer.enable { 
          systemd.services.dServer = {
            description = "dServer instance";
            wantedBy = [ "multi-user.target" ];
            serviceConfig = {
              ExecStart = "${config.services.dServer.package}/bin/server";
              Restart = "always";
            };
          };
        };
      };
    };
}
