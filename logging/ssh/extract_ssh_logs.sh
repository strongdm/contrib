#!/usr/bin/python3 
#  strongDM - Daily SSH session commands extract (from strongDM servers) 
#
#  Extracts all previous day ssh events from strongDM logs
#  Creates a stdout file for each ssh session-id
#
#  Cron daily, at midnight. crontab 0 0 * * *

import base64
import datetime
import os
import re
import subprocess
import sys

LOG_DIR="logs"
LOG_PATTERN="relay*.log"

def main():
    successful = generate_ssh_session_files()
    if not successful:
        print(f"Error splitting log files from {LOG_DIR}")
        exit(-1) 
    print_logs()

def generate_ssh_session_files():
    ssh_split_cmd = f"""
    rm *.ssh && \
    find {LOG_DIR} -name '{LOG_PATTERN}' | while read line; do 
      sdm ssh split $line 1>&2
    done
    ls *.ssh &> /dev/null
    """
    result = os.system(ssh_split_cmd)
    return result == 0

def print_logs():
    init_date = get_init_date()
    print("session_id,start_time,end_time,user,cmd(new_line=|#|)")
    for count, line in enumerate(run_command(f"sdm audit ssh --from {init_date}")):
        if count == 0:
            continue

        is_relay, session_id, start_time, user = extract_session_info(line)
        if is_relay != "true":
            continue

        if not os.path.isfile(f"{session_id}.ssh"):
            print("Session: $session_id is not present in the provided log", file=sys.stderr)
            continue

        print_session_logs(session_id, start_time, user)

def get_init_date():
    yesterday = datetime.datetime.now() - datetime.timedelta(days = 1)
    return yesterday.strftime("%Y-%m-%d")

def run_command(command):
    p = subprocess.Popen(command.split(), stdout = subprocess.PIPE, stderr = subprocess.STDOUT)
    return iter(p.stdout.readline, b'')

def extract_session_info(line):
    session_info = line.decode("utf-8").split(",")
    return session_info[8], session_info[6], session_info[0], session_info[4]

def print_session_logs(session_id, start_time, user):
    full_cmd_entry = ""
    total_elapsed_millis = 0
    start_time_regular = zulu_date_to_regular(start_time)
    for line in run_command(f"cat {session_id}.ssh"):
        elapsed_millis, cmd_entry = extract_cmd_entry_info(line)
        one_line_cmd_entry = cmd_entry.replace("\r", "").replace("\n", "|#|") 
        full_cmd_entry = f"{full_cmd_entry}{one_line_cmd_entry}" 
        total_elapsed_millis += elapsed_millis

        if not end_of_line(cmd_entry):
            continue
  
        end_time_regular = add_millis(start_time_regular, total_elapsed_millis)
        print(f"{session_id},{start_time_regular},{end_time_regular},{user},{full_cmd_entry}")

        full_cmd_entry = ""
        total_elapsed_millis = 0
        start_time_regular = end_time_regular

def zulu_date_to_regular(input_date):
    return re.sub(r"[0-9]{,3} \+[0-9]{,4} UTC$", "", input_date)

def extract_cmd_entry_info(line):
    cmd_entry_info = line.decode("utf-8").split(",")
    return int(cmd_entry_info[0]), base64.b64decode(cmd_entry_info[1]).decode("utf-8")

def end_of_line(cmd_entry):
    return True if "\n" in cmd_entry else False

def add_millis(input_date, input_millis):
    new_date = datetime.datetime.strptime(input_date, "%Y-%m-%d %H:%M:%S.%f") + datetime.timedelta(milliseconds = input_millis)
    return str(new_date)

main()
