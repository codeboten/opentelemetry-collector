window.BENCHMARK_DATA = {
  "lastUpdate": 1687448562847,
  "repoUrl": "https://github.com/codeboten/opentelemetry-collector",
  "entries": {
    "Benchmark": [
      {
        "commit": {
          "author": {
            "email": "aboten@lightstep.com",
            "name": "Alex Boten",
            "username": "codeboten"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "235f62537f02965fbed10d542f9f675a229907a4",
          "message": "Update perf.yml",
          "timestamp": "2023-06-22T08:38:37-07:00",
          "tree_id": "ccd2dc7f18672d7f5648c7fa790a3cf8874b9fcf",
          "url": "https://github.com/codeboten/opentelemetry-collector/commit/235f62537f02965fbed10d542f9f675a229907a4"
        },
        "date": 1687448558839,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkBoundedQueue",
            "value": 452.7,
            "unit": "ns/op",
            "extra": "2680562 times\n2 procs"
          },
          {
            "name": "BenchmarkBoundedQueueWithFactory",
            "value": 486,
            "unit": "ns/op",
            "extra": "2713843 times\n2 procs"
          },
          {
            "name": "BenchmarkTraceSizeSpanCount",
            "value": 5.173,
            "unit": "ns/op",
            "extra": "232239184 times\n2 procs"
          },
          {
            "name": "BenchmarkBatchMetricProcessor",
            "value": 917943,
            "unit": "ns/op",
            "extra": "1332 times\n2 procs"
          },
          {
            "name": "BenchmarkMultiBatchMetricProcessor",
            "value": 932576,
            "unit": "ns/op",
            "extra": "1252 times\n2 procs"
          },
          {
            "name": "BenchmarkHttpRequest/HTTP/2.0,_shared_client_(like_load_balancer)",
            "value": 68187,
            "unit": "ns/op",
            "extra": "16945 times\n2 procs"
          },
          {
            "name": "BenchmarkHttpRequest/HTTP/1.1,_shared_client_(like_load_balancer)",
            "value": 36871,
            "unit": "ns/op",
            "extra": "30718 times\n2 procs"
          },
          {
            "name": "BenchmarkHttpRequest/HTTP/2.0,_client_per_thread_(like_single_app)",
            "value": 68442,
            "unit": "ns/op",
            "extra": "17293 times\n2 procs"
          },
          {
            "name": "BenchmarkHttpRequest/HTTP/1.1,_client_per_thread_(like_single_app)",
            "value": 37628,
            "unit": "ns/op",
            "extra": "26876 times\n2 procs"
          },
          {
            "name": "BenchmarkCompressors/sm_log_request/raw_bytes_160/compressed_bytes_162/compressor_gzip",
            "value": 35660,
            "unit": "ns/op",
            "extra": "33445 times\n2 procs"
          },
          {
            "name": "BenchmarkJSONUnmarshal",
            "value": 3015,
            "unit": "ns/op\t    1040 B/op\t      55 allocs/op",
            "extra": "541371 times\n2 procs"
          },
          {
            "name": "BenchmarkLogsToProto",
            "value": 7473,
            "unit": "ns/op",
            "extra": "155601 times\n2 procs"
          },
          {
            "name": "BenchmarkLogsFromProto",
            "value": 16695,
            "unit": "ns/op\t   16624 B/op\t     141 allocs/op",
            "extra": "72013 times\n2 procs"
          },
          {
            "name": "BenchmarkOtlpToFromInternal_PassThrough",
            "value": 0.3356,
            "unit": "ns/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "BenchmarkOtlpToFromInternal_Gauge_MutateOneLabel",
            "value": 43.48,
            "unit": "ns/op",
            "extra": "26957906 times\n2 procs"
          },
          {
            "name": "BenchmarkOtlpToFromInternal_Sum_MutateOneLabel",
            "value": 42.98,
            "unit": "ns/op",
            "extra": "27015198 times\n2 procs"
          },
          {
            "name": "BenchmarkOtlpToFromInternal_HistogramPoints_MutateOneLabel",
            "value": 42.9,
            "unit": "ns/op",
            "extra": "27103918 times\n2 procs"
          },
          {
            "name": "BenchmarkMetricsToProto",
            "value": 13393,
            "unit": "ns/op",
            "extra": "89181 times\n2 procs"
          },
          {
            "name": "BenchmarkMetricsFromProto",
            "value": 56600,
            "unit": "ns/op\t   31984 B/op\t     909 allocs/op",
            "extra": "24208 times\n2 procs"
          },
          {
            "name": "BenchmarkJSONUnmarshal",
            "value": 8478,
            "unit": "ns/op\t    4392 B/op\t     217 allocs/op",
            "extra": "146142 times\n2 procs"
          },
          {
            "name": "BenchmarkTracesToProto",
            "value": 10061,
            "unit": "ns/op",
            "extra": "117222 times\n2 procs"
          },
          {
            "name": "BenchmarkTracesFromProto",
            "value": 27960,
            "unit": "ns/op\t   30960 B/op\t     269 allocs/op",
            "extra": "42799 times\n2 procs"
          }
        ]
      }
    ]
  }
}