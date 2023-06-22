window.BENCHMARK_DATA = {
  "lastUpdate": 1687449421927,
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
          "id": "e77ea6288871025431a0481798357004750ecc84",
          "message": "Update perf.yml",
          "timestamp": "2023-06-22T08:52:18-07:00",
          "tree_id": "9387593acffd1ce38200ec92249375af9140ccb8",
          "url": "https://github.com/codeboten/opentelemetry-collector/commit/e77ea6288871025431a0481798357004750ecc84"
        },
        "date": 1687449416443,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkBoundedQueue",
            "value": 563.5,
            "unit": "ns/op",
            "extra": "2222745 times\n2 procs"
          },
          {
            "name": "BenchmarkBoundedQueueWithFactory",
            "value": 521.1,
            "unit": "ns/op",
            "extra": "2908728 times\n2 procs"
          },
          {
            "name": "BenchmarkTraceSizeSpanCount",
            "value": 6.471,
            "unit": "ns/op",
            "extra": "185308569 times\n2 procs"
          },
          {
            "name": "BenchmarkBatchMetricProcessor",
            "value": 1516728,
            "unit": "ns/op",
            "extra": "1028 times\n2 procs"
          },
          {
            "name": "BenchmarkMultiBatchMetricProcessor",
            "value": 1102591,
            "unit": "ns/op",
            "extra": "1086 times\n2 procs"
          },
          {
            "name": "BenchmarkHttpRequest/HTTP/2.0,_shared_client_(like_load_balancer)",
            "value": 84648,
            "unit": "ns/op",
            "extra": "13563 times\n2 procs"
          },
          {
            "name": "BenchmarkHttpRequest/HTTP/1.1,_shared_client_(like_load_balancer)",
            "value": 52257,
            "unit": "ns/op",
            "extra": "27291 times\n2 procs"
          },
          {
            "name": "BenchmarkHttpRequest/HTTP/2.0,_client_per_thread_(like_single_app)",
            "value": 84892,
            "unit": "ns/op",
            "extra": "13881 times\n2 procs"
          },
          {
            "name": "BenchmarkHttpRequest/HTTP/1.1,_client_per_thread_(like_single_app)",
            "value": 44254,
            "unit": "ns/op",
            "extra": "23773 times\n2 procs"
          },
          {
            "name": "BenchmarkCompressors/sm_log_request/raw_bytes_160/compressed_bytes_162/compressor_gzip",
            "value": 41968,
            "unit": "ns/op",
            "extra": "28315 times\n2 procs"
          },
          {
            "name": "BenchmarkJSONUnmarshal",
            "value": 2704,
            "unit": "ns/op\t    1040 B/op\t      55 allocs/op",
            "extra": "443095 times\n2 procs"
          },
          {
            "name": "BenchmarkLogsToProto",
            "value": 9514,
            "unit": "ns/op",
            "extra": "127680 times\n2 procs"
          },
          {
            "name": "BenchmarkLogsFromProto",
            "value": 20711,
            "unit": "ns/op\t   16624 B/op\t     141 allocs/op",
            "extra": "56589 times\n2 procs"
          },
          {
            "name": "BenchmarkOtlpToFromInternal_PassThrough",
            "value": 0.4187,
            "unit": "ns/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "BenchmarkOtlpToFromInternal_Gauge_MutateOneLabel",
            "value": 54.35,
            "unit": "ns/op",
            "extra": "21861579 times\n2 procs"
          },
          {
            "name": "BenchmarkOtlpToFromInternal_Sum_MutateOneLabel",
            "value": 54.28,
            "unit": "ns/op",
            "extra": "21401434 times\n2 procs"
          },
          {
            "name": "BenchmarkOtlpToFromInternal_HistogramPoints_MutateOneLabel",
            "value": 54.13,
            "unit": "ns/op",
            "extra": "21460725 times\n2 procs"
          },
          {
            "name": "BenchmarkMetricsToProto",
            "value": 18739,
            "unit": "ns/op",
            "extra": "70240 times\n2 procs"
          },
          {
            "name": "BenchmarkMetricsFromProto",
            "value": 62023,
            "unit": "ns/op\t   31984 B/op\t     909 allocs/op",
            "extra": "19326 times\n2 procs"
          },
          {
            "name": "BenchmarkJSONUnmarshal",
            "value": 10371,
            "unit": "ns/op\t    4392 B/op\t     217 allocs/op",
            "extra": "118621 times\n2 procs"
          },
          {
            "name": "BenchmarkTracesToProto",
            "value": 12745,
            "unit": "ns/op",
            "extra": "90326 times\n2 procs"
          },
          {
            "name": "BenchmarkTracesFromProto",
            "value": 34759,
            "unit": "ns/op\t   30960 B/op\t     269 allocs/op",
            "extra": "34370 times\n2 procs"
          }
        ]
      }
    ]
  }
}