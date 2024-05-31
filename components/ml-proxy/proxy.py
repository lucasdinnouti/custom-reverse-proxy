import onnxruntime as rt
import numpy as np

import http.server
import http.client

TARGETS = ['a', 'b', 'c']
TARGET_PORT = 8083

sess = rt.InferenceSession("test_2.onnx", providers=["CPUExecutionProvider"])
input_name = sess.get_inputs()[0].name
label_name = sess.get_outputs()[0].name

def select_target():
    message_type = 0
    processor_a_cpu = 4.346822
    processor_b_cpu = 6.074483
    processor_c_cpu = 14.455157
    processor_a_mem = 10.81543
    processor_b_mem = 10.229492
    processor_c_mem = 18.200684

    input_data = [[message_type, 0, processor_a_cpu, processor_b_cpu, processor_c_cpu, processor_a_mem, processor_b_mem, processor_c_mem],
                  [message_type, 1, processor_a_cpu, processor_b_cpu, processor_c_cpu, processor_a_mem, processor_b_mem, processor_c_mem],
                  [message_type, 2, processor_a_cpu, processor_b_cpu, processor_c_cpu, processor_a_mem, processor_b_mem, processor_c_mem]]
    
    lat_predictions = sess.run([label_name], {input_name: input_data})
    index_of_min = np.argmin(lat_predictions)
    selected = TARGETS[index_of_min]

    return f'http://processor-{selected}.default.svc.cluster.local'

class ProxyHandler(http.server.BaseHTTPRequestHandler):
    def do_POST(self):
        self.handle_request('POST')

    def handle_request(self, method):
        target = select_target()

        conn = http.client.HTTPConnection(target, TARGET_PORT)
        conn.request(method, self.path, headers=self.headers)

        response = conn.getresponse()
        self.send_response(response.status)
        
        for header, value in response.getheaders():
            self.send_header(header, value)
        self.end_headers()
        
        self.wfile.write(response.read())
        
        conn.close()


server_address = ('', 8082)
httpd = http.server.HTTPServer(server_address, ProxyHandler)
print('Reverse proxy server running on port 8082...')
httpd.serve_forever()