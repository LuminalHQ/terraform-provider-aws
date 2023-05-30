# List Resources

The AWS provider is very old, and unfortunately some libraries it uses were
changed in backwards incompatible ways.  To avoid linking with this, instead of
providing an interface in code we provide a command which generates a list of
types.
