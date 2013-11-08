#!/dev/null

if ! test "${#}" -eq 0 ; then
	echo "[ee] invalid arguments; aborting!" >&2
	exit 1
fi

cd "${_go_sources}"

find -L . -type f \( -name '*.go' -o -name '*.c' -o -name '*.h' \) -print \
| while read _file ; do
	if test "${_file}" -nt "${_outputs}/go/src/${_file}" ; then
		_file_dirname="$( dirname -- "${_file}" )"
		if ! test -e "${_outputs}/go/src/${_file_dirname}" ; then
			mkdir -p -- "${_outputs}/go/src/${_file_dirname}"
		fi
		cp -T -- "${_file}" "${_outputs}/go/src/${_file}"
	fi
done

cd "${_outputs}/go/src"

find -L . -type f \( -name '*.go' -o -name '*.c' -o -name '*.h' \) -print \
| while read _file ; do
	if ! test -e "$( dirname -- "${_go_sources}/${_file}" )" ; then
		continue
	fi
	if ! test -e "${_go_sources}/${_file}" ; then
		rm -- "${_file}"
	fi
done

cd "${_outputs}/go"

if test ! -e "${_applications_elf}" ; then
	mkdir -p -- "${_applications_elf}"
fi

while read _application _main ; do
	echo "[ii] building \`${_application}\`..." >&2
	env "${_go_env[@]}" "${_go_bin}" build -o "${_applications_elf}/${_application}.elf" "./src/${_main}"
done <"${_go_sources}/applications.txt"

exit 0
