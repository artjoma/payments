# Payment processor
* Go 1.13
* PostgreSQL
* Docker

Initialization:
* `docker-compose up -d db`
* `cd database`
* `sh create.sh` - prepare database 
* `cd ..`
* `docker-compose stop`
* `docker-compose up -d`

Stop/start:
* `docker-compose up`
* `docker-compose stop`

For local development start only db:
`docker-compose start db`

Test:
```
curl -H "Content-Type: application/json" -H "Source-Type: client" --request POST --data '{"state": "win", "amount": "12.5", "transactionId": "0a732c81-3a38-4641-aef4-a4971e7b45vx"}' http://127.0.0.1:8095/payment/new
```

Load test: /src/client/load_test.go