# restRoberto
[![Go Report Card](https://goreportcard.com/badge/github.com/TheTipo01/restRoberto)](https://goreportcard.com/report/github.com/TheTipo01/restRoberto)
[![Build Status](https://travis-ci.com/TheTipo01/restRoberto.svg?branch=master)](https://travis-ci.com/TheTipo01/restRoberto)

restRoberto - Simple HTTP API that generates audio file with the (not so) famous Roberto voice.

## Usage
Make a GET request to the `/audio` endpoint with the following paramaters: 
- `token` with one of the tokens you specified in the config.yml
- `text` with the text you want to be generated

The server will respond with code `202` and the path on the server where you can find the newly generated audio file, in the form of `/temp/UUID.mp3`.
