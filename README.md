# httpDump

A completely hacky mostly designed to play with via a debugger tool. There are no tests, it's designed to be changed when I need to. Dumps requests sent to it to the console

Allows testing various web hooks and easily viewing requests, can also dump the request in curl form

Useful with ngrok

Defaults to port 9000

## flags
* `--port` override `PORT` env var
* `--curl` really poor mans curl output, I forget why I did this
* `--response` fixed reply alternative
* `--no-colour` disable colour

```bash
httpDump --help
```
