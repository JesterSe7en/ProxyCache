# ProxyCache

## Overview 📖

The ProxyCache CLI Tool is a lightweight command-line application designed to act as an intermediary between clients and external web services. It enhances response times and reduces the load on these services by leveraging a caching mechanism powered by Redis. The tool intelligently stores and retrieves HTTP responses based on client requests, ensuring optimal performance and efficient resource management.

Additionally, it implements a rate throttling mechanism using token buckets, allowing for controlled request rates to the external APIs. This ensures fair usage and prevents overloading the services. With its flexibility, the ProxyCache tool can proxy requests to any web API, seamlessly forwarding client requests to the specified redirect URL while maintaining performance and reliability.

## Features 🌟

- **Caching Layer**: Automatically caches HTTP responses for frequently accessed resources via Redis, minimizing redundant requests to external services.
- **Redis Integration**: Utilizes Redis as a fast, in-memory key-value store to manage cached responses, ensuring quick access and efficient data retrieval.
- **Expiration Management**: Manages expiration times for cached entries to prevent serving stale data while optimizing memory usage.
- **Error Handling**: Robust error management to gracefully handle scenarios where the external service or Redis is unavailable.
- **Command-Line Interface**: CLI for easy configuration and interaction with the caching service.
- **Rate Limiting and Throttling**: Supports rate limiting and throttling mechanisms (token bucket) to ensure optimal throughput and resource utilization.

## Installation 🛠️

1. **Clone the Repository**:
   ```powershell
   git clone https://github.com/JesterSe7en/ProxyCache.git
   cd ProxyCache
   ```
2. **Install Dependencies**: Ensure you have Go installed, then run:
   ```powershell
   go mod tidy
   ```
3. **Set up Redis**: Install Redis and ensure it is running on your local machine or configure it to connect to a remote Redis server.  This tool assumes Redis instance is running on port 6379.  Optionally, you can create a Redis docker container.  Here's a docker compose config I used on my raspberry pi.  You can use any other Redis instance.
   ```yaml
   version: '3.8'
   services:
   redis:
      image: redis:latest
      restart: always
      ports:
         - '6379:6379'
      command: redis-server --save 60 1 --loglevel warning --requirepass <some_redis_password>
      volumes:
         - ./data:/data
   ```
4. **Build the Binary**:
   ```powershell
   go build -o ProxyCache.exe
   ```



## Usage 📝

### Prerequisites
 Ensure that the environment variables `REDIS_URL` and `REDIS_PASSWORD` are set correctly.
 Also, ensure that your Redis instance is running on port 6379.

Run the ProxyCache CLI tool with the following command:
```powershell
./ProxyCache.exe -port <port> -redirectURL <redirectURL>
```

### Configuration Options
- -port: The port on which the server will listen for incoming requests.
- -redirectURL: The URL of the external service to be proxied.

## How It Works 🔎
1. **Request Handling**: When a client sends a request, the tool checks if the response is already cached in Redis.
2. **Cache Lookup**: If a valid cached response is found, it is returned immediately, bypassing the external service.
3. **Fetching from External Service**: If no valid response is cached, the tool forwards the request to the external service, caches the response in Redis, and returns it to the client.
4. **Expiration Management**: Cached entries are automatically removed  after 24 hours to ensure freshness.

## Why This Project? 🤔
  -  Backend Engineering Skills: This project demonstrates my ability to design and implement scalable backend systems using Go and Redis.
  -  Understanding of Caching Mechanisms: It showcases my knowledge of caching strategies, which are crucial for optimizing performance in web applications.
  -  Error Handling Proficiency: The project emphasizes my approach to robust error management and handling in backend applications.

## Future Plans 🚀

- Add support for other caching mechanisms, such as Memcached.
- Add support for custom Redis configuration such as port and data expiration.
- ~~Add support for rate limiting and throttling.~~  *Added on Oct 21, 2024*

## License 📜

This project is licensed under the MIT License. See the LICENSE file for more details.

## Acknowledgements 🙏

- [Redis](https://redis.io/)
- [Go](https://golang.org/)
- [Docker](https://www.docker.com/)
- [Roadmap.sh](https://roadmap.sh/)
