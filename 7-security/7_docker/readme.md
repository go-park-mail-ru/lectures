docker build -f Dockerfile -t docker-golang-image .
docker build -f Dockerfile.Multistage -t docker-golang-image .
docker-compose up

docker run --rm -p 8080:8080 docker-golang-image

<http://127.0.0.1:8080/>
<http://127.0.0.1:8080/images/sad.jpeg>
<http://127.0.0.1:8080/albums/selfies?user_id=10>
