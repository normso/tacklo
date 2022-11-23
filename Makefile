build_server:
	@echo "Compiling Server >>>>>"
	mkdir ./builds
	cd ./src/server/;\
		go build -o ../../builds/main main.go
	@echo "Compiling Done Now Running >>>>"

run_server:
	./builds/main

