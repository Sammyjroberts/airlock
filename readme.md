## AIRLOCK

The goal of this repo is to create a simple and easy to use way to implement and
utilize plugins. For our purpose, a plugin is user generated code you would like
to utilize in your project. Existing solutions work using GRPC or RPC, but that
overhead can be a bit much for constant back and forth communication. Or super
performance critical applications (nanosecond scale)

Currently, The repo is int a prototype state, with a simple proof of concept.

The next goal is to use reflection to automatically generate the interface
implementation using the pipe server dynamically with go generics, but its a bit
hard to do that with the current state of go generics.

Following that will do the same for go recieving plugins, and python recieving
plugins.
