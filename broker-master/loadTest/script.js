import http from 'k6/http';
import { sleep } from 'k6';
import grpc from 'k6/net/grpc';

// export const options = {
//
//     vus: 10,
//
//     duration: '30s',
//
// };
//
//
// export default function () {
//
//     http.get('http://localhost:5050');
//
//     sleep

const client = new grpc.Client();
client.load(['../api/proto'], 'broker.proto');

export default () => {
    client.connect('localhost:5050', {plaintext: true});
    const publishRequest = {subject: "ali",
        body: JSON.stringify(JSON.stringify[110, 110]),
        expirationSeconds: 100000000,}

    for (let i = 0; i < 10000; i++) {
        const response = client.invoke('broker.Broker/Publish', publishRequest)
        console.log(JSON.stringify(response.message))
    }

    client.close()
    sleep(1)

}