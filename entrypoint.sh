#! /bin/bash

if [ "$LITESTREAM_ACCESS_KEY_ID" != "" ]; then
    echo "dbs:" > /etc/litestream.yml
    echo "  - path: /app/arkstat.db" >> /etc/litestream.yml
    echo "    replicas:" >> /etc/litestream.yml
    echo "      - url: $LITESTREAM_REPLICA_URL" >> /etc/litestream.yml
    echo "        retention: 24h" >> /etc/litestream.yml
    echo "        sync-interval: 1m" >> /etc/litestream.yml
    litestream restore -if-replica-exists -v
    litestream replicate -exec /app/arkstat
else 
    /app/arkstat
fi

