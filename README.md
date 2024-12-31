# jmake
A FreeBSD 14.0+ jail script generator. (WIP)

## Philosophy
`jmake` is *not* a jail manager. Instead of trying to abstract the task of jail management away from the vanilla FreeBSD utilities, `jmake` instead generates commands and configuration for use with these vanilla utilities. This philosophy ensures that `jmake` doesn't tie its users or their systems to any non-vanilla tooling. You can choose to stop using jmake and write these commands/configurations with another method (either manually or via different tool\[s\]) at *any time*.
