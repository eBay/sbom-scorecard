tests: examples/julia.spdx.json
	go run . examples/julia.spdx.json

examples/julia.spdx.json:
	curl -Lo examples/julia.spdx.json https://github.com/JuliaLang/julia/raw/master/julia.spdx.json
