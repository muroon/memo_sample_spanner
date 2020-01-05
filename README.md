# memo sample spanner

The sample code of Cloud Spanner using "The Clean Architecture".
This is updated for Cloud Spanner based on [memo_sample](https://github.com/muroon/memo_sample).

[description of this project(Japanese)](https://gist.github.com/muroon/7daf23236777991a058544bd01ab9cc0)

### environment variable

| env | content |
----|----
| SPN_PROJECT_ID | GCP project ID |
| SPN_INSTANCE_ID | instance ID of Cloud Spanner |
| SPN_DATABASE_ID | database in instance  |

### use spanner emulator

You can use local [spanner emulator](https://github.com/gcpug/handy-spanner).
So you must set an environment parameter SPANNER_EMULATOR_HOST.

```
go build
env SPANNER_EMULATOR_HOST=localhost:9999 ./memo_sample_spanner -local
```

### static analysis tool

static analysis tool dedicated to this project

https://github.com/muroon/memo_sample_spanner_analyzer

