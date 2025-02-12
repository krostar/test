{
  deps,
  lib,
  ...
}: let
  invalidGoCode = [
    "internal/code/internal/testdata" # dir contains unparsable go code for test purprose
  ];
in {
  dev.formatters.treefmt.settings.formatter = {
    gci.excludes = builtins.map (f: "${f}/*") invalidGoCode;
    goimports.excludes = builtins.map (f: "${f}/*") invalidGoCode;
    gofumpt.excludes = builtins.map (f: "${f}/*") invalidGoCode;
  };
  ci.linters.golangci-lint = {
    issues = {
      exclude-rules = [
        {
          # to test the package we hardcode things like !true to get representation of such expression
          path = "internal/message/from_bool_test.go";
          linters = ["gocritic" "govet" "revive" "staticcheck" "stylecheck"];
          text = "dupSubExpr|nilness|bool-literal-in-expr|constant-logical-expr|ST1019|SA4000";
        }
        {
          # to test the package we hardcode things like !true to get representation of such expression
          path = "internal/compare/testify/compare_test.go";
          linters = ["gocritic" "revive" "staticcheck" "testifylint"];
          text = "dupSubExpr|bool-literal-in-expr|constant-logical-expr|SA4000|useless-assert|require-error";
        }
      ];
      exclude-dirs = invalidGoCode;
    };
    linters-settings = {
      depguard = {rules.test = lib.mkForce null;}; # test rules are meant to protect using other dependencies
      revive.rules =
        builtins.map (
          v:
            if v.name == "function-result-limit"
            then v // {arguments = [4];}
            else v
        )
        deps.harmony.result.data.harmony.ci.linters.golangci-lint.linters-settings.revive.rules;
    };
  };
}
