#!/bin/sh
sp_pid=`ps -ef | grep tetris-server | grep -v grep | awk '{print $2}'`
if [ -z "$sp_pid" ];
then
 echo "[ not find tetris-server pid ]"
else
 echo "find result: $sp_pid "
 kill -9 $sp_pid
fi