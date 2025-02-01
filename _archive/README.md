# Obsolete Experiments

These are some initial experiments that are abandoned but kept for posterity.

## dotnet (1)

This implemented the solution using dotnet and YARP.
The implementation is very simple and straightforward, brings in a lot of powerful features, and leverages the use of a well-maintained library that is specfically for this purpose.
But the baseline memory consumption is about 25MB, which isn't terrible but is a lot higher than a basic Go implementation.

## dotnet2

For fun I decided to try implementing the reverse proxy with more basic dotnet primitives, to see how low the memory footprint could go.
This implementation brings down the base memory footprint to 10MB. This is good, but still higher than the Go implementation.

## rust

In theory Rust is the most attractive option.
The lack of a garbage collector should result in a much smaller binary as well as better performance.
I tried to get started with Rust, but then thought to myself:
"What the heck am I doing? I've never programmed in Rust before? Why torture myself like this?"
