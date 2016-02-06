#Overview

Here we have two programs:
  - A back-end command line collector to retrieve your connection status;
  - A web server to display the data collected.

The data is stored in a SQLite database, but since everything is automated, you don't have to worry about that.

Unfortunately, the program uses a program available only to Linux: [speed-test](github.com/sivel/speedtest-cli).

# Usage

## CLI

The command line program should run periodically in order to collect data over time. I recommend creating an entry in cron run it every 15 or 30 minutes.

The only configuration required is the server code to where you want to connect to run the tests. That can be retrieved using the following commands:

```bash
$ speedtest-cli --list | grep Brussels
5151) Combell (Brussels, Belgium) [3.61 km]
```

In this example, the code is `5151`. So, run the collector:

```bash
$ check-my-speed-cli -server 5151
```

## Web Server

To see the data collected, start the web server with the following command:

```bash
$ speedtest-web
```

The charts are available at [http://localhost:8088/speed]() by default. If you wish to use a custom domain or port, use the parameter `url`:

```bash
$ speedtest-web -url domain_or_ip:port_number
```

# Reports

Since this is the very first version, the reports are very ugly, but functional. There are 4 gauge images, showing, from left to right, the lowest value collected, the highest one, the average value and the last value collected. The gauges assume a default nominal connection of 50 Mbit/s, where the red section and the green section meet.

There is also a line graph with the blue line representing the download speed over time and the red line, the upload speed. Ideally, both should have constant value matching your nominal connection, but that's hardly the case.

I'm using the Google Charts Javascript Lib to render the graphs and I plan to improve the look'n'feel of this report (perhapes, even create some others?), but I'm not really a designer, so any suggestion are welcome.
