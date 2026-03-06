{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable-small";
    synergy = {
      url = "github:krostar/synergy";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = {synergy, ...} @ inputs:
    synergy.lib.mkFlake {
      inherit inputs;
      src = ./.nix;
      eval.synergy.restrictDependenciesUnits.synergy = ["harmony"];
    };
}
