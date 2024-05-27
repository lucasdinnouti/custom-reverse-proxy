import requests
import json
import pandas as pd


container_resource = {}

instance_type_resources = {
    'small': { 'cpu': 1, 'mem': 16777216 },
    'medium': { 'cpu': 1, 'mem': 33554432 },
    'large': { 'cpu': 1, 'mem': 67108864 }
}

# Query Prometheus for cpu usage over time 

url = 'http://localhost:9090/api/v1/query_range'
params = {
    'query': 'container_cpu_usage_seconds_total{container=~"processor-."}',
    'start': '2024-05-23T00:30:00.000Z',
    'end': '2024-05-23T01:30:00.000Z',
    'step': '1s'
}

x = requests.post(url, headers={'Content-Type': 'application/x-www-form-urlencoded'}, data=params)

response_json = json.loads(x.text)

# save snapshot of metric values
f = open("cpu_metrics.json", "w")
f.write(response_json)
f.close()

for result in response_json['data']['result']:
    container = result['metric']['container']
    
    for (timestamp, value) in result['values']:
        if timestamp not in container_resource:
            container_resource[timestamp] = {}

        if container not in container_resource[timestamp]:
            container_resource[timestamp][container] = {}

        container_resource[timestamp][container]['cpu'] = value


# Query Prometheus for mem usage over time 

url = 'http://localhost:9090/api/v1/query_range'
params = {
    'query': 'container_memory_usage_bytes{container=~"processor-."}',
    'start': '2024-05-23T00:30:00.000Z',
    'end': '2024-05-23T01:30:00.000Z',
    'step': '1s'
}

x = requests.post(url, headers={'Content-Type': 'application/x-www-form-urlencoded'}, data=params)

response_json = json.loads(x.text)

# save snapshot of metric values
f = open("mem_metrics.json", "w")
f.write(response_json)
f.close()

for result in response_json['data']['result']:
    container = result['metric']['container']
    
    for (timestamp, value) in result['values']:
        if timestamp not in container_resource:
            container_resource[timestamp] = {}

        if container not in container_resource[timestamp]:
            container_resource[timestamp][container] = {}

        container_resource[timestamp][container]['mem'] = value

# Read file containing loadtest results 

def normalize_resource(instance_type, resource_type, amount):
    return(amount / instance_type_resources[instance_type][resource_type])

df = pd.read_csv('result_10.csv', names=['timestamp', 'elapsed', 'message_type', 'instance'])

for i, row in df.iterrows():
    ts = int(row['timestamp'])
    inst_type = row['instance'].split(sep='-')[0]
    print('instance_type', inst_type)

    df.at[i,'processor_a_cpu'] = normalize_resource(inst_type, 'cpu', container_resource[ts]['processor-a']['cpu'])
    df.at[i,'processor_b_cpu'] = normalize_resource(inst_type, 'cpu', container_resource[ts]['processor-b']['cpu'])
    df.at[i,'processor_a_mem'] = normalize_resource(inst_type, 'mem', container_resource[ts]['processor-a']['mem'])
    df.at[i,'processor_b_mem'] = normalize_resource(inst_type, 'mem', container_resource[ts]['processor-b']['mem'])

print(df.to_string()) 

df.to_csv('result_10_r.csv')





# print(container_resource)
