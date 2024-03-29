# MQTT Subscriber Application

## Description

This is a Go application that connects to an MQTT broker, subscribes to a specific topic, and saves the received messages to a PostgreSQL database. It also provides a REST API to retrieve the stored messages.

## Installation

To install the application, follow these steps:

1. Clone the repository: `git clone <repository-url>`
2. Navigate to the project directory: `cd <project-directory>`

## Usage

To run the application, follow these steps:

1. Build the Docker image: `docker-compose up --build`
2. Access the REST API at `http://localhost:8080/messages`


## Dependencies

The application relies on the following dependencies:

1. Gin: Web framework for building RESTful APIs.
2. Eclipse Paho MQTT Client: MQTT client library for Go.
3. Go-PG: PostgreSQL ORM for Go.
