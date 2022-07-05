{
  description = "Nhost Hasura Storage";

  inputs = {
    nix-common.url = "github:dbarrosop/nix-common/dbarroso/hasura-storage";
    nixpkgs = {
      inputs.nixpkgs.follows = "nix-common/nixpkgs";
    };
    flake-utils.url = "github:numtide/flake-utils";
    nix-filter.url = "github:numtide/nix-filter";
  };

  outputs = { self, nix-common, nixpkgs, flake-utils, nix-filter }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        name = "hasura-storage";
        version = pkgs.lib.fileContents ./VERSION;
        module = "github.com/nhost/hasura-storage";

        pkgs = import nixpkgs {
          inherit system;
          overlays = [
            nix-common.overlays.default
            (final: prev: rec {
              vips = prev.vips.overrideAttrs (oldAttrs: rec {
                buildInputs = [
                  final.glib
                  final.libxml2
                  final.expat
                  final.libjpeg
                  final.libpng
                  final.libwebp
                  final.openjpeg
                ];
              });
            })
          ];
        };

        go-src = nix-filter.lib.filter {
          root = ./.;
        };

        nix-src = nix-filter.lib.filter {
          root = ./.;
          include = [
            (nix-filter.lib.matchExt "nix")
          ];
        };

        buildInputs = with pkgs; [
          vips
        ];

        nativeBuildInputs = with pkgs; [
          go
          clang
          pkg-config
        ];

        checkBuildInputs = with pkgs; [
          docker-client
        ];


        ldflags = [
          "-X ${module}/controller.buildVersion=${version}"
        ];

        tags = [ "integration" ];

      in
      {
        checks = nix-common.checks.go {
          inherit pkgs buildInputs nativeBuildInputs checkBuildInputs tags ldflags go-src nix-src;

          preCheck = ''
            export HASURA_AUTH_BEARER=$(make dev-jwt)
            export TEST_S3_ACCESS_KEY=$(make dev-s3-access-key)
            export TEST_S3_SECRET_KEY=$(make dev-s3-secret-key)
            export GIN_MODE=release
          '';
        };

        devShells = flake-utils.lib.flattenTree rec {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [
              nixpkgs-fmt
              golangci-lint
              docker-client
              docker-compose
              go-migrate
              gnumake
              gnused
              richgo
              ccls
            ] ++ buildInputs ++ nativeBuildInputs;
          };
        };

        packages = flake-utils.lib.flattenTree
          rec {
            hasura-storage = pkgs.buildGoModule {
              inherit version ldflags buildInputs nativeBuildInputs;

              pname = name;

              src = go-src;

              vendorSha256 = null;

              doCheck = false;

              subPackages = [ "." ];

              meta = with pkgs.lib; {
                description = "Hasura Storage is awesome";
                homepage = "https://github.com/nhost/hasura-storage";
                license = licenses.mit;
                maintainers = [ "nhost" ];
                platforms = platforms.linux ++ platforms.darwin;
              };
            };

            docker-image = pkgs.dockerTools.buildLayeredImage {
              inherit name;
              tag = version;
              created = "now";

              contents = [
                pkgs.cacert
              ] ++ buildInputs;
              config = {
                Env = [
                  "TMPDIR=/"
                  "MALLOC_ARENA_MAX=2"
                ];
                Entrypoint = [
                  "${self.packages.${system}.hasura-storage}/bin/hasura-storage"
                ];
              };
            };

            default = hasura-storage;

          };
      }
    );
}
