# restRoberto
[![Go Report Card](https://goreportcard.com/badge/github.com/TheTipo01/restRoberto)](https://goreportcard.com/report/github.com/TheTipo01/restRoberto)

restRoberto - Simple HTTP API that generates audio file with the (not so) famous Roberto voice.

## Endpoints

### GET `/audio`

Generates audio from the provided text and replies with status code 202 and the audio as raw PCM.

Query parameters:

- `token`: the authorization token
- `text`: the text used to generate the audio
- `voice`: the voice to use (Roberto or Paola). Default is Roberto.

Example query: `GET https://rest.roberto.site/audio?token=valid_token&text=nyanpasu`

# Dockerfile
There's now a dockerfile that helps you to deploy restRoberto easily using Docker under wine.
You need to provide
- `roberto_setup.exe`, the Loquendo Roberto voice installer, in the same directory as the Dockerfile
- `LoqTTS6.dll`, the patched DLL, as it's practically impossible to activate the voice headlessly
