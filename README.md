# minitwitgo
Minitwit build with Go

Find dependencies: go list -f '{{ join .Imports "\n" }}'


Docker: 

docker build . -t minitwit
docker-compose up -d