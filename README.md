> **Warning** This implementation is still a heavy work in progress. It is not ready for use.

# go-polar
[![license](https://img.shields.io/github/license/Minestom/MinestomDataGenerator.svg)](LICENSE)

Go implementation of [Polar](https://github.com/hollow-cube/polar). Polar is a world format designed 
for simpler and smaller handling of small worlds, particularly for user generated content where size 
matters.

Polar generally should not be used for large worlds, since it stores worlds in a single file and does 
not allow random access of chunks (the entire world must be loaded to read chunks). As a general rule 
of thumb, Polar should only be used for worlds small enough that they are OK being completely kept in 
memory.

The Polar format is described in detail [here](https://github.com/hollow-cube/polar/blob/main/FORMAT.md)

## Features
todo
k
## Install
```shell
go get -u github.com/hollow-cube/go-polar
```

## Usage
todo

## Comparison to others
todo

## Contributing
Contributions via PRs and issues are always welcome.

## License
This project is licensed under the [MIT License](./LICENSE).
