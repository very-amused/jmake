# jmake
A FreeBSD 14.0+ jail script generator. (WIP)

## Philosophy
`jmake` is *not* a jail manager. Instead of trying to abstract the task of jail management away from the vanilla FreeBSD utilities, `jmake` instead generates commands and configuration for use with these vanilla utilities. This philosophy ensures that `jmake` doesn't tie its users or their systems to any non-vanilla tooling. You can choose to stop using jmake and write these commands/configurations with another method (either manually or via different tool\[s\]) at *any time*.

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
