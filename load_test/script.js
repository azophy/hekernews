import { check } from 'k6';
import http from 'k6/http';

const baseUrl = 'http://172.17.0.1:3000'

export function testRegister(params) {
  let data = params || {
    name: 'test',
    email: 'test@example.com',
    username: 'test',
    password: 'test',
  };
  let res = http.post(`${baseUrl}/register`, data);
  //check(res, {
    //'is status 200': (r) => r.status === 200,
  //});
}

export function testLogin(params) {
  let data = params || {
    username: 'test',
    password: 'test',
  };
  let res = http.post(`${baseUrl}/login`, data);
  //check(res, {
    //'is status 200': (r) => r.status === 200,
  //});
}

export function testLogout(params) {
  http.post(`${baseUrl}/member/logout`);
}

// =================================
export const options = {
  stages: [
    { duration: '30s', target: 20 },
    //{ duration: '1m30s', target: 10 },
    //{ duration: '20s', target: 0 },
  ],
};

export default function () {
  testRegister();
  testLogin();
  testLogout();
}

