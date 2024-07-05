docker build -t gateway:latest -f ../gateway/Dockerfile ../gateway

docker build -t stocks:latest -f ../stocks/Dockerfile ../stocks

docker build -t order:latest -f ../order/Dockerfile ../order