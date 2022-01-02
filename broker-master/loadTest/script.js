import http from 'k6/http';
import { sleep } from 'k6';
import grpc from 'k6/net/grpc';

// export const options = {
//
//     vus: 5,
//
//     duration: '300s',
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

export const options = {
    vus: 5,
    duration: '300s',
};

export default () => {
    client.connect('localhost:5050', {plaintext: true});
    const publishRequest1 = {subject: "ali",
        body: JSON.stringify(JSON.stringify[110, 110]),
        expirationSeconds: 1,}

    const publishRequest2 = {subject: "hello",
        body: JSON.stringify(JSON.stringify[110, 110]),
        expirationSeconds: 1,}

    const publishRequest3 = {subject: "mahdi",
        body: JSON.stringify(JSON.stringify[110, 110]),
        expirationSeconds: 1,}


    const publishRequest4 = {subject: "amir",
        body: JSON.stringify(JSON.stringify[110, 110]),
        expirationSeconds: 1,}

    const publishRequest5 = {subject: "hi",
        body: JSON.stringify(JSON.stringify[110, 110]),
        expirationSeconds: 1,}

    const publishRequest6 = {subject: "ali2",
        body: JSON.stringify(JSON.stringify[110, 110]),
        expirationSeconds: 1,}

    const publishRequest7 = {subject: "hello2",
        body: JSON.stringify(JSON.stringify[110, 110]),
        expirationSeconds: 1,}

    const publishRequest8 = {subject: "mahdi2",
        body: JSON.stringify(JSON.stringify[110, 110]),
        expirationSeconds: 1,}


    const publishRequest9 = {subject: "amir2",
        body: JSON.stringify(JSON.stringify[110, 110]),
        expirationSeconds: 1,}

    const publishRequest10 = {subject: "hi2",
        body: JSON.stringify(JSON.stringify[110, 110]),
        expirationSeconds: 1,}

    //for (let i = 0; i < 100; i++) {
        const response1 = client.invoke('broker.Broker/Publish', publishRequest1)
        // console.log(JSON.stringify(response.message))
    //}

    //for (let i = 0; i < 100; i++) {
        const response2 = client.invoke('broker.Broker/Publish', publishRequest2)
        // console.log(JSON.stringify(response.message))
    //}

    //for (let i = 0; i < 100; i++) {
        const response3 = client.invoke('broker.Broker/Publish', publishRequest3)
        // console.log(JSON.stringify(response.message))
    //}

    //for (let i = 0; i < 100; i++) {
        const response4 = client.invoke('broker.Broker/Publish', publishRequest4)
        // console.log(JSON.stringify(response.message))
    //}

    //for (let i = 0; i < 100; i++) {
        const response5 = client.invoke('broker.Broker/Publish', publishRequest5)
        // console.log(JSON.stringify(response.message))
    //}





    //for (let i = 0; i < 100; i++) {
    const response6 = client.invoke('broker.Broker/Publish', publishRequest6)
    // console.log(JSON.stringify(response.message))
    //}

    //for (let i = 0; i < 100; i++) {
    const response7 = client.invoke('broker.Broker/Publish', publishRequest7)
    // console.log(JSON.stringify(response.message))
    //}

    //for (let i = 0; i < 100; i++) {
    const response8 = client.invoke('broker.Broker/Publish', publishRequest8)
    // console.log(JSON.stringify(response.message))
    //}

    //for (let i = 0; i < 100; i++) {
    const response9 = client.invoke('broker.Broker/Publish', publishRequest9)
    // console.log(JSON.stringify(response.message))
    //}

    //for (let i = 0; i < 100; i++) {
    const response10 = client.invoke('broker.Broker/Publish', publishRequest10)
    // console.log(JSON.stringify(response.message))
    //}

    client.close()

}