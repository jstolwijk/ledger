install k6: `brew install k6`

Lets run a load test with 1000 concurrent users, with each user executing 10 requests:
`k6 run --vus 1000 --iterations 10000 test.js`
