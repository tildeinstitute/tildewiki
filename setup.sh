#!/bin/bash


install_tildewiki()
{
  [[ $(id -u) != 0 ]] && noroot_die "I can't be installed unless this is run as root"
  display_logo
  echo
  echo Building TildeWiki ...
  go clean
  go build -mod vendor
  echo Creating user/group ...
  setup_usergrp
  echo Copying files ...
  copy_files
  install_startup_script
  echo Cleaning up ...
  go clean
  cat > /dev/stdout << EOF
    
    TildeWiki data has been installed to
      /usr/local/tildewiki

    You may now feel free to add content!

    To start TildeWiki, run:
      /usr/local/bin/tildewiki

EOF
}

noroot_die()
{
  echo
  echo -e "${1:-"Error"}" >&2
  exit 1
}

setup_usergrp()
{
  adduser --home /usr/local/tildewiki --system --group tildewiki
}

copy_files()
{
  mkdir -p /usr/local/tildewiki/
  cp ./tildewiki /usr/local/tildewiki/
  cp -r pages /usr/local/tildewiki/
  cp -r assets /usr/local/tildewiki/
  cp tildewiki.yaml /usr/local/tildewiki/
  chown -R tildewiki:tildewiki /usr/local/tildewiki
}

install_startup_script()
{
  cat > /usr/local/bin/tildewiki << EOF
#!/usr/bin/env bash
error_exit() {
  echo -e "\${1:-"Unknown Error"}" >&2
  exit 1
}
[[ \$(id -u) != 0 ]] && error_exit "I can't daemonize unless I'm run as root :("
/usr/sbin/daemonize -c /usr/local/tildewiki -u tildewiki /usr/local/tildewiki/tildewiki
EOF
  chmod 755 /usr/local/bin/tildewiki
}

display_logo()
{
  cat > /dev/stdout << EOF
    __  _ __    __             _ __   _
   / /_(_) /___/ /__ _      __(_) /__(_)
  / __/ / / __  / _ \ | /| / / / //_/ /
 / /_/ / / /_/ /  __/ |/ |/ / / ,< / /
 \__/_/_/\__,_/\___/|__/|__/_/_/|_/_/

        :: TildeWiki v0.5.4 ::
     (c)2019 Ben Morrison (gbmor)
              GPL v3
   https://github.com/gbmor/tildewiki
     All Contributions Appreciated!
EOF
}

uninstall_tildewiki()
{
  [[ $(id -u) != 0 ]] && noroot_die "I can't be uninstalled unless this is run as root"
  display_logo
  echo
  echo Removing files ...
  rm -rf /usr/local/tildewiki
  rm -f /usr/local/bin/tildewiki
  echo Removing user/group ...
  userdel tildewiki
  echo TildeWiki successfully uninstalled!
  echo
}

display_help()
{
  display_logo
  cat >/dev/stdout<<EOF


     TildeWiki Installation Script

  install   | Installs TildeWiki data to /usr/local/tildewiki
              Places a start-up script at /usr/local/bin/tildewiki
            
  uninstall | Removes TildeWiki from the system
            
  help, -h  | Displays this message
EOF
}

case "$1" in
  install)
    install_tildewiki
    ;;
  uninstall)
    uninstall_tildewiki
    ;;
  help)
    display_help
    ;;
  -h)
    display_help
    ;;
  *)
    display_help
    ;;
esac

