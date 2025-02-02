package main

import (
	"fmt"
	"os"
	"path"
)

// JailConfiguration - Configuration for configuring jail systems/templates (jail-{name}-configure.sh) or a FreeBSD image (img-bootstrap.sh)
type JailConfiguration struct {
	Packages   []string // FreeBSD packages to be installed on the target using `pkg`
	User       *UserConfigurations
	PasswdFile string // Path to bootstrap user password file (default: .secret/jmake.passwd)

	userKeys []string // Key order for User field (TODO)
}

type UserConfigurations map[string]*UserConfiguration

// Information used to configure a user account. See PasswdFile for setting this user's password
type UserConfiguration struct {
	Username string `toml:"-"` // Username [parsed from toml key]
	FullName string // User full name (default: none)

	Uid int // Explicitly provided UID (default: let FreeBSD choose our UID)

	LoginGroup string   // Login group to create for the user (default: let FreeBSD create a login group with the same name as the user)
	Groups     []string // Supplementary groups to invite the user to (default: [])

	LoginClass string // Login class (default: "default" [set by FreeBSD])

	Shell string // Shell command for user (default: "sh" [set by FreeBSD]) NOTE: this different from Linux. FBSD wants an executable name, NOT a full path

	CreateHome *bool       // Whether to make a home directory for the user (default: true)
	HomeDir    string      // Location for user's home directory (default: "/home/{{.Name}}")
	HomePerms  os.FileMode // User homedir permissions (default: 0o700)

	GenPasswd bool // Whether to let FreeBSD generate a random password instead of prompting for a user-supplied password. (default: false)
	// NOTE: If GenPasswd is set in an *image* bootstrap config, each *jail* deployed from that image will receive its own unique password generated for that user.
	// Bootstrap user accounts are *locked* on images to ensure password security across different jails.

	Files []string // List of files to copy from the host to the user's homedir in the jail.
}

// Set default values where applicable
func (u *UserConfiguration) setDefaults() {
	if u.CreateHome == nil {
		u.CreateHome = new(bool)
		*u.CreateHome = true
	}
	if u.HomeDir == "" {
		u.HomeDir = path.Join("/home", u.Username)
	}
	if u.HomePerms == 0 {
		u.HomePerms = 0o700
	}
}

// Get HomePerms as an octal for use in pw useradd/usermod commands
func (u *UserConfiguration) HomePermString() string {
	return fmt.Sprintf("%o", u.HomePerms)
}
