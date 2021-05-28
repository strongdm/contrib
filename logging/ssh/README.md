# Extract SSH Logs

This folder contains a shell/python script that can be used for extracting SSH full logs

## Requirements
* Python3
* SDM logs

## Configuration
You can adjust the following variables at the top of the script:
* LOG_DIR. Folder where the relay logs are located
* LOG_PATTERN. Logs name pattern

Ideally, configure the script as a CRONJOB. For example (daily config):
```
0 0 * * *
```

## Sample
```
$ ./extract_ssh_logs.sh 2>/dev/null
session_id,start_time,end_time,user,cmd(new_line=|#|)
s1r6n6RRiECGwflAJ2EPkv3YkBGC,2021-04-13 09:13:51.789592,2021-04-13 09:13:52.200592,Rodolfo Campos,Welcome to OpenSSH Server|#||#|openssh-server:~$
s1r6n6RRiECGwflAJ2EPkv3YkBGC,2021-04-13 09:13:52.200592,2021-04-13 09:13:55.071592,Rodolfo Campos,ls|#|logs  ssh_host_keys  sshd.pid|#|openssh-server:~$
s1r6n6RRiECGwflAJ2EPkv3YkBGC,2021-04-13 09:13:55.071592,2021-04-13 09:13:56.126592,Rodolfo Campos,pwd|#|/config|#|openssh-server:~$
s1r6n6RRiECGwflAJ2EPkv3YkBGC,2021-04-13 09:13:56.126592,2021-04-13 09:13:59.501592,Rodolfo Campos,echo "hello world"|#|hello world|#|openssh-server:~$
s1r6n6RRiECGwflAJ2EPkv3YkBGC,2021-04-13 09:13:59.501592,2021-04-13 09:14:02.427592,Rodolfo Campos,ls|#|logs  ssh_host_keys  sshd.pid|#|openssh-server:~$
s1r6n6RRiECGwflAJ2EPkv3YkBGC,2021-04-13 09:14:02.427592,2021-04-13 09:14:02.831592,Rodolfo Campos,ls|#|logs  ssh_host_keys  sshd.pid|#|openssh-server:~$
s1r6n6RRiECGwflAJ2EPkv3YkBGC,2021-04-13 09:14:02.831592,2021-04-13 09:14:04.111592,Rodolfo Campos,exit|#|logout|#|
s1rWwpo8PrK5nJU9GnI5y3hDGYF9,2021-04-22 15:28:58.497968,2021-04-22 15:28:58.871968,Rodolfo Campos,Welcome to OpenSSH Server|#||#|
s1rWwpo8PrK5nJU9GnI5y3hDGYF9,2021-04-22 15:28:58.871968,2021-04-22 15:29:02.232968,Rodolfo Campos,openssh-server:~$ exit|#|logout|#|
```

Considerations:
* Commands with several lines are delimited by `|#|` 
* SSH session files will be stored, and deleted in the next iteration, in the folder where the script gets executed
