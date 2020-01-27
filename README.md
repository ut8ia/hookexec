# HTTP hooks executor

HTTP Server that catch hooks from services and execute pre-configured scripts.
- POST or GET hooks
- Custom param, header and Auth tokens
- Python, bash or whatever you need as execution worker
- goLang powered : cross platform and minimal resources 

Typical use cases : 
- catch GitLab or GitHub hook to refresh environment
- build automation in CI\CD
- relay events to N messengers
 
### config example
```yaml
server:
  host: "0.0.0.0" # default address
  port: 8000 # port that listen hooks
  bodyLimit: 4096 # body limits for post hooks

request:
  header: "X-Gitlab-token" # Auth header name
  token: "llwixfry82347r6bx23874bvr6238x2423kk" # Auth header token
  param: "hook" # URL param that means a hook name, like ?hook=yourhookname

# list of your hooks 
hooks:
  'default':  # default hook name
    executor: "/bin/python" # system executor of the hook
    script: "./scripts/python_example.py" # script that will be executed
  'example':
    executor: "/bin/sh"
    script: "./scripts/bash_example.sh"
```

### hooks test examples
```bash
$ curl -X GET http://localhost:8000  -H "X-Gitlab-token: llwixfry82347r6bx23874bvr6238x2423kk"
```
Result : default hook script executed

```bash
$ curl -X GET http://localhost:8000?hook=example  -H "X-Gitlab-token: llwixfry82347r6bx23874bvr6238x2423kk"
```
Result : example hook script executed ( bash_example.sh - according to sample config )

```bash
$ curl -X GET http://localhost:8000?hook=unknownhookname  -H "X-Gitlab-token: llwixfry82347r6bx23874bvr6238x2423kk"
```
Result : default hook script executed ( hook name not resolved )

```bash
$ curl -X GET http://localhost:8000?hook=unknownhookname
```
Result : HTTP 401 Unauthorized ( Auth header missing or invalid )

```bash
curl -X POST http://localhost:8000?hook=example -d '{"mydata":1234234}' -H "X-Gitlab-token: llwixfry82347r6bx23874bvr6238x2423kk" 
```
Result : Post data successfully relayed to execution script as argument