# Distributed Mandelbrot Generator

|![mandelbrot](saved_images/mandelbrot.png)|![seahorse](saved_images/seahorse.png)|
|:---:|:---:|:---:|
|![valley](saved_images/valley.png)|![crown](saved_images/crown.png)|

## Performance 

### 29/12/2022

CPU Intel Core i5-9600K (6 Cores) @ 3.70GHz, RAM 16.0 Go, computing Seahorse preset at 2048*2048px resolution, 8192 iterations.

No containers, no cpu usage limits (~100%), no chunks : 
  2m0.2565745s
  
Single container, no cpu usage limits (~100%), no chunks :
  2m7.7364957s
  
6 containers, each limited to 1.0 cpu usage, 64 chunks :
  longest chunk: 3m52.2680694s
  
6 containers, no cpu limits, 64 chunks :
  longest chunk: 2m10.5041407s
  
Note: When computing chunks, a single container is assigned several chunks requests concurently. A commercial NGINX subscription is required to limit containers to a single request with a queue.
See: http://nginx.org/en/docs/http/ngx_http_upstream_module.html#queue
