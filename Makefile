build:
	go build -o ./out/chip8-go ./main.go 

ibm: build
	./out/chip8-go load ./roms/IBM_Logo.ch8