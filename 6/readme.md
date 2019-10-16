docker run -p 3306:3306 --name some-mysql -e MYSQL_ROOT_PASSWORD=1234 -d mysql:5
docker run -p 6379:6379 --name some-redis -d redis
docker run -p 11211:11211 --name my-memcache -d memcached
docker run -p 5672:5672 -d --hostname my-rabbit --name some-rabbit rabbitmq:3
docker run -p 27017:27017 --name some-mongo -d mongo

docker build --tag=mytnt .
docker run --name mytnt-inst -p 3301:3301 -d mytnt
docker exec -t -i mytnt-inst console