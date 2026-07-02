{
  inputs = {
    nixpkgs.follows = "synergy/nixpkgs";
    synergy.url = "github:krostar/synergy";
  };

  outputs = {synergy, ...} @ inputs:
    synergy.lib.mkFlake {
      inherit inputs;
      src = ./.nix;
      eval.synergy.restrictDependenciesUnits.synergy = ["harmony" "krostar"];
    };
}
