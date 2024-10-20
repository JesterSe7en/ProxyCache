# ProxyCache

## Overview 📖

The ProxyCache CLI Tool is a lightweight command-line application designed to act as an intermediary between clients and external web services. It enhances response times and reduces load on the external services by leveraging a caching mechanism powered by Redis. The tool intelligently stores and retrieves HTTP responses based on client requests, ensuring optimal performance and efficient resource management.

## Features 🌟

- **Caching Layer**: Automatically caches HTTP responses for frequently accessed resources, minimizing redundant requests to external services.
- **Redis Integration**: Utilizes Redis as a fast, in-memory key-value store to manage cached responses, ensuring quick access and efficient data retrieval.
- **Expiration Management**: Configurable expiration times for cached entries to prevent serving stale data while optimizing memory usage.
- **Error Handling**: Robust error management to gracefully handle scenarios where the external service or Redis is unavailable.
- **Command-Line Interface**: User-friendly CLI for easy configuration and interaction with the caching service.

## Installation 🛠️

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/JesterSe7en/ProxyCache.git
   cd cache-proxy-cli
   ```
2. ** Install Dependencies**: Ensure you ahve Go installed, then run:
   ```bash
   go mod tidy
   ```
3. Set up Redis: Install Redis and ensure it is running on your local machine or configure it to connect to a remote Redis server.


## Usage 📝

### Prerequisites
 Ensure that the environment variables `REDIS_URL` and `REDIS_PASSWORD` are set correctly.

Run the ProxyCache CLI tool with the following command:
```bash
cacheProxy --port <port> --redirectURL <redirectURL>
```

### Configuration Options
- --port: The port on which the server will listen for incoming requests. Default is 6379.
- --redirectURL: The URL of the external service to be proxied.

## How It Works 🔎
1. **Request Handling**: When a client sends a request, the tool checks if the response is already cached in Redis.
2. **Cache Lookup**: If a valid cached response is found, it is returned immediately, bypassing the external service.
3. **Fetching from External Service**: If no valid response is cached, the tool forwards the request to the external service, caches the response in Redis, and returns it to the client.
4. **Expiration Management**: Cached entries are automatically removed based on configured expiration settings to ensure freshness.

## Why This Project? 🤔
  -  Backend Engineering Skills: This project demonstrates my ability to design and implement scalable backend systems using Go and Redis.
  -  Understanding of Caching Mechanisms: It showcases my knowledge of caching strategies, which are crucial for optimizing performance in web applications.
  -  Error Handling Proficiency: The project emphasizes my approach to robust error management and handling in backend applications.

## Future Plans 🚀

- Add support for other caching mechanisms, such as Memcached or DynamoDB.
- Add support for other web services, such as GraphQL or RESTful APIs.
- Add support for rate limiting and throttling.

## License 📜

This project is licensed under the MIT License. See the LICENSE file for more details.

## Acknowledgements 🙏

- [Redis](https://redis.io/)
- [Go](https://golang.org/)
- [Docker](https://www.docker.com/)
- [Roadmap.sh](https://roadmap.sh/)
