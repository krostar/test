{
  deps,
  lib,
  ...
}: let
  invalidGoCode = [
    "internal/code/testdata" # dir contains unparsable go code for test purprose
  ];
in {
  dev.formatters.treefmt.settings.formatter = {
    gci.excludes = builtins.map (f: "${f}/*") invalidGoCode;
    goimports.excludes = builtins.map (f: "${f}/*") invalidGoCode;
    gofumpt.excludes = builtins.map (f: "${f}/*") invalidGoCode;
  };
  ci.linters.golangci-lint.linters = {
    exclusions = {
      rules = [
        {
          # to test the package we hardcode things like !true to get representation of such expression
          path = "internal/message/from_bool_test\.go";
          linters = ["gocritic" "govet" "revive" "staticcheck" "stylecheck" "depguard"];
          text = "boolExprSimplify|dupSubExpr|nilness|bool-literal-in-expr|constant-logical-expr|dupArg|ST1019|SA4000|QF1001|SA4013|import 'reflect' is not allowed";
        }
        {
          # to test the package we hardcode things like !true to get representation of such expression
          path = "internal/compare/testify/compare_test\.go";
          linters = ["gocritic" "revive" "staticcheck" "testifylint"];
          text = "dupSubExpr|bool-literal-in-expr|constant-logical-expr|SA4000|S1040|useless-assert|require-error";
        }
        {
          path = "internal/code/get_caller_call_expr_test\.go|internal/code/parse_package_ast_test\.go";
          linters = ["staticcheck"];
          text = "SA5011";
        }
        {
          path = "internal/testing\.go";
          linters = ["iface"];
          text = "interface TestingT is declared but not used within the package";
        }
        {
          path = "double/spy_record\.go";
          linters = ["depguard"];
          text = "import 'reflect' is not allowed";
        }
      ];
      paths = invalidGoCode;
    };
    settings = {
      depguard = {rules.test = lib.mkForce null;}; # test rules are meant to protect against using other dependencies
      revive.rules =
        builtins.map (
          v:
            if v.name == "function-result-limit"
            then v // {arguments = [4];}
            else v
        )
        deps.synergy.result.data.harmony.ci.linters.golangci-lint.linters.settings.revive.rules;
    };
  };
}
