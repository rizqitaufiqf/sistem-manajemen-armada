# Sistem Manajemen Armada
## Technologies Used
*  **Golang**: For backend service development.
*  **MQTT (Eclipse Mosquitto)**: For receiving vehicle location data.
*  **PostgreSQL**: For storing vehicle location data.
*  **RabbitMQ**: For processing geofence-based events.
*  **Docker**: For running the entire system in containers.
## How to Run the Application
Ensure you have Docker and Docker Compose installed on your system.
1.  **Clone Repository:**
```bash
git clone https://github.com/rizqitaufiqf/sistem-manajemen-armada
cd sistem-manajemen-armada
```
2.  **Build and Run Services with Docker Compose:**
```bash
docker-compose up -d
```
3.  **Run MQTT Publisher:**
You can run the MQTT publisher from within the container or from your local host. 
> **Note:** For default, MQTT Publisher is already running in your docker app.

To run MQTT Publisher from your local host: 
1. ensure the Mosquitto broker is accessible from outside docker 
2. Set up your environment configuration:
    - rename `env.example` to `.env` 
    - Adjust the credentials and other configuration values as needed inside the `.env` file
3. enable `godenv` in `cmd/mqtt_publisher/main.go`
```
godotenv.Load()
```
4. run command below:
```bash
go run cmd/mqtt_publisher/main.go
```
4.  **Monitor your MQTT Subscriber and Your RabbitMQ Worker**
##### MQTT Subscriber
```
docker logs -f mqtt_subscriber
```
##### RabbitMQ Worker
```
docker logs -f rabbitmq_worker
```
## Testing the API
You can test the API endpoints using Postman.
### 1. Import Postman Collection:
Import file `sistem-manajemen-armada.postman_collection.json` file into your Postman.
### 2. API Endpoints:
#### Get a vehicle last location:
**Request:**
`GET http://localhost:8080/vehicles/:vehicle_id/location`

**Example:** `http://localhost:8080/vehicles/B1234XYZ/location`
#### Get vehicle location history within a specific time range:
**Request:**
`GET http://localhost:8080/vehicles/:vehicle_id/history?start=:timestamp_start&end=:timestamp_end`
**Parameters:**
- `start`: Start timestamp (UNIX epoch time)
- `end`: End timestamp (UNIX epoch time)

**Example:** `http://localhost:8080/vehicles/B1234XYZ/history?start=1715000000&end=1715009999` (Adjust timestamps according to your sent data)
> **Note:** You can find working examples with sample data in the imported Postman collection. Adjust the timestamps according to your test data.