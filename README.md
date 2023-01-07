# Distributed Mandelbrot Generator
 
This application is a demonstrator used to showcase the capability of Docker and NGINX to distribute workloads on several separate servers, by computing the [Mandelbrot set](https://en.wikipedia.org/wiki/Mandelbrot_set).

![Streamlit frontend](/saved_images/frontend.png)
<img src="saved_images/mandelbrot.png" width="24.6%"></img>
<img src="saved_images/seahorse.png" width="24.6%"></img>
<img src="saved_images/valley.png" width="24.6%"></img>
<img src="saved_images/crown.png" width="24.6%"></img> 

## Installation

### With Docker

1. From the root folder, run `docker compose up`. This command will start building the images for the frontend, load balancer and backend services. After building them, it will also create the related containers and start them.
1.  Once started, the streamlit application will be available at the following address: http://localhost:8501.

## Usage

## Architecture

### Frontend

The frontend application was develloped in python using [Streamlit](https://streamlit.io/). It uses a regular form template to upload the parameter set by the user to the backend which takes care of computing whether a point is within the mandelbrot set, and renders the image returned. A checkbox within that form enable the image to be computed as chunks, spreading the load of the computation over all server available to the application.

### Go backend (go routines, etc...)

Note that the Go program uses Goroutines to spread the computation on all available cores of the host machine; it is therefore unnecessary to have several instances of the Go service running on the same machine.

### Docker container & NGINX load balancing & Azure Container Instances.

blablabla

## References

https://docs.docker.com/cloud/aci-integration/