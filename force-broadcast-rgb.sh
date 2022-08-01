#!/usr/bin/env bash
set -e
declare -A PROP_SETS

DRMCLI_ARGS="-D /dev/dri/card0"
PROPTEST_ARGS="-D /dev/dri/card0 -M i915"
PROP_SETS[239]=1

main() {
	resources="$(drmcli resources)"
	for connectorID in $(j "$resources" '.Connectors[]'); {
		info="$(drmcli connector info $connectorID)"
		[[ $(j "$info" '.Status') != connected ]] && continue
		[[ $(j "$info" '.Type') == eDP ]] && continue
		
		for propertyID in "${!PROP_SETS[@]}"; {
			proptest $PROPTEST_ARGS $connectorID connector $propertyID ${PROP_SETS[$propertyID]}
			echo "connector $connectorID: $propertyID = ${PROP_SETS[$propertyID]}"
		}
	}
}

j() {
	jq -r "${@:2}" <<< "$1"
}

drmcli() {
	if [[ $(type drmcli) == file ]]; then
		command drmcli $DRMCLI_ARGS "$@"
		return $?
	fi

	if command -v go &> /dev/null; then
		go run github.com/diamondburned/force-broadcast-rgb/cmd/drmcli $DRMCLI_ARGS "$@"
		return $?
	fi

	echo "drmcli not installed"
	return 1
}

command -v proptest &> /dev/null || {
	echo "proptest not installed"
	exit 1
}

main "$@"
