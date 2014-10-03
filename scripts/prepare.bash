#!/dev/null
## chunk::59c4750c295e27c2e75c134980f1e813::begin ##

if ! test "${#}" -eq 0 ; then
	echo "[ee] invalid arguments; aborting!" >&2
	exit 1
fi

## chunk::3c8b019c663118b00172b22aeae97568::begin ##
if test ! -e "${_temporary}" ; then
	mkdir -- "${_temporary}"
fi
if test ! -e "${_outputs}" ; then
	mkdir -- "${_outputs}"
fi
## chunk::3c8b019c663118b00172b22aeae97568::end ##

exit 0
## chunk::59c4750c295e27c2e75c134980f1e813::end ##
