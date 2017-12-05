#!/bin/bash
case $1 in
'PROCESS_STATE_RUNNING')
	/bin/echo "$3 run" >> /root/event.log
	if [ $3 == 'micro_teacher_server' ];then
		supervisorctl restart micro_teacher_api
	fi
;;

'PROCESS_STATE_STOPPED')
# supervisorctl stop 
	/bin/echo "$3 stoped" >> /root/event.log
;;

'PROCESS_STATE_EXITED')
	/bin/echo "$3 exited" >> /root/event.log

;;

'PROCESS_STATE_FATAL')
	/bin/echo "$3 fatal" >> /root/event.log
;;

esac