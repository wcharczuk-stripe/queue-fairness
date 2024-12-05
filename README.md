queue fairness simulator
========================

To run, simply:
> go run main.go

Example output:

```$ go run main.go --queue-type=fairness                                                                                                                                 queue-fairness main
using task queue type:          fairness
using simulation duration:      1h0m0s
using compaction interval:      15m0s
using tick interval:            500ms

compaction (15m0s) ts=22:04 elapsed=15m0s tql=337786 tp=2904360
compaction (15m0s) ts=22:19 elapsed=30m0s tql=578225 tp=2913183
compaction (15m0s) ts=22:34 elapsed=45m0s tql=777228 tp=2913011
compaction (15m0s) ts=22:49 elapsed=1h0m0s tql=1057720 tp=2913001

tasks processed: 11643555
queued for      p95: 13m26.5s   avg: 2m35.469s

queued for by priority "P0" [12044]             p95: 0s avg: 0s
queued for by priority "P1" [118875]            p95: 0s avg: 0s
queued for by priority "P2" [11858446]          p95: 4m14.5s    avg: 1m1.91s
queued for by priority "P3" [118681]            p95: 55m50.5s   avg: 28m59.839s
queued for by priority "P4" [593229]            p95: 55m55s     avg: 29m3.015s
queued for by fairness key "high" [751856]      p95: 5m10.5s    avg: 1m36.508s
queued for by fairness key "low" [3629465]      p95: 13m48.5s   avg: 2m40.196s
queued for by fairness key "medium" [7262234]   p95: 13m39s     avg: 2m39.013s
```