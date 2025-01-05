package main

// Additional check functions that can be added to any context. See [checks.go] for details.
type ContextChecks struct{}

// Check that a script is being run as root
func (_ *ContextChecks) NeedsRoot() string {
	return "if [ `id -u` != 0 ]; then echo 'This script must be run as root.'; exit 1; fi"
}

// Check that the previously run command succeeded
func (_ *ContextChecks) CheckResult() string {
	return "[ \"$?\" == 0 ] || exit 1"
}
