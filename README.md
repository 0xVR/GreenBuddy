# GreenBuddy
GreenBuddy is a smart plant monitoring system that uses a custom soil moisture sensor and a Raspberry Pi to track soil moisture levels. The system sends real-time moisture data to a cloud-based serverless architecture built on AWS, providing automated notifications to users for plant care.

## Features
- Custom-built capacitive soil moisture sensor.
- Real-time soil moisture data collection and monitoring.
- Serverless cloud architecture using AWS Lambda, SNS, and DynamoDB.
- Automated notifications when plants need watering.
- Scalable and low-cost solution for efficient plant care.

## Hardware Used:
- Raspberry Pi
- Capacitive soil moisture sensor
- MCP3008 (Analog to Digital Converter)
- Jumper wires and breadboard

## Installation and Usage

Hardware setup has been omitted for simplicity.

1. Move `plant_monitor.py` to the Raspberry Pi through the repo:

```sh
git clone https://github.com/0xVR/GreenBuddy.git
```

2. Install the required packages.

```sh
cd GreenBuddy && pip install -r requirements.txt
```

3. Setup AWS Lambda, SNS, and DynamoDB using Terraform:

```sh
cd terraform && terraform init && terraform apply
```

4. Update the Python script with the links to your AWS API Gateway and SNS topic

5. Subscribe to the SNS topic on AWS

6. Done!