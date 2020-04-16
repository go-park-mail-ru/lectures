docker build -t docker_0 .

docker run -it docker_0

docker run -it -p 8080:8080 docker_0

docker-compose -f docker-compose.yml up -f

docker-compose -f docker-compose.yml up
docker-compose -f docker-compose.yml down



docker-compose -f docker-compose.yml down && docker-compose -f docker-compose.yml up


https://stackoverflow.com/questions/29145370/how-can-i-initialize-a-mysql-database-with-schema-in-a-docker-container
https://forums.docker.com/t/mysql-create-database-and-import-sql-file-with-dockerfile/30300

