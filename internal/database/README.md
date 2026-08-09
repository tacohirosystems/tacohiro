# database

Integrates with SQLite's C library (`libsqlite`) through `cgo` and sets up pools
for read and write connections.

It takes a slightly different approach from the existing SQLite Go libraries
(through `cgo`) by:

1. Avoiding manually getting and releasing the connection from the pool through
the `db.Read` and `db.Write` helper functions.
2. Abstracting away as little as possible while keeping access to the lower level
`libsqlite` functions.
3. Using precompiled `.so` files for faster compile times but needs to be statically
linked through `nix` to build the server binary.

Why? Because idk. I want a clean `go.sum`/`go.mod` file lol.
