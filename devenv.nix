{
  pkgs,
  lib,
  config,
  ...
}:

let
  # Lifited from devenv's `languages.go` module.
  # Override the buildGoModule function to use the specified Go package.
  buildGoModule = pkgs.buildGoModule.override { go = config.languages.go.package; };
  buildWithSpecificGo = pkg: pkg.override { inherit buildGoModule; };
in
{
  languages.go = {
    enable = true;
    enableHardeningWorkaround = true;
  };

  packages = [
    pkgs.treefmt
    pkgs.reuse
    (buildWithSpecificGo pkgs.xcaddy)
    (buildWithSpecificGo pkgs.pkgsite)
  ];

  pre-commit.hooks = {
    treefmt = {
      enable = true;
      settings.formatters = [
        config.languages.go.package
        pkgs.nixfmt-rfc-style
        pkgs.toml-sort
        pkgs.typos
      ];
    };
    reuse.enable = true;
  };

  difftastic.enable = true;
}
