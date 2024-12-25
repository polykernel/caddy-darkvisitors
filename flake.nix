{
  description = "Flake for Caddy plugin development";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

    devenv.url = "github:cachix/devenv";
    devenv.inputs.nixpkgs.follows = "nixpkgs";

    flake-utils.url = "github:numtide/flake-utils";
  };

  nixConfig = {
    extra-trusted-public-keys = "devenv.cachix.org-1:w1cLUi8dv3hnoSPGAuibQv+f9TZLr6cv/Hm9XgU50cw=";
    extra-substituters = "https://devenv.cachix.org";
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
        in
        nixpkgs.lib.fix (_self: {
          packages = {
            devenv-up = _self.devShells.default.config.procfileScript;
            devenv-test = _self.devShells.default.config.test;
          };

          devShells = {
            default = inputs.devenv.lib.mkShell {
              inherit inputs pkgs;
              modules = [
                (
                  { lib, ... }:
                  {
                    devenv.root =
                      let
                        devenvRootFileContent = builtins.readFile ./devenv_root;
                      in
                      pkgs.lib.mkIf (devenvRootFileContent != "") devenvRootFileContent;
                  }
                )

                ./devenv.nix
              ];
            };
          };
        });
    in
    inputs.flake-utils.lib.eachSystem supportedSystems perSystem;
}
