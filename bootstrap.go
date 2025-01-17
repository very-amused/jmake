package main

import "os"

// BootstrapConfig - Configuration for bootstrapping jail(s) (jail-{name}-bootstrap.sh) or a FreeBSD image (img-bootstrap.sh)
type BootstrapConfig struct {
	Packages   []string // FreeBSD packages to be installed on the target using `pkg`
	User       *BootstrapUsers
	PasswdFile string // Path to bootstrap user password file (default: .secret/jmake.passwd)

	userKeys []string // Key order for User field
}

type BootstrapUsers map[string]*BootstrapUser

// Information used to boostrap a user account. See PasswdFile for setting this user's password
type BootstrapUser struct {
	Username string `toml:"-"` // Username [parsed from toml key]
	FullName string // User full name (default: none)

	Uid int // Explicitly provided UID (default: let FreeBSD choose our UID)

	LoginGroup string   // Login group to create for the user (default: let FreeBSD create a login group with the same name as the user)
	Groups     []string // Supplementary groups to invite the user to (default: [])

	LoginClass string // Login class (default: "default" [set by FreeBSD])

	Shell string // Shell command for user (default: "sh" [set by FreeBSD]) NOTE: this different from Linux. FBSD wants an executable name, NOT a full path

	CreateHome bool        // Whether to make a home directory for the user (default: true)
	HomeDir    string      // Location for user's home directory (default: "/home/{{.Name}}" [set by FreeBSD])
	HomePerms  os.FileMode // User homedir permissions (default: 0o700)

	GenPasswd bool // Whether to let FreeBSD generate a random password instead of prompting for a user-supplied password.
	// NOTE: If GenPasswd is set in an *image* bootstrap config, each *jail* deployed from that image will receive its own unique password generated for that user.
	// Bootstrap user accounts are *locked* on images to ensure password security across different jails.
}
