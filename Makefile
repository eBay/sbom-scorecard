phony:
	@echo Use specific targets to download individual needed files.

examples/julia.spdx.json:
	curl -Lo examples/julia.spdx.json https://github.com/JuliaLang/julia/raw/master/julia.spdx.json
