# jmake
A FreeBSD 14.0+ jail script generator. (WIP)

## Philosophy
`jmake` is *not* a jail manager. Instead of trying to abstract the task of jail management away from the vanilla FreeBSD utilities, `jmake` instead generates commands and configuration for use with these vanilla utilities. This philosophy ensures that `jmake` doesn't tie its users or their systems to any non-vanilla tooling. You can choose to stop using jmake and write these commands/configurations with another method (either manually or via different tool\[s\]) at *any time*.

## Example Config
```toml
# It's recommended to place the [host] section at the top in order to
# make any jmake.toml quickly identifiable and asssociable with its intended deployment host(s)

[host] # Configuration describing the jail deployment host
domain = "ncrypt.ing" # Network domain the host resides on (jails will be configured with this same domain)

[zfs] # Config for ZFS backed jail storage (ZFS is required for jmake)
dataset = "zroot/var/jail" # Name of the ZFS dataset used for jail storage
mountpoint = "/var/jail" # Mountpoint for the ZFS dataset used for jail storage

# Config for the FreeBSD image which will be used as a template for jail deployment
[img]
# FreeBSD release version to deploy jails from
release = "14.2-RELEASE"
# Name of the ZFS template snapshot to deploy jails from
snapshot = "base-custom0"
# FreeBSD architecture string
# (optional, default: "amd64/amd64")
arch = "amd64/amd64"
# FreeBSD download mirror URL used for fetching release media
# (optional, default: "https://download.freebsd.org/ftp/releases")
mirror = "https://download.freebsd.org/ftp/releases"

# Config for bootstrapping the jail deployment template snapshot. See bootstrap.go for all available options
# (optional)
#
# Image bootstrapping is useful to describe common options (packages, users, etc.) that are expected for *all* jails
[img.bootstrap]
packages = ["bash", "bash-completion", "sudo", "neovim"]
user.alice = {groups = ["wheel"], shell = "bash", files = [".bashrc"]} # Alice is an admin who prefers bash as a shell
user.bob = {shell = "tcsh", files = [".cshrc"]} # Bob is an unprivileged user who prefers tcsh as a shell

# Jail configuration for a redundant web server setup
[jail.dmz] # This jail gets a hostname of dmz.ncrypt.ing
description = "Reverse proxy to application servers"
network = "192.168.120.3/24"
[jail.dmz.bootstrap]
packages = ["haproxy"]

[jail.www1]
description = "Web server 1"
network = "192.168.120.4/24"
[jail.www1.bootstrap]
packages = ["node", "npm", "nginx"]
user.www = {home_dir = "/var/www", home_dir_perms = 0o755}

[jail.www2]
description = "Web server 2"
network = "192.168.120.5/24"
[jail.www2.bootstrap]
packages = ["node", "npm", "nginx"]
user.www = {home_dir = "/var/www", home_dir_perms = 0o755}

[jail.api1]
description = "API server 1"
network = "192.168.120.6/24"
[jail.api1.bootstrap]
packages = ["go"]
user.api = {home_dir = "/var/www/api", home_dir_perms = 0o700}

[jail.api2]
description = "API server 2"
network = "192.168.120.7/24"
[jail.api2.bootstrap]
packages = ["go"]
user.api = {home_dir = "/var/www/api", home_dir_perms = 0o700}
```

## Usage

### Generated Scripts
`jmake` generates its commands in the form of `*.sh` files. The name of these files follows the format `{section}-{action}.sh`, where {section} is a top level config section (i.e `zfs`, `img`, `bridge`) and {action} describes what the script *does* in relation to its section. For example, `zfs-init.sh` initializes the ZFS pool configured in the `[zfs]` section, which can then be used to store FreeBSD images and jails. Conversely, `zfs-destroy.sh` destroys this ZFS pool.

These generated scripts are intentionally made not executable, because they're meant to be run using `sh -xv {script}` rather than directly executed. These scripts are intended to simply store jail management commands, they don't produce robust or useful output on their own to see what's going on (this may change in the future). `sh -vx` will cause the content of these scripts (jail management commands and comments describing what they're doing) to be echoed as they run.

### Editing /etc/rc.conf
`jmake` may generate rc config blocks in `*.rc.conf` files (i.e `bridge.rc.conf` for bridge networking configuration). These files are intended to have their contents appended to `/etc/rc.conf`, which can be achieved by using `sudo tee -a` to append with root privileges.
```sh
# Example: append the contents of bridge.rc.conf to /etc/rc.conf
cat bridge.rc.conf | sudo tee -a /etc/rc.conf
# Example: append the contents of all generated *.rc.conf files to /etc/rc.conf
cat *.rc.conf | sudo tee -a /etc/rc.conf
```
