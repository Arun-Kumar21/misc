# Concurrent URL Health Checker

Takes a string of URLs from the user and uses the `http.Get` method to check the status of each URL.

When we don't use concurrency, i.e. sequential URL fetching, it takes a very long time to fetch 100s of URLs. If some failed to fetch, `http.Get` tries to fetch it for a maximum of 10 times, which means more waiting :( .

`main.go` contains code for sequential fetching and shows the difference between concurrency and sequential.

For achieving concurrency, I created a `worker` function that takes `id`, `jobs`, `results`:
- `id` - It is basically to uniquely identify a worker
- `jobs` - A channel of string, and its direction is receiver (means can only get data from the channel)
- `results` - A channel of string, and its direction is sender (means can only put data in the channel)

I created 5 workers in the `main` function and initialized `jobs` and `results` channels with a capacity of 5. When running main with concurrency, 5 workers are created which are goroutines and they wait for an item to appear in jobs, then I pass a batch of 5 URLs into jobs which are picked up by workers, and after those 5 URLs are processed, each worker puts its result in the results channel.

Then the next batch of 5 URLs will be processed, and so on until all URLs are processed.

Finally, I print a comparison time for both fetching methods at the end of the main func.


Output :- 
```
❯ go run main.go
https://google.com https://arun.space https://x.com https://meta.com https://go.dev https://archlinux.org https://github.com https://en.wikipedia.org https://youtube.com https://cornhub.com 
Sequential fetching
200 OK
Failed to fetch
200 OK
200 OK
200 OK
200 OK
200 OK
403 Forbidden
200 OK
Failed to fetch
Failed to fetch
Concurrency fetching
Worker  0 started job: https://x.com
200 OK
Worker  2 started job: https://go.dev
200 OK
Worker  0 started job: https://archlinux.org
200 OK
Worker  2 started job: https://github.com
200 OK
Worker  0 started job: https://en.wikipedia.org
403 Forbidden
Worker  4 started job: https://google.com
Worker  4 started job: 
200 OK
 -> ERROR: Get "": unsupported protocol scheme ""
Worker  2 started job: https://youtube.com
200 OK
Worker  3 started job: https://meta.com
200 OK
Worker  1 started job: https://arun.space
https://arun.space -> ERROR: Get "https://arun.space": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
Worker  0 started job: https://cornhub.com
https://cornhub.com -> ERROR: Get "https://cornhub.com": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
Sequential: 19.430494503s | Concurrent: 5.986655459s
```

I just added this so that I don't have to type these 100 times when running my program, you can also use these :)

```bash
https://google.com https://arun.space https://x.com https://meta.com https://go.dev https://archlinux.org https://github.com https://en.wikipedia.org https://youtube.com https://cornhub.com 
```
