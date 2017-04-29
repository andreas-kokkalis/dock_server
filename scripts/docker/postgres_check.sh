#!/usr/bin/env bash

# func postgres_is_up - checks if postgres is up and accepting connecctions
    # Retry until the reponse is "localhost:5432 - accepting connections"
postgres_is_up(){
    set +x
    local is_up='false'
    for i in {1..40}
    do
        sleep 2
        local response=$(docker exec -i dockserver_postgresdb_1 pg_isready -h localhost -p 5432)
        if [[ $? -eq 0 && $response == 'localhost:5432 - accepting connections' ]]; then
            is_up='true'
            break
        else
            continue
        fi
    done
    set -x
    echo $is_up
}
