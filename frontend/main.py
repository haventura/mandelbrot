import streamlit as st
import requests
import json
import dataclasses
import base64
from PIL import Image
import io
import math
from aiohttp import ClientSession, ClientTimeout
import asyncio
import os
import time
import logging

@dataclasses.dataclass
class Image_data:
    width:         int   = 512
    height:        int   = 512
    max_iteration: int   = 256
    min_r:         float = -2.0
    max_r:         float = 0.5
    min_i:         float = -1.25
    max_i:         float = 1.25
    colormap_name: str   = "BWY"

# run with python -m streamlit run main.py

async def main():
    st.set_page_config(page_title="Mandelbrot Generator", page_icon="🚀", layout="wide")
    st.title("Mandelbrot Generator")

    col1, col2, col3 = st.columns([1,2,1])
    if 'center_r' not in st.session_state:
        st.session_state['center_r'] = -.75
        st.session_state["center_i"] = .0
        st.session_state["diameter"] = 2.5
        st.session_state["image_placeholder"] = col2.empty()

    with col1.form("my_form"):
        image_data = Image_data()
        st.header("🧪 Parameters")
        center_r = st.number_input("Center Re(c)", min_value=-3., max_value=3., step=0.1, format="%.15g", key="center_r")
        center_i = st.number_input("Center Im(c)", min_value=-3., max_value=3., step=0.1, format="%.15g", key="center_i")
        diameter = st.number_input("Diameter", min_value=0., max_value=6., step=0.1, format="%.15g", key="diameter")
        image_data.max_iteration = st.select_slider("Iterations", get_resolution_steps(1, 65536), value=256, key="max_iteration")
        resolution = st.select_slider("Resolution", get_resolution_steps(128, 16384), value=512, key="resolution")
        image_data.colormap_name = st.selectbox("Color map", options=["BWY","BWR","ICEFIRE","TWILIGHT"], key="colormap")
        chunks = st.checkbox('Compute in chunks')
        chunks_amount = st.select_slider("Number of chunks", get_chunks_steps(4096), value = 16)
        image_data.min_r = center_r - diameter / 2
        image_data.max_r = center_r + diameter / 2
        image_data.min_i = center_i - diameter / 2
        image_data.max_i = center_i + diameter / 2

        image_data.width = resolution
        image_data.height = resolution

        submitted = st.form_submit_button("Submit")

        if submitted:
            try:
                addr = os.getenv('BACKEND_ADDR')
                if not addr:
                    addr = "nginx:80"
            except KeyError:
                addr = "nginx:80"
            url = f'http://{addr}/compute'
            logging.info(image_data)
            if not chunks:
                response = requests.post(url, data = json.dumps(dataclasses.asdict(image_data)))
                response_as_base64 = response.content.decode("utf-8")
                image = base64.b64decode(response_as_base64.split(',')[1])
                st.session_state["image_placeholder"].image(image, use_column_width="always")
                col1.download_button("💾 Download", image, "mandelbrot.png", key="download_button")
            else:      
                chunks = divide_in_chunk(image_data, chunks_amount)
                image = Image.new("RGBA", (resolution, resolution))
                async with ClientSession(timeout=ClientTimeout(total=3600)) as session:
                    coros = [task_coro(session, chunk_position, chunk_data, url, image) for chunk_position, chunk_data in chunks.items()]
                    await asyncio.gather(*coros)
                img_byte_arr = io.BytesIO()
                image.save(img_byte_arr, format='PNG')
                img_byte_arr = img_byte_arr.getvalue()
                col1.download_button("💾 Download", img_byte_arr, "mandelbrot.png", key="download_button")

    with col3:
        st.header("Some interesting features")
        st.subheader("🏠 Home:")
        st.button("Try me!", key="home_button", on_click=update_coordinates, args=(-0.75, 0.0, 2.5))
        st.subheader("🌊🐎 Seahorse:")
        st.button("Try me!", key="seahorse_button", on_click=update_coordinates, args=(-0.74303, 0.126433, 0.01611))
        st.subheader("🦎 Tail:")
        st.button("Try me!", key="tail_button", on_click=update_coordinates, args=(-0.7436499, 0.13188204, 0.00073801))
        st.subheader("👑 Crown:")
        st.button("Try me!", key="crown_button", on_click=update_coordinates, args=(-0.743643135, 0.131825963, 0.000014628))
        st.subheader("🛰️ Satellite:")
        st.button("Try me!", key="satellite_button", on_click=update_coordinates, args=(-0.743644786, 0.1318252536, 0.0000029336))
        st.subheader("⛰️ Valley:")
        st.button("Try me!", key="valley_button", on_click=update_coordinates, args=(-0.74364386269, 0.13182590271, 0.00000013526))

    st.write('<style>div.block-container{padding-top:1rem;}</style>', unsafe_allow_html=True)
    hide_footer_style = """
            <style>
            #MainMenu {visibility: hidden;}
            footer {visibility: hidden;}
            </style>
            """
    st.markdown(hide_footer_style, unsafe_allow_html=True) 
    st.caption("Made with Go and Streamlit by Andrea Ventura. <a href=https://github.com/haventura/mandelbrot>GitHub project</a>", unsafe_allow_html=True)

def update_coordinates(center_r, center_i, diameter):
    st.session_state["center_r"] = center_r
    st.session_state["center_i"] = center_i
    st.session_state["diameter"] = diameter    

def divide_in_chunk(image_data: Image_data, n_chunk):
    n_horizontal_chunk = n_vertical_chunk = int(math.sqrt(n_chunk))
    chunk_pixel_width = image_data.width // n_horizontal_chunk
    chunk_pixel_height = image_data.height // n_vertical_chunk    
    chunk_real_width = (image_data.max_r - image_data.min_r) / n_horizontal_chunk
    chunk_imaginary_width = (image_data.max_i - image_data.min_i) / n_vertical_chunk
    dict_of_chunk = {}
    for i in range(n_vertical_chunk):
        for j in range(n_horizontal_chunk):
            chunk_min_x = chunk_pixel_width * j
            chunk_min_y = chunk_pixel_height * i
            image_chunk_data = Image_data()
            image_chunk_data.width = chunk_pixel_width
            image_chunk_data.height = chunk_pixel_height
            image_chunk_data.min_r = image_data.min_r + (chunk_real_width * j)
            image_chunk_data.max_r = image_data.min_r + (chunk_real_width * (j + 1))
            image_chunk_data.min_i = image_data.max_i - (chunk_imaginary_width * (i + 1))       # !!! inverted y axis for pixels !!!
            image_chunk_data.max_i = image_data.max_i - (chunk_imaginary_width * (i))
            image_chunk_data.max_iteration = image_data.max_iteration
            image_chunk_data.colormap_name = image_data.colormap_name
            dict_of_chunk[(chunk_min_x, chunk_min_y)] = image_chunk_data
    return dict_of_chunk

def get_resolution_steps(min, max):
    output = []
    a = min
    while True:
        if a <= max:
            output.append(a)
            a *= 2
            continue
        return output

def get_chunks_steps(max):
    output = []
    a = 4
    while True:
        if a <= max:
            output.append(a)
            a *= 4
            continue
        return output

async def task_coro(session, chunk_position, chunk_data, url, image):
    async with session.post(url, data = json.dumps(dataclasses.asdict(chunk_data))) as response:
        awaited_response = await response.read()
        response_as_base64 = awaited_response.decode("utf-8")
        bytes_response = base64.b64decode(response_as_base64.split(',')[1])
        sub_image = Image.open(io.BytesIO(bytes_response))
        image.paste(sub_image, (chunk_position[0], chunk_position[1]))
        sub_image.close()
        start = time.time()
        st.session_state["image_placeholder"].image(image, use_column_width="always")
        end = time.time()
        logging.info(end - start)

if __name__ == '__main__':
    asyncio.run(main())