# tepcollect

A tool to get electricity prices and usage from [くらしTEPCO web](https://www.app.kurashi.tepco.co.jp/).

## Usage
```bash
docker build -t tepcollect .
docker run -it -e TEPCO_USERNAME=***** -e TEPCO_PASSWORD=***** tepcollect
```