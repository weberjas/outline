# Outline

### Summary

Outline is a command line utility designed to give the user a fast, high-level overview of a GO project. It does this by scanning all the source files and building a graph of types, struct, methods, and functions. This graph is then printed to stdout using text. It is not meant as an exhaustive reference, merely a quick way to understand a project.

You may pass the flag -showOnlyExp=true to omit unexported functions from the printout
