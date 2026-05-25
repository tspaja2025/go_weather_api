# Go Weather API

A simple weather API project built with Go. This API fetches and returns weather data from a 3rd party API.

## Roadmap.sh beginner project
This project was created as a part of weather API beginner project.
Check out the project details [roadmap.sh](https://roadmap.sh/projects/weather-api-wrapper-service)

---

## Installation

Clone the repository:

```bash
git clone https://github.com/tspaja2025/go_weather_api.git
cd go_weather_api
```

Run the application:

```bash
go run main.go
```

```bash
curl "http://localhost:3000/api/weather?city=London"
```

---

## Technologies Used

* Go
* Standard library packages:
  * `encoding/json`
  * `fmt`
  * `log`
  * `net/http`
  * `net/url`
  * `os`
  * `strings`
  * `time`
	* `github.com/joho/godotenv`

## Learning Goals

This project was built to practice:

* Working with 3rd party APIs
* Caching
* Environment variables
* API Structure
* Requests

---

## License

This project is open source and available under the MIT License.
