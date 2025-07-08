```bash
docker run -d --name redis-stack-server -p 6379:6379 redis/redis-stack-server:latest
```

In WSL:
```bash
sudo apt update
sudo apt install redis-tools
redis-cli ping
```