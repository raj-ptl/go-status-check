## Go Status Check

Endpoints : 

> post URLs to track
>
`POST` /websites 

Sample request body :

```
{
  "websites" : [
    "google.com",
    "twitter.com",
    "localhost:9090/",
    "localhost:9090/undefinedEndpoint"
    ]
}
```
---
> get status of added URLs
>
`GET` /websites 

Sample response

```
{
  "StatusArray": [
    {
      "URL": "twitter.com",
      "Status": "UP",
      "LastChecked": "2022-08-24T20:15:02.3489142+05:30"
    },
    {
      "URL": "localhost:9090/undefinedEndpoint",
      "Status": "DOWN",
      "LastChecked": "2022-08-24T20:15:01.8757576+05:30"
    },
    {
      "URL": "localhost:9090/",
      "Status": "UP",
      "LastChecked": "2022-08-24T20:15:01.8755563+05:30"
    },
    {
      "URL": "google.com",
      "Status": "UP",
      "LastChecked": "2022-08-24T20:15:02.0752074+05:30"
    }
  ]
}
```
---
#### Polling Rate

Polling rate can be modified in `status/status.go`