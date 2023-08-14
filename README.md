# restRoberto
[![Go Report Card](https://goreportcard.com/badge/github.com/TheTipo01/restRoberto)](https://goreportcard.com/report/github.com/TheTipo01/restRoberto)

restRoberto - Simple HTTP API that generates audio file with the (not so) famous Roberto voice.

## Endpoints

### GET `/audio`

Generates audio from the provided text and replies with status code 202 and the audio path.

Query parameters:

- `token`: the authorization token
- `text`: the text used to generate the audio

Example query: `GET https://rest.roberto.site/audio?token=valid_token&text=nyanpasu`

### GET `/temp/:uuid.mp3`

Get previously generated audio.
This is not meant to be used manually because the enpoint at `/audio` replies automatically with the full audio path.

This endpoint does not require authentication.
