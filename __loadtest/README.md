# Load Test

## Dependences
- K6 (https://grafana.com/docs/k6/latest/set-up/install-k6/)

## Run
Pack create, cancel and update status flow:
```
k6 run packs.js
```

Pack events:
```
k6 run pack_events.js
```

## Config
**Service URL**: To change the service URL modify the `baseURL` constant in the ./config.js file.
