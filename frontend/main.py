import streamlit as st
import requests
import json
import dataclasses
import base64
from PIL import Image
import io
import math

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

@dataclasses.dataclass
class Image_chunk_data:
	total_width:   int  = 512 
	total_height:  int  = 512
	chunk_width:  int   = 512
	chunk_height: int   = 512
	chunk_min_x:  int   = 0
	chunk_min_y:  int   = 0
	chunk_min_r:  float = -2.0
	chunk_max_r:  float = 0.5
	chunk_min_i:  float = -1.25
	chunk_max_i:  float = 1.25
	max_iteration: int  = 256
	colormap_name: str  = "BWY"

# run with python -m streamlit run main.py

def main():
    st.set_page_config(page_title="Mandelbrot Generator", page_icon="🚀", layout="wide")
    st.title("Mandelbrot Generator")

    col1, col2, col3 = st.columns([1,2,1])
    image_placeholder = col2.empty()
    if 'center_r' not in st.session_state:
        st.session_state['center_r'] = -.75
        st.session_state["center_i"] = .0
        st.session_state["diameter"] = 2.5

    with col1.form("my_form"):
        image_data = Image_data()
        st.header("🧪 Parameters")
        center_r = st.number_input("Center Re(c)", min_value=-3., max_value=3., step=0.1, format="%.15g", key="center_r")
        center_i = st.number_input("Center Im(c)", min_value=-3., max_value=3., step=0.1, format="%.15g", key="center_i")
        diameter = st.number_input("Diameter", min_value=0., max_value=6., step=0.1, format="%.15g", key="diameter")
        image_data.max_iteration = st.select_slider("Iterations", get_resolution_steps(1, 65536), value=256, key="max_iteration")
        resolution = st.select_slider("Resolution", get_resolution_steps(128, 16384), value=512, key="resolution")
        image_data.colormap_name = st.text_input("Color map", value="BWY", key="colormap")
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
            if not chunks:
                url = 'http://localhost:5000/compute/single'
                print(image_data)
                response = requests.post(url, data = json.dumps(dataclasses.asdict(image_data)))
                response_as_base64 = response.content.decode("utf-8")
                image = base64.b64decode(response_as_base64.split(',')[1])
                col2.image(image, use_column_width="always")
                col1.download_button("💾 Download", image, "mandelbrot.png", key="download_button")
            else:
                url = 'http://localhost:5000/compute/chunk'
                print(image_data)
                chunks = divide_in_chunk(image_data, chunks_amount)
                image = Image.new("RGBA", (resolution, resolution))
                for chunk_data in chunks:
                    print(chunk_data)
                    response = requests.post(url, data = json.dumps(dataclasses.asdict(chunk_data)))
                    response_as_base64 = response.content.decode("utf-8")
                    bytes_response = base64.b64decode(response_as_base64.split(',')[1])
                    sub_image = Image.open(io.BytesIO(bytes_response))
                    image.paste(sub_image, (chunk_data.chunk_min_x, chunk_data.chunk_min_y))
                    sub_image.close()
                    image_placeholder.image(image, use_column_width="always")

                img_byte_arr = io.BytesIO()
                image.save(img_byte_arr, format='PNG')
                img_byte_arr = img_byte_arr.getvalue()
                col1.download_button("💾 Download", img_byte_arr, "mandelbrot.png", key="download_button")

    with col3:
        st.header("Some interesting features")
        st.subheader("🌊🐎 Seahorse:")
        st.button("Try me!", key="seahorse_button", on_click=update_coordinates, args=(-0.74303, 0.126433, 0.01611))
        st.subheader("🦎 Tail:")
        st.button("Try me!", key="tail_button", on_click=update_coordinates, args=(-0.7436499, 0.13188204, 0.00073801))
        st.subheader("👑 Crowns:")
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
    st.caption("Made with Go and Streamlit by Andrea Ventura")

def update_coordinates(center_r, center_i, diameter):
    st.session_state["center_r"] = center_r
    st.session_state["center_i"] = center_i
    st.session_state["diameter"] = diameter    

def divide_in_chunk(image_data: Image_data, n_chunk) -> list[Image_chunk_data]:
    n_horizontal_chunk = int(math.sqrt(n_chunk))
    n_vertical_chunk = int(math.sqrt(n_chunk))
    list_of_chunk = []
    for i in range(n_vertical_chunk):
        for j in range(n_horizontal_chunk):
            image_chunk_data = Image_chunk_data()
            image_chunk_data.total_width = image_data.width
            image_chunk_data.total_height = image_data.height
            image_chunk_data.max_iteration = image_data.max_iteration
            image_chunk_data.colormap_name = image_data.colormap_name

            image_chunk_data.chunk_width = image_data.width // n_horizontal_chunk
            image_chunk_data.chunk_height = image_data.height // n_vertical_chunk
            image_chunk_data.chunk_min_x = image_chunk_data.chunk_width * j
            image_chunk_data.chunk_min_y = image_chunk_data.chunk_height * i

            r_width = (image_data.max_r - image_data.min_r) / n_horizontal_chunk
            i_width = (image_data.max_i - image_data.min_i) / n_vertical_chunk

            image_chunk_data.chunk_min_r = image_data.min_r + (r_width * j)
            image_chunk_data.chunk_max_r = image_data.min_r + (r_width * (j + 1))
            image_chunk_data.chunk_max_i = image_data.min_i + (i_width * i)       # !!! inverted y axis for pixels !!!
            image_chunk_data.chunk_min_i = image_data.min_i + (i_width * (i + 1))

            list_of_chunk.append(image_chunk_data)

    return list_of_chunk

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



if __name__ == '__main__':
    main()