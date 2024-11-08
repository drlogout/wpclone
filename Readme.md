# wplcone

<img src="you_better_use_wpclone.png" alt="you_better_use_wpclone" />

## YOU BETTER USE *wpclone*!

CLI tool to clone Wordpress. 

## Installation

```bash
curl https://file.noltech.net/wpclone/install.sh | bash
```

## Init

```bash
wpclone init
```

This generates a `wpclone.yml` which looks as follows :

```yaml
---

local:
  path: /Users/mofa/tmp/wordpress/local
  url: http://example.test
  db:
    name: wpclone
    user: wpclone
    password: wpclone
    # host: localhost # (default: 127.0.0.1)
    # port: 33060 # (default: 3306)
  # wp_cli: wp8.1 # (default: wp)
  # docker: 
    # all: true # (default: false)
    # ssl: true # (default: false)
    # db_only: true # (default: false)

remote:
  path: /home/example/example.com
  url: https://example.com
  ssh:
    host: example.com
    user: example
    # key: ~/.ssh/tux@croox.com # (default: ~/.ssh/id_ed25519 or ~/.ssh/id_rsa)
    # password: a-ssh-password # (if set, key will be ignored)
    # port: 2222 # (default: 22)
  # wp_cli: wp8.1 # (default: wp)
  # db: # (optional, only needed if wp-config.php does not exist on remote)
    # name: wpclone
    # user: wpclone
    # password: wpclone
    # host: localhost # (default: 127.0.0.1)
    # port: 33060 # (default: 3306)

```

## Project Directory

The project directory is a central place to store wpclone.yml files of different projects. By default the project directory is `~/wpclone`. It can be set with the `WPCLONE_PROJECT_DIR` environment variable. 

To init a wpclone.yml in the project directory specify the global `-p` flag, for example

```bash
wpclone -p croox.com init
wpclone -p croox.com pull
```

To show all projects in the project directory use `wpclone project ls`.

## Pull & Push

#### If Wordpress is installed on the remote host, the site can be pulled from remote:

```bash
wpclone -p croox.com pull
```

And after changes are made locally pushed to remote

```bash
wpclone -p croox.com push
```

## Install local Wordpress

If there is no Wordpress on the remote host it can be installed on the local machine and pushed afterwards

```bash
wpclone -p croox.com install
```

Because there is no `wp-config.php` on the remote server, we need to provide the remote DB credentials:

````yaml
remote:
  path: /home/example/example.com
  url: https://example.com
	...
  db: # (optional, only needed if wp-config.php does not exist on remote or is force pushed)
    name: wpclone
    user: wpclone
    password: wpclone
    # host: localhost # (default: 127.0.0.1)
    # port: 33060 # (default: 3306)
````

Then we can push to remote:

```bash
wpclone -p croox.com push
```
