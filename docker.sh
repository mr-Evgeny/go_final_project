docker build . --tag todo
docker run -v ./db:/app/db -p 7540:7540 todo
