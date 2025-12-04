purpleair2mqtt
==============

Patrick Wagstrom &lt;160672+pridkett@users.noreply.github.com&gt;<br>
June 2022

Overview
--------

This program connects to local web server on a PurpleAir air quality monitor
and publishes the data to an MQTT broker and optionally logs the data into an
influxdb database for further analysis. There are already a few libraries out there that connect the PurpleAir API, but I'm the kind of person that wants local interfaces to local devices

**US EPA AQI Calculation**: This program now calculates the US EPA Air Quality Index (AQI) based on PM2.5 and PM10 concentrations reported by the PurpleAir sensor, following the official EPA guidelines.

Background
----------

Once you have your PurpleAir monitor connected to your local network, you can access it by just going to its IP address and making an HTTP request. This provides a user friendly page of what is going on with the monitor. If you'd like to get the data structured, make a request to `/json`.

This provides an excellent way to get real-time information from a local device without needing to manage API keys or even make calls out to the public internet. It's this JSON payload that the program parses.

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

Configuration
-------------

This program uses TOML as it's configuration file because all configuraiton file formats are terrible and TOML it terrible in the least amount of conflicting ways. You'll need to figure out the IP address or hostname of your Purple Air monitor on your local network.

First, define the location of the device and the polling rate in seconds. I find that you don't need to read it much more than once every two minutes. This will allow the application to pull the information from the Purple Air sensor.

```toml
[purpleair]
    url = "http://192.168.1.24/json"
    poll_rate = 120

```

Next, define the information needed for wherever your MQTT broker is running. I haven't done anything fancy here, but if it helps, I use Mosquitto as my MQTT broker. You'll need the hostname. I don't use any username or password here. If `topic_prefix` is left as null it will default to `airquality` and if `topic` is left as null it will default to the `geo` identifier of your PurpleAir sensor.

```toml
[mqtt]
    broker_host = "mqttbroker.local"
    broker_port = 1883
    client_id = "purpleair2mqtt"
    topic_prefix = "airquality"
    topic = ""
```

If you want Home Assistant integration, this chunk _should_ work to support Home Assistant auto discovery from MQTT. I need to do some more testing to see how well it works. I don't actually use Home Assistant that much.

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

Building the Application
------------------------

This should be able to be built like most other straightforward golang applications.

```bash
go build
```

Running the Application
-----------------------

```bash
./purpleair2mqtt -config config.toml
```

Grafana Integration
-------------------

This application has some lightweight Grafana integration, but it's not what I'd call fancy. I'll document that more in the future.

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

Running with Docker
-------------------

One of the advantages of running this application with Docker is that you can use `docker-compose` and then not worry about needing to restart the system all the time. This makes this way less of a concern when problems happen.

### Building the Container

```bash
docker build -t pridkett/purpleair2mqtt .
```

### Running the Container as a One Off

### Running the Container from `docker-compose`

## Authors & License

Copyright (c) 2022 [Patrick Wagstrom](https://github.com/pridkett); modifications (c) 2025 [Chris Dzombak](https://github.com/cdzombak)

Licensed under the terms of the MIT License; see [LICENSE](LICENSE) in this repo.
