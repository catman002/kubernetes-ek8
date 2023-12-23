#!/bin/sh

# graceful shutdown
trap "graceful_shutdown; exit 0;" SIGTERM SIGINT
graceful_shutdown() {
    echo "[INFO] SIGTERM caught, terminating <keepalived> process..."
    stop_process
    echo "[INFO] keepalived graceful shutdown successfully."
}
stop_process() {
    stop_keepalived
}

# start keepalived
KEEPALIVED_LAUNCH="keepalived --dont-fork --dump-conf --log-console --log-detail --log-facility 7 --vrrp -f /etc/keepalived/keepalived.conf"
start_keepalived() {
    echo "[INFO] Keepalived is starting."
    eval ${KEEPALIVED_LAUNCH} &
}

# stop keepalived
stop_keepalived() {
    count=1
    k_pid=$(pidof keepalived)
    while true; do
    	kill -TERM $k_pid > /dev/null 2>&1
        sleep 3
        
        k_pid=$(pidof keepalived)
        if [ ! -n "$k_id" ]; then
           break
        fi
       
        if [ ${count} -gt 5 ]; then
            echo "[ERROR] Keepalived stop failed."
            return
        fi
        count=$[count+1]
     done
     echo "[INFO] Keepalived terminated." 
#    kill -TERM $(cat /var/run/vrrp.pid)
#    kill -TERM $(cat /var/run/keepalived.pid)
}

start_keepalived

sleep 10

# while-loop to ensure  keepalived health.
while true; do
    k_pid=$(pidof keepalived)

    if [ ! -n "$k_pid" ]; then
        # when one crash, graceful shutdown
        echo "[ERROR] Keepalived crashed, shutdown all process, exit 1"
        stop_process
        break
    fi
    sleep 5
done

exit 1


