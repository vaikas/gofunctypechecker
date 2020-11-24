# Go function detection library

This repo has a WIP library for being able to scan files and see if they match
a particular signature. The motivation for this was from learning to use buildpacks
and in particular having a detection logic that would be able to decide a cooler
user experience by inspecting the given go source code and give a highest abstraction
level supported by various buildpacks. Some examples are being able to only create
a handler method for [CloudEvents](https://github.com/cloudevents/sdk-go) library and
have the build know this is a supported signature and create the necessary scaffolding
to simply just call this function.

# Specifying supported function signatures

You can control which function signatures are supported either by populating the datastructures, or
by creating a config file (in json, later yaml & toml). For examples on how to use this, best place for
now is to look at the [test code](./pkg/detect/detect_test.go).