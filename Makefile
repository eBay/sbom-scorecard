phony:
	@echo Use specific targets to download individual needed files.

examples/julia.spdx.json:
	curl -Lo examples/julia.spdx.json https://github.com/JuliaLang/julia/raw/master/julia.spdx.json

examples/dropwizard.cyclonedx.json:
	curl -Lo examples/dropwizard.cyclonedx.json https://github.com/CycloneDX/bom-examples/raw/master/SBOM/dropwizard-1.3.15/bom.json
