#!/dev/null

if ! test "${#}" -eq 0 ; then
	echo "[ee] invalid arguments; aborting!" >&2
	exit 1
fi

if test ! -e "${_outputs}" ; then
	mkdir -- "${_outputs}"
fi

exit 0
