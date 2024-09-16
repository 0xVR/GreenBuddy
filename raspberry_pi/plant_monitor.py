import time
import spidev
import RPi.GPIO as GPIO
import requests

# AWS API Gateway endpoint for the Lambda function
api_endpoint = 'https://<your-api-id>.execute-api.<region>.amazonaws.com/trigger'

# Setup for MCP3008
spi = spidev.SpiDev()
spi.open(0, 0)
spi.max_speed_hz = 1350000

def read_channel(channel):
    adc = spi.xfer2([1, (8 + channel) << 4, 0])
    data = ((adc[1] & 3) << 8) + adc[2]
    return data

def get_moisture_level(sensor_value):
    # Calculate voltage ratio to determine soil moisture %
    dry_value = 3.3
    wet_value = 5
    moisture_level = (sensor_value - wet_value) / (dry_value - wet_value) * 100
    return max(0, min(100, moisture_level))

def send_to_aws(moisture_data):
    payload = {'moisture_level': moisture_data}
    try:
        response = requests.post(api_endpoint, json=payload)
        response.raise_for_status()
        print(f"Data sent successfully: {response.status_code}")
    except requests.exceptions.HTTPError as err:
        print(f"Error sending data: {err}")

# Main loop
try:
    while True:
        sensor_value = read_channel(0)  # Assuming the sensor is connected to channel 0
        moisture_level = get_moisture_level(sensor_value)
        print(f"Moisture level: {moisture_level}%")

        send_to_aws(moisture_level)

        # Wait 10 minutes before next reading
        time.sleep(600)

except KeyboardInterrupt:
    print("Program stopped by User")
finally:
    spi.close()
    GPIO.cleanup()