HekerNews
=========

HackerNews clone for personal exercise in high performance backend service

based on:
- go
- Labstack's Echo framework
- MVP.css
- AlpineJs

## Running load test
```
$ docker run --rm -i grafana/k6 run - <load_test/script.js
$ docker run --rm \
  -p 5665:5665 \
  -i ghcr.io/grafana/xk6-dashboard:0.7.2 run \
  --out web-dashboard 
  - <load_test/script.js
```
