#!/usr/bin/env bash

DAEMON_PATH=${GOBIN}

HELLOWORLD_DAEMON=${DAEMON_PATH}/helloworld

usage() {
    cat <<EOF

usage:
    $(basename $0) start
    $(basename $0) restart
    $(basename $0) stop
    $(basename $0) status

EOF
    exit 1
}

COMMAND=""
case "$1" in
    start)
        COMMAND=start
        ;;

    stop)
        COMMAND=stop
        ;;

    restart)
        COMMAND=restart
        ;;

    status)
        COMMAND=status
        ;;

    *)
        usage
        ;;
esac
if [ ${COMMAND} != "" ];
then
    echo "Running command: ${COMMAND} for ${HELLOWORLD_DAEMON}"
    ${HELLOWORLD_DAEMON} ${COMMAND}
fi
