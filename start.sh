#!/bin/sh
base_path=/app
exec_path=${base_path}/tetris-server
conf_path=${base_path}/bin/conf/server.json
log_path=${base_path}/logs
bi_path=${base_path}/bi
stop_path=${base_path}
sh ${stop_path}/stop.sh
${exec_path} -conf ${conf_path} -log ${log_path} -bi ${bi_path} -wd ${base_path} 2>&1