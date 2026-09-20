#!/bin/sh
set -e

if [ "$1" = "./socksserver" ] || [ "$1" = "socksserver" ] || [ "$1" = "/socksserver" ]; then
    shift
    exec /socksserver "$@"
fi

if [ "${1#-}" != "$1" ]; then
    exec /socksserver "$@"
fi

exec "$@"
