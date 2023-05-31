{
  description = "neocode flake";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        name = "carapace";
        package = (with pkgs; (makeOverridable callPackage self { }));
      in
      {
        defaultPackage = package;
        packages.${name} = package;

        devShell =
          pkgs.mkShell {
            buildInputs = with pkgs; [
              gh
              go
              gopls
            ];
          };
      });
}
