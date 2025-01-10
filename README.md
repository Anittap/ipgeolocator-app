# IP Geolocation Application

This project is a scalable and efficient IP Geolocation API built using **Go**, **Flask**, **Memcached**, and **AWS Secrets Manager**. It is containerized with Docker for seamless deployment and scalability.

## Features

- **Go Backend**: Handles geolocation requests using the Gin framework and fetches data from an external API.
- **AWS Secrets Manager**: Securely manages API keys, ensuring sensitive data is never hardcoded.
- **Memcached Integration**: Implements caching for improved performance and reduced API call overhead.
- **Flask Frontend**: Provides a user-friendly interface for submitting IP addresses and viewing geolocation results.
- **Docker Containerization**: Ensures the app is portable and deployable across various environments.
- **Scalable Architecture**: Designed for high performance and easy scalability.

---

## Architecture

The application is divided into the following components:

1. **Backend (Go)**:
   - Built using the Gin framework.
   - Fetches geolocation data from an external API.
   - Implements caching with Memcached to store and retrieve results efficiently.

2. **Frontend (Flask)**:
   - A web interface to submit IP addresses and display geolocation information.
   - Validates user inputs before forwarding requests to the backend.

3. **Memcached**:
   - Caches geolocation results to minimize redundant external API calls and speed up response times.

4. **Secrets Management**:
   - Uses AWS Secrets Manager to securely fetch API keys via an attached IAM role instead of access key environment variables.

---

## Prerequisites

- **Docker**: Installed and running.
- **AWS Account**: Access to AWS Secrets Manager with a configured secret for the API key and an IAM role attached to the instance.
- **External Geolocation API Key**: Stored in AWS Secrets Manager.

---

## Getting Started

### Clone the Repository
```bash
git clone https://github.com/Anittap/ipgeolocator-app.git
cd ipgeolocator-app
```

### Build and Run the Application
Use Docker Compose to build and run the application:
```bash
docker-compose up --build
```

### Access the Application
- **Frontend**: Visit `http://localhost:80` to use the IP Geolocation app.
- **Backend**: Accessible via `http://localhost:8080` (used internally by the frontend).

---

## Note on AWS Region

The AWS region for Secrets Manager is hardcoded in the `docker-compose.yml` file as `us-east-1`:
```yaml
      - REGION_NAME=us-east-1
```
Ensure your AWS resources are configured in this region, or modify the value in the `docker-compose.yml` file to match your preferred region.

---

## Project Structure

```plaintext
.
├── Dockerfile         # Multi-stage Docker build for the Go backend
├── docker-compose.yml # Docker Compose configuration
├── frontend/          # Flask-based frontend code
│   ├── Dockerfile     # Dockerfile for the Flask frontend
│   ├── main.py        # Entry point for the Flask app
│   ├── requirements.txt # Python dependencies
│   └── templates/     # HTML templates for the frontend
│       ├── error.html # Error page template
│       └── index.html # Main page template
├── main.go            # Go-based backend code
```

---

## Key Technologies

- **Go**: High-performance backend with Gin framework.
- **Flask**: Lightweight and user-friendly frontend.
- **Memcached**: In-memory caching for improved API response times.
- **AWS Secrets Manager**: Securely manages sensitive API keys.
- **Docker**: Ensures portability and ease of deployment.
- **Docker Compose**: Simplifies multi-container application setup.

---

## Benefits

- **Caching**: Reduces external API calls, improving response times.
- **Security**: Sensitive data is securely managed with AWS Secrets Manager.
- **Scalability**: The containerized architecture ensures seamless scaling.
- **Performance**: Optimized with Memcached for low-latency responses.

---

## Future Improvements

- **Load Balancing**: Introduce load balancers for better scalability.
- **Kubernetes Deployment**: Enhance scalability with Kubernetes orchestration.
- **Rate Limiting**: Implement rate-limiting to prevent abuse of the API.

---

