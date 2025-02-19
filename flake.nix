{
  description = "Flake for Caddy plugin development";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

    pre-commit-hooks.url = "github:cachix/git-hooks.nix";

    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    { self, nixpkgs, ... }@inputs:
    let
      supportedSystems = [
        "x86_64-linux"
        "aarch64-linux"
      ];

      perSystem =
        system:
        let
          pkgs = import nixpkgs {
            inherit system;
            config = { };
            overlays = [ ];
          };

          goPackage = pkgs.go;
          buildGoModule = pkgs.buildGoModule.override { go = goPackage; };
          buildWithSpecificGo = pkg: pkg.override { inherit buildGoModule; };
        in
        nixpkgs.lib.fix (_self: {
          checks = {
            pre-commit-check = inputs.pre-commit-hooks.lib.${system}.run {
              src = ./.;
              hooks = {
                treefmt = {
                  enable = true;
                  settings.formatters = [
                    goPackage
                    pkgs.nixfmt-rfc-style
                    pkgs.typos
                    pkgs.toml-sort
                  ];
                };
                reuse.enable = true;
              };
            };
          };

          devShells.default = pkgs.mkShell {
            inherit (_self.checks.pre-commit-check) shellHook;
            nativeBuildInputs = [
              goPackage
              (buildWithSpecificGo pkgs.xcaddy)
              (buildWithSpecificGo pkgs.pkgsite)
            ] ++ _self.checks.pre-commit-check.enabledPackages;
          };
        });
    in
    inputs.flake-utils.lib.eachSystem supportedSystems perSystem;
}
