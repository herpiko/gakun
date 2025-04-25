# Gakun

Gakun is a SSH key manager where you can manage and switch your SSH key easily. Instead of managing the keys in separated directory and utilizing symbolic link like [skm](https://github.com/TimothyYe/skm), Gakun manages the SSH keys directly within the `~/.ssh` directory, as it should be. So, whenever you want to go back to manual management, Gakun ensure that you will still having a great control over your keys.

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
