import requests
import json
import pandas as pd


container_resource = {}

# Query Prometheus for cpu usage over time 

url = 'http://localhost:9090/api/v1/query_range'
params = {
    'query': '(sum(rate(container_cpu_usage_seconds_total{container=~"processor-."}[1m])) by (container) / (sum(container_spec_cpu_quota{container=~"processor-."}) by (container) / sum(container_spec_cpu_period{container=~"processor-."}) by (container))) * 100',
    'start': '2024-05-28T00:00:00.000Z',
    'end': '2024-05-28T02:30:00.000Z',
    'step': '1s'
}

# save snapshot of metric values
x = requests.post(url, headers={'Content-Type': 'application/x-www-form-urlencoded'}, data=params)
f = open("cpu_metrics.json", "w")
f.write(x.text)
f.close()

response_json = json.loads(x.text)

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
    'query': '(sum(container_memory_usage_bytes{container=~"processor-."}) by (container) / sum (container_spec_memory_limit_bytes{container=~"processor-."}) by (container)) * 100',
    'start': '2024-05-28T00:00:00.000Z',
    'end': '2024-05-28T02:30:00.000Z',
    'step': '1s'
}

# save snapshot of metric values
x = requests.post(url, headers={'Content-Type': 'application/x-www-form-urlencoded'}, data=params)
f = open("mem_metrics.json", "w")
f.write(x.text)
f.close()

response_json = json.loads(x.text)

for result in response_json['data']['result']:
    container = result['metric']['container']
    
    for (timestamp, value) in result['values']:
        if timestamp not in container_resource:
            container_resource[timestamp] = {}

        if container not in container_resource[timestamp]:
            container_resource[timestamp][container] = {}

        container_resource[timestamp][container]['mem'] = value

# Read file containing loadtest results 
df = pd.read_csv('result_100.csv', names=['timestamp', 'elapsed', 'message_type', 'instance_id', 'instance_type'])

for i, row in df.iterrows():
    ts = int(row['timestamp'])

    print(ts, i)
    df.at[i,'processor_a_cpu'] = float(container_resource[ts]['processor-a']['cpu'])
    df.at[i,'processor_b_cpu'] = float(container_resource[ts]['processor-b']['cpu'])
    df.at[i,'processor_c_cpu'] = float(container_resource[ts]['processor-c']['cpu'])
    df.at[i,'processor_a_mem'] = float(container_resource[ts]['processor-a']['mem'])
    df.at[i,'processor_b_mem'] = float(container_resource[ts]['processor-b']['mem'])
    df.at[i,'processor_c_mem'] = float(container_resource[ts]['processor-c']['mem'])

print(df.to_string()) 

df.to_csv('result_100_r.csv')





# print(container_resource)
