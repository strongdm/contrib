import os
import datetime
import subprocess
import platform
import json
import csv


def get_accounts():
    '''
     Gather Accounts from the sdm CLI
      * Requires being logged into the sdm client
    '''
    sdm_command = ['sdm',
                   'audit',
                   'users',
                   'list',
                   '-j']
    result = subprocess.run(sdm_command, stdout=subprocess.PIPE, check=False)
    accounts = format_cli_results(result)
    return accounts


def get_activities(lookback_period):
    '''
     Gather Activities using the sdm cli
      * Requires being logged into the sdm client
    '''
    today = datetime.date.today()
    from_date = today - datetime.timedelta(days=lookback_period)
    from_date = from_date.strftime('%Y/%m/%d')
    sdm_command = ['sdm',
                   'audit',
                   'activities',
                   '-e',
                   '-j',
                   '--from',
                   from_date]
    result = subprocess.run(sdm_command, stdout=subprocess.PIPE, check=False)
    activities = format_cli_results(result)
    return activities


def format_cli_results(result):
    '''
    Convert CLI results into Python Dict
    '''
    result_string = result.stdout.decode('utf-8')
    # Json results include newlines
    # So splitting lines gets indivudal json objects
    individual_record_strings = result_string.splitlines()
    formatted_cli_results = []
    # Convert individual json string objects into python dicts
    for record in individual_record_strings:
        json_record = json.loads(record)
        formatted_cli_results.append(json_record)
    return formatted_cli_results


def prepare_stats_for_collection(accounts):
    '''
    Create stats Dict
    '''
    stats = {}
    for account in accounts:
        # Skip Suspended Accounts
        if account['strongRole'] == 'suspended':
            continue
        # Differentiate between Service and Human account types
        if account['firstName'] == 'Service Account':
            stats[account['id']] = {"name": account['lastName'], "type": "service_account"}
        else:
            stats[account['id']] = {"name": account['email'], "type": "human_account"}
    for account in stats.items():
        stats[account[0]].update({"last_login": None})
    return stats


def generate_stats(activities, accounts):
    '''Stats Dict for data profiling'''
    stats = prepare_stats_for_collection(accounts)
    for activity in activities:
        # Skip events:
        #   * That were not logins
        #   * That were from support users, etc
        #   * That were from deleted users
        if activity['activity'].startswith('user logged into') is False or activity['actorUserID'] not in stats:
            continue
        if stats[activity['actorUserID']]['last_login'] is None:
            stats[activity['actorUserID']]['last_login'] = activity['timestamp']
        elif stats[activity['actorUserID']]['last_login'] < activity['timestamp']:
            stats[activity['actorUserID']]['last_login'] = activity['timestamp']
    for account in stats.items():
        if account[1]['last_login'] is None: 
            stats[account[0]].update({"last_login": "No logins for the last 90 days."})
    return stats


def format_stats_for_csv(stats):
    '''
    Convert single stats Dict to List of Dicts
    to simplify csv creation
    '''
    formatted_stats =  []
    for key in stats:
        formatted_stats.append(stats[key])
    return formatted_stats


def create_csv(stats):
    '''Write stats to CSV'''
    formatted_stats = format_stats_for_csv(stats)
    field_names = ['name', 'type', 'last_login']
    with open('StrongDM_Logins_Past_90_Days.csv', 'w', encoding='utf-8') as csvfile:
        writer = csv.DictWriter(csvfile, fieldnames=field_names)
        writer.writeheader()
        writer.writerows(formatted_stats)

def open_csv():
  ''' Open CSV on completion '''
  filepath='StrongDM_Logins_Past_90_Days.csv'
  if platform.system() == 'Darwin':        # Mac
    subprocess.call(('open', filepath))
  elif platform.system() == 'Windows':    # Windows
      os.startfile(filepath)
  else:                                   # linux variants
      subprocess.call(('xdg-open', filepath))


def main():
    ''' Run it all!'''
    accounts = get_accounts()
    activities = get_activities(90)
    stats = generate_stats(activities, accounts)
    create_csv(stats)
    open_csv()
    exit()


if __name__ == "__main__":
    main()
