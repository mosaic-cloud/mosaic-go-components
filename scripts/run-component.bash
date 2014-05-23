#!/dev/null

_identifier="${1:-00000000cc01bfbe028de269636921dadcf2999c}"
_fqdn="${mosaic_node_fqdn:-mosaic-0.loopback.vnet}"
_ip="${mosaic_node_ip:-127.0.155.0}"

if test -n "${mosaic_component_temporary:-}" ; then
	_tmp="${mosaic_component_temporary}"
elif test -n "${mosaic_temporary:-}" ; then
	_tmp="${mosaic_temporary}/components/${_identifier}"
else
	_tmp="/tmp/mosaic/components/${_identifier}"
fi
if test "${_identifier}" == 00000000cc01bfbe028de269636921dadcf2999c ; then
	_tmp="${_tmp}--$( date +%s )"
fi

_run_bin="${_applications_elf}/component-backend.elf"
_run_env=(
		mosaic_component_identifier="${_identifier}"
		mosaic_component_temporary="${_tmp}"
		mosaic_node_fqdn="${_fqdn}"
		mosaic_node_ip="${_ip}"
)

case "${_identifier}" in
	
	( 00000000190a256e5dcaa1825e8c17117d5415ad )
		if ! test "${#}" -ge 2 ; then
			echo "[ee] invalid arguments; aborting!" >&2
			exit 1
		fi
		_run_args=(
				component-me2-init
				"${@:2}"
		)
	;;
	
	( 00000000cc01bfbe028de269636921dadcf2999c )
		if ! test "${#}" -eq 0 ; then
			echo "[ee] invalid arguments; aborting!" >&2
			exit 1
		fi
		_run_args=(
				standalone
		)
	;;
	
	( * )
		if ! test "${#}" -ge 1 ; then
			echo "[ee] invalid arguments; aborting!" >&2
			exit 1
		fi
		_run_args=(
				component
				"${_identifier}"
				"${@:2}"
		)
	;;
esac

mkdir -p -- "${_tmp}"
cd -- "${_tmp}"

exec env "${_run_env[@]}" "${_run_bin}" "${_run_args[@]}"

exit 1
