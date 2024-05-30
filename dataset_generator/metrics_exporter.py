import requests
import json
import pandas as pd
import sys

from datetime import datetime, timezone, timedelta

container_resource = {}
args = {}

since = ((datetime.now(timezone.utc) - timedelta(hours=2))).strftime('%G-%m-%dT%X.000Z')
until = (datetime.now(timezone.utc)).strftime('%G-%m-%dT%X.000Z')

cpu_metric = '(sum(rate(container_cpu_usage_seconds_total{container=~"processor-."}[1m])) by (container) / (sum(container_spec_cpu_quota{container=~"processor-."}) by (container) / sum(container_spec_cpu_period{container=~"processor-."}) by (container))) * 100'
mem_metric = '(sum(container_memory_usage_bytes{container=~"processor-."}) by (container) / sum (container_spec_memory_limit_bytes{container=~"processor-."}) by (container)) * 100'

def parse_args():
    global args
    global since
    global until
        
    for i in range(1, (len(sys.argv) - 1), 2):
        args[sys.argv[i]] = sys.argv[i + 1]

    if '--since' in args:
        since = (datetime.now(timezone.utc) - timedelta(hours=int(args['--since']))).strftime('%G-%m-%dT%X.000Z')
    
    if '--until' in args:
        until = (datetime.now(timezone.utc) - timedelta(hours=int(args['--until']))).strftime('%G-%m-%dT%X.000Z')

    print(args)

def query_prometheus(query, resource_type):

    url = 'http://localhost:9090/api/v1/query_range'
    params = {
        'query': query,
        'start': since,
        'end': until,
        'step': '1s'
    }

    # save snapshot of metric values
    x = requests.post(url, headers={'Content-Type': 'application/x-www-form-urlencoded'}, data=params)
    f = open(resource_type + "_metrics.json", "w")
    f.write(x.text)
    f.close()

    response_json = json.loads(x.text)

    if 'error' in response_json:
        raise Exception('prometheus exception', response_json['error'])

    for result in response_json['data']['result']:
        container = result['metric']['container']
        
        for (timestamp, value) in result['values']:
            if timestamp not in container_resource:
                container_resource[timestamp] = {}

            if container not in container_resource[timestamp]:
                container_resource[timestamp][container] = {}

            container_resource[timestamp][container][resource_type] = value

def enrich_loadtest_metrics():
    # Read file containing loadtest results 
    df = pd.read_csv(args['--csv'], names=['timestamp', 'elapsed', 'message_type', 'instance_id', 'instance_type'])

    for i, row in df.iterrows():
        ts = int(row['timestamp'])

        if ts in container_resource:
            for res_type in ['cpu', 'mem']:
                for instance in ['a', 'b', 'c']:
                    # e.g.: df.at[i,'processor_a_cpu'] = float(container_resource[ts]['processor-a']['cpu'])
                    df.at[i,'_'.join(['processor', instance, res_type])] = float(container_resource[ts]['processor-' + instance][res_type])

        else:
            print('warn: ', ts, 'not covered by prometheus metrics')

    print(df.to_string())

    df.to_csv('metrics_' + args['--csv'])

if __name__ == '__main__':
    parse_args()

    query_prometheus(cpu_metric, 'cpu')
    query_prometheus(mem_metric, 'mem')

    enrich_loadtest_metrics()
