{
  description = "MCP server for D&D 5e Italian SRD data";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "quintaedizione-mcp";
          version = "0.1.0";
          src = ./.;
          vendorHash = "sha256-1KmCcX7n7KiqVCzWExjhKmfUXfdu+4PIx9npqjsvtLI=";
          meta = with pkgs.lib; {
            description = "MCP server for D&D 5e Italian SRD data";
            homepage = "https://github.com/emiliopalmerini/quintaedizione-mcp";
            license = licenses.mit;
            mainProgram = "quintaedizione-mcp";
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [ go gopls ];
        };
      }
    );
}
