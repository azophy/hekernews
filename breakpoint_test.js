import { check } from 'k6';
import { htmlReport } from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";
import http from 'k6/http';

//const baseUrl = 'http://host.docker.internal:3000'
const baseUrl = 'http://localhost:3000'

function randomString(length) {
   var result           = '';
   var characters       = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
   var charactersLength = characters.length;
   for ( var i = 0; i < length; i++ ) {
      result += characters.charAt(Math.floor(Math.random() * charactersLength));
   }
   return result;
}

export function handleSummary(data) {
  var datetime = Date.now()
  var filename = "./reports/" + datetime + "_summary.html"
  console.log('filename: ' + filename)
  return {
    [filename] : htmlReport(data),
  };
}

export function testRegister(username) {
  let data =  {
    name: 'test',
    email: 'test@example.com',
    username: username,
    password: 'test',
  };
  let res = http.post(`${baseUrl}/register`, data);
  check(res, {
    'testRegister is status 200': (r) => r.status === 200,
  });
}

export function testLogin(username) {
  let data = {
    username: username,
    password: 'test',
  };
  let res = http.post(`${baseUrl}/login`, data);
  check(res, {
    'testLogin is status 200': (r) => r.status === 200,
  });
}

export function testLogout(params) {
  let res = http.get(`${baseUrl}/member/logout`);
  check(res, {
    'testLogout is status 200': (r) => r.status === 200,
  });
}

// =================================
export const options = {
  // threshold for breakpoint test
  //executor: 'ramping-arrival-rate', //Assure load increase if the system slows
  stages: [
    //{ duration: '10s', target: 10 }, // just slowly ramp-up to a HUGE load
    { duration: '15m', target: 1000 }, // just slowly ramp-up to a HUGE load
  ],
  thresholds: {
    //http_req_failed: [{ threshold: 'rate<0.05', abortOnFail: true, delayAbortEval: '10s' }],
    //http_req_duration: [{ threshold: 'p(95) < 750', abortOnFail: true, delayAbortEval: '10s' }],
    http_req_failed: [{ threshold: 'rate<0.1', abortOnFail: true, delayAbortEval: '10s' }],
    http_req_duration: [{ threshold: 'p(80) < 1000', abortOnFail: true, delayAbortEval: '10s' }],
  },
};

export default function () {
  //var username = randomString(20)
  //testRegister(username);
  testLogin('test');
  testLogout();
}
