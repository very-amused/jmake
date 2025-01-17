# jmake
A FreeBSD 14.0+ jail script generator. (WIP)

## Philosophy
`jmake` is *not* a jail manager. Instead of trying to abstract the task of jail management away from the vanilla FreeBSD utilities, `jmake` instead generates commands and configuration for use with these vanilla utilities. This philosophy ensures that `jmake` doesn't tie its users or their systems to any non-vanilla tooling. You can choose to stop using jmake and write these commands/configurations with another method (either manually or via different tool\[s\]) at *any time*.

### Editing /etc/rc.conf
`jmake` may generate rc config blocks in `*.rc.conf` files (i.e `bridge.rc.conf` for bridge networking configuration). These files are intended to have their contents appended to `/etc/rc.conf`, which can be achieved by using `sudo tee -a` to append with root privileges.
```sh
# Example: append the contents of bridge.rc.conf to /etc/rc.conf
cat bridge.rc.conf | sudo tee -a /etc/rc.conf
# Example: append the contents of all generated *.rc.conf files to /etc/rc.conf
cat *.rc.conf | sudo tee -a /etc/rc.conf
```
