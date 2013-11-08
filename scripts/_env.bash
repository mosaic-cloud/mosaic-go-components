#!/dev/null

set -e -E -u -o pipefail -o noclobber -o noglob +o braceexpand || exit 1
trap 'printf "[ee] failed: %s\n" "${BASH_COMMAND}" >&2' ERR || exit 1
export -n BASH_ENV

_workbench="$( readlink -e -- . )"
_sources="${_workbench}/sources"
_scripts="${_workbench}/scripts"
_tools="${mosaic_distribution_tools:-${_workbench}/.tools}"
_outputs="${_workbench}/.outputs"
_temporary="${mosaic_distribution_temporary:-/tmp}"
_applications_elf="${_outputs}/applications-elf"

_GOOS=linux
_GOARCH=386
_GOROOT="${_tools}/pkg/go"
_GOPATH="${_outputs}/go"

_PATH="${_GOROOT}/bin:${_tools}/bin:${PATH}"

_go_bin="go"
_go_bin="$( PATH="${_PATH}" type -P -- "${_go_bin}" || printf -- "${_go_bin}" )"
if test -z "${_go_bin}" ; then
	echo "[ww] missing \`${_go_bin}\` (Go tool) executable in path: \`${_PATH}\`; ignoring!" >&2
fi

_generic_env=(
		PATH="${_PATH}"
		TMPDIR="${_temporary}"
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
_package_version="${mosaic_distribution_version:-0.7.0}"
