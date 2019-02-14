# GitHub SSH AuthorizedKeysCommand Utility

Use the SSH public keys you've added to GitHub to log in to your machines!

## 0. Compile for your target OS

Currently, ghsshauth has been tested on Debian and FreeBSD, but should basically
support all UNIX-like systems. Feel free to send patches.

```
$ git clone git@github.com:benburwell/ghsshauth.git
$ cd ghsshauth
$ sudo make install
```

## 1. Configure host

Open `/etc/ssh/sshd_config`, find the `AuthorizedKeysCommand` line, and change
it to:

```
AuthorizedKeysCommand /usr/local/sbin/ghsshauth %h
AuthorizedKeysCommandUser root
```

(the `%h` represents the home directory of the user being authenticated).

In your home directory, create the file `.ssh/authorized_github_users` and add
your GitHub username (and any other username you want to have access) to the
file, one per line. You can begin lines with the `#` character to have them be
ignored.

**IMPORTANT:** You'll need to make sure that the `AuthorizedKeysCommandUser` has
read access to the entire path up to your `authorized_github_users` file. The
easy way to do this is to make set `AuthorizedKeysCommandUser root`. If you'd
rather use `AuthorizedKeysCommandUser nobody`, then you'll need to make sure
the `nobody` user has read access to `~/.ssh/authorized_github_users`. This
means your home directory needs to be `chmod 755` as does your `.ssh` directory.
**If you go this route, be sure that any secret keys in your `.ssh` directory
such as `id_rsa` are `chmod 600`, else secret keys they will no longer be!**
