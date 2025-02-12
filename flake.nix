{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    harmony = {
      url = "git+file:///Users/alexis.destrez/Projects/Private/harmony";
      inputs = {
        synergy.follows = "synergy";
        nixpkgs.follows = "nixpkgs";
      };
    };
    synergy = {
      url = "git+file:///Users/alexis.destrez/Projects/Private/synergy";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = {synergy, ...} @ inputs:
    synergy.lib.mkFlake {
      inherit inputs;
      src = ./nix;
      eval.synergy.restrictDependenciesUnits.harmony = ["harmony"];
    };
}
