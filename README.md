# Gakun

Gakun is an SSH key manager that allows you to manage and switch your SSH keys per host easily. Unlike [skm](https://github.com/TimothyYe/skm), which manages keys in separate directories and utilizes symbolic links, Gakun manages SSH keys directly within the `~/.ssh` directory, as it should be. This approach ensures that whenever you want to return to manual management, you will still have great control over your keys.

This software is still in heavy development. Please expect breaking changes and use it at your own risk.

## Usage
```
NAME:
   gakun - SSH key manager

USAGE:
   gakun [global options] [command [command options]]

COMMANDS:
   add      Add host and key to a profile. Example: 'gakun add work gitlab.com ~/.ssh/id_rsa_work'
   use      Use SSH key for certain host. Example: 'gakun use work -h gitlab.com'
   ls       List profiles
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```


# License

MIT
