
### Note
run the go program in WSL aswell

also, for some reason, i need to run the program before running redis-cli ping (idk wtf)
### Running Redis in WSL

### 1. Start Redis Container
```bash
docker run -d --name redis-stack-server -p 6379:6379 redis/redis-stack-server:latest
```

### 2. Install Redis Tools
```bash
sudo apt update
sudo apt install redis-tools
```

### 3. Verify Connection
Important: You need to run your Go program first before checking the connection (reason unknown).

```bash
redis-cli ping
```

### 🛑 Stopping the Redis Container
To stop the container:
```bash
docker stop redis-stack-server
```

To remove the container:
```bash
docker rm redis-stack-server
```
