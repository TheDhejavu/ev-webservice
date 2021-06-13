
#!/bin/bash -e

exec > >(tee -a /var/log/app/entry.log|logger -t server -s 2>/dev/console) 2>&1

APP_ENV=${APP_ENV:-local}

echo "[`date`] Running entrypoint script in the '${APP_ENV}' environment..."

CONFIG_PATH=./config
CONFIG_FILE=${CONFIG_PATH}/${APP_ENV}.yml

if [[ -z ${DB_SOURCE} ]]; then
  export DB_SOURCE=`sed -n 's/^dbsource:[[:space:]]*"\(.*\)"/\1/p' ${CONFIG_FILE}`
fi

echo "[`date`] Starting server..."
./server -config ${CONFIG_FILE} >> /var/log/app/server.log 2>&1