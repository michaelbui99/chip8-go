# chip8-go

<img src="./doc/emulator_running_ibm_logo_program.png">

Chip-8 emulator written in Go.<br>
The goal was to get more comfortable with Go and learning the basics of emulation.

Still in development, currently only supports enough instructions to run the IBM Logo program.

# Requirements

-   SDL2
-   Go (for installing the emulator)

# Installation

```console
$ go install github.com/michaelbui99/chip8-go@latest
```

# Load and run a Chip-8 ROM

```console
$ chip8-go load /PATH/TO/ROM
```

# See all available commands

```console
$ chip8-go -h
```

# References

## High level overview over components + specifications

-   https://tobiasvl.github.io/blog/write-a-chip-8-emulator/
