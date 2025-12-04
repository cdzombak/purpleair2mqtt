# purpleair2mqtt

Bridge PurpleAir air quality sensors to MQTT and InfluxDB.

This program connects to local web server on a PurpleAir air quality monitor and publishes the data to an MQTT broker and optionally logs the data into an InfluxDB database for further analysis. 

**US EPA AQI Calculation**: This program also calculates the US EPA Air Quality Index (AQI) based on PM2.5 and PM10 concentrations reported by the PurpleAir sensor, following the official EPA guidelines.

## Background

Once you have your PurpleAir monitor connected to your local network, you can access it by just going to its IP address and making an HTTP request. This provides a user friendly page of what is going on with the monitor. If you'd like to get the data structured, make a request to `/json`.

This provides an excellent way to get real-time information from a local device without needing to manage API keys or even make calls out to the public internet. It's this JSON payload that the program parses.

## Installation

### Homebrew (macOS or Linux)

```shell
brew install cdzombak/oss/purpleair2mqtt
```

### Debian/Ubuntu via apt repository

Install my Debian repository if you haven't already:

```shell
sudo apt-get install ca-certificates curl gnupg
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://dist.cdzombak.net/deb.key | sudo gpg --dearmor -o /etc/apt/keyrings/dist-cdzombak-net.gpg
sudo chmod 0644 /etc/apt/keyrings/dist-cdzombak-net.gpg
echo -e "deb [signed-by=/etc/apt/keyrings/dist-cdzombak-net.gpg] https://dist.cdzombak.net/deb/oss any oss\n" | sudo tee -a /etc/apt/sources.list.d/dist-cdzombak-net.list > /dev/null
sudo apt-get update
```

Then install `purpleair2mqtt`:

```shell
sudo apt-get install purpleair2mqtt
```

### Docker

Docker images are available for `linux/amd64`, `linux/arm64`, `linux/386`, `linux/arm/v7`, and `linux/arm/v6`, from both GHCR and Docker Hub:

```shell
docker pull ghcr.io/cdzombak/purpleair2mqtt:1
# or:
docker pull cdzombak/purpleair2mqtt:1
```

### Build from source

See "Building the Application" below.

## Configuration

- This program uses TOML as its configuration format.
- You'll need to figure out the IP address or hostname of your Purple Air monitor on your local network. I recommend assigning it a static IP in your router.

First, define the location of the device and the polling rate in seconds. I find that you don't need to read it much more than once every two minutes. This will allow the application to pull the information from the Purple Air sensor.

```toml
[purpleair]
    url = "http://192.168.1.24/json"
    poll_rate = 120

```

Next, define the information needed for wherever your MQTT broker is running. You'll need the hostname. If your MQTT broker requires authentication, you can optionally specify `broker_username` and `broker_password`. If `topic_prefix` is left as null it will default to `airquality` and if `topic` is left as null it will default to the `geo` identifier of your PurpleAir sensor.

```toml
[mqtt]
    broker_host = "mqttbroker.local"
    broker_port = 1883
    broker_username = ""  # optional
    broker_password = ""  # optional
    client_id = "purpleair2mqtt"
    topic_prefix = "airquality"
    topic = ""
```

If you want Home Assistant integration, this chunk _should_ work to support Home Assistant auto discovery from MQTT.

```toml
[hass]
    discovery = true
    discovery_prefix = "homeassistant"
    device_model = "pa-sd-ii"
    device_name = "pa-sd-ii"
    # if you don't set object_id then you'll get end up with the MAC as your id
    object_id = "pa-sd-ii"
```

Finally, if you'd like to use the native InfluxDB integration, this section should work for you. You'll need to supply the `hostname`, create the database, which defaults to `purpleair` and define the username and password to write to that database.

```toml
[influx]
    hostname = "influxdb.local"
    port = 8086
    database = "purpleair"
    username = "YOUR_USERNAME"
    password = "YOUR_PASSWORD"
```

## Building the Application

To build for the current platform:

```bash
make build
```

The binary will be placed in `./out/purpleair2mqtt`.

To build for all supported platforms:

```bash
make all
```

To build and create Debian packages (requires [fpm](https://fpm.readthedocs.io)):

```bash
make package
```

Run `make help` to see all available targets.

## Running the Application

```bash
./purpleair2mqtt -config config.toml
```

### Command-Line Arguments

- `-config <file>`: Path to the TOML configuration file (required)
- `-version`: Print version and exit

## Running with Docker

Docker images are published to both Docker Hub and GHCR. See the "Installation" section above for pull commands.

### Running the Container

You'll need to mount a configuration file into the container. The entrypoint expects the config at `/config.toml`:

```bash
docker run -d \
  --name purpleair2mqtt \
  --restart unless-stopped \
  --network host \
  -v /path/to/your/config.toml:/config.toml:ro \
  cdzombak/purpleair2mqtt:1
```

### Docker Compose Example

```yaml
services:
  purpleair2mqtt:
    image: cdzombak/purpleair2mqtt:1
    container_name: purpleair2mqtt
    restart: unless-stopped
    network_mode: host
    volumes:
      - ./config.toml:/config.toml:ro
```

### Building the Container Locally

```bash
docker build \
  --build-arg BIN_VERSION=$(./.version.sh) \
  -t purpleair2mqtt:local .
```

## US EPA AQI Calculation

The program calculates the US EPA Air Quality Index (AQI) based on the PM2.5 and PM10 concentration values reported by the PurpleAir sensor. The calculation follows the official EPA breakpoints and formulas as specified in the [Technical Assistance Document for the Reporting of Daily Air Quality](https://document.airnow.gov/technical-assistance-document-for-the-reporting-of-daily-air-quailty.pdf).

The following AQI values are calculated and published:

- **Overall EPA AQI**: The highest AQI value between PM2.5 and PM10
- **PM2.5 AQI**: AQI calculated from PM2.5 concentration
- **PM10 AQI**: AQI calculated from PM10 concentration
- **AQI Category**: Good, Moderate, Unhealthy for Sensitive Groups, Unhealthy, Very Unhealthy, or Hazardous
- **AQI Color**: English color name (Green, Yellow, Orange, Red, Purple, Maroon)
- **AQI Color RGB**: RGB color value (e.g., `rgb(0,228,0)` for Good, `rgb(255,0,0)` for Unhealthy)

These values are published via MQTT and stored in InfluxDB alongside the existing PurpleAir data.

## MQTT Topics

The application publishes data to the following MQTT topics (assuming default `airquality` prefix):

**Status Topics** (overall sensor values):
- `airquality/{sensor_name}/EPAAQI` - US EPA AQI value (highest of PM2.5 and PM10)
- `airquality/{sensor_name}/EPAPM25AQI` - US EPA PM2.5 AQI
- `airquality/{sensor_name}/EPAPM10AQI` - US EPA PM10 AQI
- `airquality/{sensor_name}/EPAAQICategory` - AQI category (e.g., "Good", "Moderate")
- `airquality/{sensor_name}/EPAAQIColor` - AQI color name (e.g., "Green", "Yellow")
- `airquality/{sensor_name}/EPAAQIColorRGB` - AQI color as RGB string (e.g., `rgb(0,228,0)`)

**Individual Sensor Topics** (for sensor A and B):
- `airquality/{sensor_name}/sensor_A/epa_aqi` - EPA AQI for sensor A
- `airquality/{sensor_name}/sensor_A/epa_pm25_aqi` - EPA PM2.5 AQI for sensor A
- `airquality/{sensor_name}/sensor_A/epa_pm10_aqi` - EPA PM10 AQI for sensor A
- `airquality/{sensor_name}/sensor_A/epa_aqi_category` - AQI category for sensor A
- `airquality/{sensor_name}/sensor_A/epa_aqi_color` - AQI color name for sensor A
- `airquality/{sensor_name}/sensor_A/epa_aqi_color_rgb` - AQI color RGB for sensor A
- (Same topics available for sensor_B)

All existing PurpleAir data topics remain unchanged.

## InfluxDB Schema

The application writes to two measurements in InfluxDB. The measurement names are configurable (see Configuration section).

### `purpleair_status` Measurement

This measurement contains overall sensor status data.

**Tags:**
- `sensorId` - MAC address of the sensor

**Fields:**
- `temperature` - Temperature in Fahrenheit
- `humidity` - Relative humidity percentage
- `pressure` - Atmospheric pressure in mmHg
- `dewpoint` - Dewpoint in Fahrenheit
- `rssi` - WiFi signal strength
- `epa_aqi` - US EPA AQI value (highest of PM2.5 and PM10)
- `epa_pm25_aqi` - US EPA PM2.5 AQI
- `epa_pm10_aqi` - US EPA PM10 AQI
- `epa_aqi_category` - AQI category string (e.g., "Good", "Moderate")
- `epa_aqi_color` - AQI color name (e.g., "Green", "Yellow")
- `epa_aqi_color_rgb` - AQI color RGB value (e.g., "rgb(0,228,0)")

### `purpleair_monitor` Measurement

This measurement contains per-sensor particle data (one entry each for sensor A and B).

**Tags:**
- `sensorId` - MAC address of the sensor
- `sensor` - Sensor identifier ("A" or "B")

**Fields:**
- `pm2.5_aqic` - PurpleAir's AQI color
- `pm2.5_aqi` - PurpleAir's AQI value
- `pm1.0_cf_1`, `pm2.5_cf_1`, `pm10.0_cf_1` - CF=1 PM values
- `pm1.0_atm`, `pm2.5_atm`, `pm10.0_atm` - ATM PM values
- `pm0.3_um`, `pm0.5_um`, `pm1.0_um`, `pm2.5_um`, `pm5.0_um`, `pm10.0_um` - Particle counts
- `key1_response`, `key2_response`, etc. - Response metrics
- `epa_aqi` - US EPA AQI value
- `epa_pm25_aqi` - US EPA PM2.5 AQI
- `epa_pm10_aqi` - US EPA PM10 AQI
- `epa_aqi_category` - AQI category string
- `epa_aqi_color` - AQI color name
- `epa_aqi_color_rgb` - AQI color RGB value

## Authors & License

Copyright (c) 2022 [Patrick Wagstrom](https://github.com/pridkett); modifications (c) 2025 [Chris Dzombak](https://github.com/cdzombak)

Licensed under the terms of the MIT License; see [LICENSE](LICENSE) in this repo.
