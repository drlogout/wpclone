#!/usr/bin/env bash

release_url() {
  echo "https://file.noltech.net/wpclone"
}

download_release_from_file_server() {
  local version="$1"
  local os_info="$2"
  local tmpdir="$3"
  local architecture
  architecture="$(uname -m)"

  local filename="wpclone__${os_info}-${architecture}"
  local download_file="$tmpdir/$filename"
  local archive_url="$(release_url)/$filename"
  info 'Downloading' "$archive_url"

  curl -k --progress-bar --show-error --location --fail "$archive_url" --output "$download_file" --write-out "$download_file"
}


info() {
  local action="$1"
  local details="$2"
  command printf '\033[1;32m%12s\033[0m %s\n' "$action" "$details" 1>&2
}

error() {
  command printf '\033[1;31mError\033[0m: %s\n\n' "$1" 1>&2
}

warning() {
  command printf '\033[1;33mWarning\033[0m: %s\n\n' "$1" 1>&2
}

request() {
  command printf '\033[1m%s\033[0m\n' "$1" 1>&2
}

eprintf() {
  command printf '%s\n' "$1" 1>&2
}

bold() {
  command printf '\033[1m%s\033[0m' "$1"
}

# returns the os name to be used in the packaged release
parse_os_info() {
  local uname_str="$1"
  local arch="$(uname -m)"

  case "$uname_str" in
    Linux)
      if [ "$arch" == "x86_64" ]; then
        echo "linux"
      elif [ "$arch" == "aarch64" ]; then
        echo "linux-arm"
      else
        error "Releases for architectures other than x64 and arm are not currently supported."
        return 1
      fi
      ;;
    Darwin)
      echo "macos"
      ;;
    *)
      return 1
  esac
  return 0
}

parse_os_pretty() {
  local uname_str="$1"

  case "$uname_str" in
    Linux)
      echo "Linux"
      ;;
    Darwin)
      echo "macOS"
      ;;
    *)
      echo "$uname_str"
  esac
}

# return true(0) if the element is contained in the input arguments
# called like:
#  if element_in "foo" "${array[@]}"; then ...
element_in() {
  local match="$1";
  shift

  local element;
  # loop over the input arguments and return when a match is found
  for element in "$@"; do
    [ "$element" == "$match" ] && return 0
  done
  return 1
}

create_tree() {
  local install_dir="$1"

  info 'Creating' "directory layout"

  # .wpclone/
  #     bin/

  mkdir -p "$install_dir" && mkdir -p "$install_dir"/bin
  if [ "$?" != 0 ]
  then
    error "Could not create directory layout. Please make sure the target directory is writeable: $install_dir"
    exit 1
  fi
}

install_version() {
  local install_dir="$1"
  local should_run_setup="$2"


  install_release "$install_dir"

  if [ "$?" == 0 ]
  then
      if [ "$should_run_setup" == "true" ]; then
        info 'Finished' "installation. Updating user profile settings."
        "$install_dir"/bin/wpclone setup
      fi
  fi
}

install_release() {
  local install_dir="$1"

  download_archive="$(download_release; exit "$?")"
  exit_status="$?"
  if [ "$exit_status" != 0 ]
  then
    error "Could not download wpclone."
    return "$exit_status"
  fi

  install_from_file "$download_archive" "$install_dir"
}

download_release() {
  local uname_str="$(uname -s)"
  local os_info
  os_info="$(parse_os_info "$uname_str")"
  if [ "$?" != 0 ]; then
    error "The current operating system ($uname_str) does not appear to be supported by wpclone."
    return 1
  fi
  local pretty_os_name="$(parse_os_pretty "$uname_str")"

  info 'Fetching' "archive for $pretty_os_name, version $version"
  # store the downloaded archive in a temporary directory
  local download_dir="$(mktemp -d)"

  download_release_from_file_server "$version" "$os_info" "$download_dir"
}

install_from_file() {
  local archive="$1"
  local install_dir="$2"

  create_tree "$install_dir"

  info 'Extracting' "wpclone binaries and launchers"

  # extract the files to the specified directory
  cp "$archive" "$install_dir"/bin/wpclone
  chmod a+x "$install_dir"/bin/wpclone
}

# return if sourced (for testing the functions above)
return 0 2>/dev/null

# default to installing the latest available version
version_to_install="latest"

# default to running setup after installing
should_run_setup="true"

# install to wpclone_HOME, defaulting to ~/.wpclone
install_dir="${wpclone_HOME:-"$HOME/.wpclone"}"


install_version "$install_dir" "$should_run_setup"
