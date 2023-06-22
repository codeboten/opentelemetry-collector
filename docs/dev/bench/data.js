window.BENCHMARK_DATA = {
  "lastUpdate": 1687467673285,
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
      },
      {
        "commit": {
          "author": {
            "email": "alex@boten.ca",
            "name": "Alex Boten",
            "username": "codeboten"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "c6b49f8b1b0b25495e09695eac4eaee0f479189b",
          "message": "Merge branch 'open-telemetry:main' into main",
          "timestamp": "2023-06-22T13:57:15-07:00",
          "tree_id": "0262bbded4b63b4715d982cb2eb99e0ce6d40b94",
          "url": "https://github.com/codeboten/opentelemetry-collector/commit/c6b49f8b1b0b25495e09695eac4eaee0f479189b"
        },
        "date": 1687467663842,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkBoundedQueue",
            "value": 477.3,
            "unit": "ns/op",
            "extra": "2670004 times\n2 procs"
          },
          {
            "name": "BenchmarkBoundedQueueWithFactory",
            "value": 255.6,
            "unit": "ns/op",
            "extra": "4180713 times\n2 procs"
          },
          {
            "name": "BenchmarkTraceSizeSpanCount",
            "value": 5.162,
            "unit": "ns/op",
            "extra": "228021067 times\n2 procs"
          },
          {
            "name": "BenchmarkBatchMetricProcessor",
            "value": 913532,
            "unit": "ns/op",
            "extra": "1239 times\n2 procs"
          },
          {
            "name": "BenchmarkMultiBatchMetricProcessor",
            "value": 929682,
            "unit": "ns/op",
            "extra": "1298 times\n2 procs"
          },
          {
            "name": "BenchmarkHttpRequest/HTTP/2.0,_shared_client_(like_load_balancer)",
            "value": 70622,
            "unit": "ns/op",
            "extra": "16936 times\n2 procs"
          },
          {
            "name": "BenchmarkHttpRequest/HTTP/1.1,_shared_client_(like_load_balancer)",
            "value": 45233,
            "unit": "ns/op",
            "extra": "31876 times\n2 procs"
          },
          {
            "name": "BenchmarkHttpRequest/HTTP/2.0,_client_per_thread_(like_single_app)",
            "value": 69251,
            "unit": "ns/op",
            "extra": "17020 times\n2 procs"
          },
          {
            "name": "BenchmarkHttpRequest/HTTP/1.1,_client_per_thread_(like_single_app)",
            "value": 39802,
            "unit": "ns/op",
            "extra": "27282 times\n2 procs"
          },
          {
            "name": "BenchmarkCompressors/sm_log_request/raw_bytes_160/compressed_bytes_162/compressor_gzip",
            "value": 35575,
            "unit": "ns/op",
            "extra": "33264 times\n2 procs"
          },
          {
            "name": "BenchmarkJSONUnmarshal",
            "value": 2227,
            "unit": "ns/op\t    1040 B/op\t      55 allocs/op",
            "extra": "462446 times\n2 procs"
          },
          {
            "name": "BenchmarkLogsToProto",
            "value": 7477,
            "unit": "ns/op",
            "extra": "159126 times\n2 procs"
          },
          {
            "name": "BenchmarkLogsFromProto",
            "value": 16690,
            "unit": "ns/op\t   16624 B/op\t     141 allocs/op",
            "extra": "71504 times\n2 procs"
          },
          {
            "name": "BenchmarkOtlpToFromInternal_PassThrough",
            "value": 0.336,
            "unit": "ns/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "BenchmarkOtlpToFromInternal_Gauge_MutateOneLabel",
            "value": 43.54,
            "unit": "ns/op",
            "extra": "27416848 times\n2 procs"
          },
          {
            "name": "BenchmarkOtlpToFromInternal_Sum_MutateOneLabel",
            "value": 43.51,
            "unit": "ns/op",
            "extra": "26043350 times\n2 procs"
          },
          {
            "name": "BenchmarkOtlpToFromInternal_HistogramPoints_MutateOneLabel",
            "value": 46.7,
            "unit": "ns/op",
            "extra": "26912157 times\n2 procs"
          },
          {
            "name": "BenchmarkMetricsToProto",
            "value": 13650,
            "unit": "ns/op",
            "extra": "88951 times\n2 procs"
          },
          {
            "name": "BenchmarkMetricsFromProto",
            "value": 50316,
            "unit": "ns/op\t   31984 B/op\t     909 allocs/op",
            "extra": "24180 times\n2 procs"
          },
          {
            "name": "BenchmarkJSONUnmarshal",
            "value": 8434,
            "unit": "ns/op\t    4392 B/op\t     217 allocs/op",
            "extra": "142906 times\n2 procs"
          },
          {
            "name": "BenchmarkTracesToProto",
            "value": 10088,
            "unit": "ns/op",
            "extra": "118888 times\n2 procs"
          },
          {
            "name": "BenchmarkTracesFromProto",
            "value": 28159,
            "unit": "ns/op\t   30960 B/op\t     269 allocs/op",
            "extra": "42454 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "alex@boten.ca",
            "name": "Alex Boten",
            "username": "codeboten"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "c6b49f8b1b0b25495e09695eac4eaee0f479189b",
          "message": "Merge branch 'open-telemetry:main' into main",
          "timestamp": "2023-06-22T13:57:15-07:00",
          "tree_id": "0262bbded4b63b4715d982cb2eb99e0ce6d40b94",
          "url": "https://github.com/codeboten/opentelemetry-collector/commit/c6b49f8b1b0b25495e09695eac4eaee0f479189b"
        },
        "date": 1687467668343,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkBoundedQueue",
            "value": 359.3,
            "unit": "ns/op",
            "extra": "8268156 times\n2 procs"
          },
          {
            "name": "BenchmarkBoundedQueueWithFactory",
            "value": 450.2,
            "unit": "ns/op",
            "extra": "2659314 times\n2 procs"
          },
          {
            "name": "BenchmarkTraceSizeSpanCount",
            "value": 5.363,
            "unit": "ns/op",
            "extra": "223597473 times\n2 procs"
          },
          {
            "name": "BenchmarkBatchMetricProcessor",
            "value": 903848,
            "unit": "ns/op",
            "extra": "1231 times\n2 procs"
          },
          {
            "name": "BenchmarkMultiBatchMetricProcessor",
            "value": 919020,
            "unit": "ns/op",
            "extra": "1222 times\n2 procs"
          },
          {
            "name": "BenchmarkHttpRequest/HTTP/2.0,_shared_client_(like_load_balancer)",
            "value": 67706,
            "unit": "ns/op",
            "extra": "16111 times\n2 procs"
          },
          {
            "name": "BenchmarkHttpRequest/HTTP/1.1,_shared_client_(like_load_balancer)",
            "value": 39944,
            "unit": "ns/op",
            "extra": "29788 times\n2 procs"
          },
          {
            "name": "BenchmarkHttpRequest/HTTP/2.0,_client_per_thread_(like_single_app)",
            "value": 70522,
            "unit": "ns/op",
            "extra": "17066 times\n2 procs"
          },
          {
            "name": "BenchmarkHttpRequest/HTTP/1.1,_client_per_thread_(like_single_app)",
            "value": 40068,
            "unit": "ns/op",
            "extra": "28660 times\n2 procs"
          },
          {
            "name": "BenchmarkCompressors/sm_log_request/raw_bytes_160/compressed_bytes_162/compressor_gzip",
            "value": 34783,
            "unit": "ns/op",
            "extra": "34128 times\n2 procs"
          },
          {
            "name": "BenchmarkJSONUnmarshal",
            "value": 2193,
            "unit": "ns/op\t    1040 B/op\t      55 allocs/op",
            "extra": "552008 times\n2 procs"
          },
          {
            "name": "BenchmarkLogsToProto",
            "value": 8259,
            "unit": "ns/op",
            "extra": "154077 times\n2 procs"
          },
          {
            "name": "BenchmarkLogsFromProto",
            "value": 16688,
            "unit": "ns/op\t   16624 B/op\t     141 allocs/op",
            "extra": "73981 times\n2 procs"
          },
          {
            "name": "BenchmarkOtlpToFromInternal_PassThrough",
            "value": 0.3374,
            "unit": "ns/op",
            "extra": "1000000000 times\n2 procs"
          },
          {
            "name": "BenchmarkOtlpToFromInternal_Gauge_MutateOneLabel",
            "value": 43.56,
            "unit": "ns/op",
            "extra": "27166670 times\n2 procs"
          },
          {
            "name": "BenchmarkOtlpToFromInternal_Sum_MutateOneLabel",
            "value": 42.59,
            "unit": "ns/op",
            "extra": "27127243 times\n2 procs"
          },
          {
            "name": "BenchmarkOtlpToFromInternal_HistogramPoints_MutateOneLabel",
            "value": 42.53,
            "unit": "ns/op",
            "extra": "26995425 times\n2 procs"
          },
          {
            "name": "BenchmarkMetricsToProto",
            "value": 13305,
            "unit": "ns/op",
            "extra": "89409 times\n2 procs"
          },
          {
            "name": "BenchmarkMetricsFromProto",
            "value": 47502,
            "unit": "ns/op\t   31984 B/op\t     909 allocs/op",
            "extra": "24994 times\n2 procs"
          },
          {
            "name": "BenchmarkJSONUnmarshal",
            "value": 9273,
            "unit": "ns/op\t    4392 B/op\t     217 allocs/op",
            "extra": "144297 times\n2 procs"
          },
          {
            "name": "BenchmarkTracesToProto",
            "value": 10382,
            "unit": "ns/op",
            "extra": "117436 times\n2 procs"
          },
          {
            "name": "BenchmarkTracesFromProto",
            "value": 26946,
            "unit": "ns/op\t   30960 B/op\t     269 allocs/op",
            "extra": "44330 times\n2 procs"
          }
        ]
      }
    ]
  }
}