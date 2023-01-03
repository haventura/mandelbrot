# Distributed Mandelbrot Generator

<img src="saved_images/mandelbrot.png" width="24.5%"></img>
<img src="saved_images/seahorse.png" width="24.5%"></img>
<img src="saved_images/valley.png" width="24.5%"></img>
<img src="saved_images/crown.png" width="24.5%"></img> 

## Architecture

### Streamlit Frontend (async requests for chunks, etc...)

blablabla

### Go backend (go routines, etc...)

blablabla

### Docker container & NGINX load balancing.

blablabla

## Deployment

### Locally

Follow this link to create an ACI context witin Azure: https://docs.docker.com/cloud/aci-integration/

### Remotely, with Azure and Docker ACI integration

Using default docker context:
docker build -t andreaventura/ecam:mandelbrot_nginx .
docker push -t andreaventura/ecam:mandelbrot_nginx .
repeat for all 3 images; nginx, backend, frontend 

Then in docker ACI context:
docker compose up

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
