#!/dev/null

set -e -E -u -o pipefail -o noclobber -o noglob +o braceexpand || exit 1
trap 'printf "[ee] failed: %s\n" "${BASH_COMMAND}" >&2' ERR || exit 1
export -n BASH_ENV

_workbench="$( readlink -e -- . )"
_sources="${_workbench}/sources"
_scripts="${_workbench}/scripts"
_tools="${pallur_tools:-${_workbench}/.tools}"
_temporary="${pallur_temporary:-${pallur_TMPDIR:-${TMPDIR:-/tmp}}}"
_outputs="${_temporary}/$( basename -- "${_workbench}" )--outputs--$( readlink -e -- "${_workbench}" | tr -d '\n' | md5sum -t | tr -d ' \n-' )"
_applications_elf="${_outputs}/applications-elf"

_PATH="${pallur_PATH:-${_tools}/bin:${PATH}}"
_HOME="${pallur_HOME:-${HOME}}"
_TMPDIR="${pallur_TMPDIR:-${TMPDIR:-${_temporary}}}"

if test -n "${pallur_pkg_go:-}" ; then
	_GOROOT="${pallur_pkg_go}"
elif test -e "${_tools}/pkg/go" ; then
	_GOROOT="${_tools}/pkg/go"
else
	_GOROOT="${GOROOT}"
fi
_GOOS="${GOOS:-linux}"
_GOARCH="${GOARCH:-386}"
_GOPATH="${_outputs}/go"

if test -n "${_GOROOT:-}" ; then
	_go_bin="${_GOROOT}/bin/go"
else
	_go_bin="$( PATH="${_PATH}" type -P -- "${_go_bin}" || printf -- "${_go_bin}" )"
fi
if test -z "${_go_bin}" ; then
	echo "[ee] missing \`${_go_bin}\` (Go tool) executable in path: \`${_PATH}\`; ignoring!" >&2
	_go_bin=false
fi

_generic_env=(
		PATH="${_PATH}"
		HOME="${_HOME}"
		TMPDIR="${_TMPDIR}"
)

_go_sources="${_sources}"
_go_env=(
		"${_generic_env[@]}"
		GOOS="${_GOOS}"
		GOARCH="${_GOARCH}"
		GOROOT="${_GOROOT}"
		GOPATH="${_GOPATH}"
)

_package_name="$( basename -- "$( readlink -e -- . )" )"
_package_version="${pallur_distribution_version:-0.7.0_dev}"
_artifacts_cache="${pallur_artifacts:-}"
