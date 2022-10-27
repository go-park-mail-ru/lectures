docker build -f Dockerfile -t docker-golang-image .
docker build -f Dockerfile.Multistage -t docker-golang-image .
docker-compose up
