import streamlit as st
import requests
import json
import dataclasses
import base64

@dataclasses.dataclass
class Data:
    width:         int   = 1024
    height:        int   = 1024
    max_iteration: int   = 128
    min_r:         float = -2
    max_r:         float = 2
    min_i:         float = -2
    max_i:         float = 2
    colormap_name: str   ="mycolormap"

# run with python -m streamlit run main.py

def main():
    data = Data()
    st.set_page_config(page_title="Mandelbrot renderer", page_icon="random", layout="wide")
    st.title("Mandelbrot Generator")
    col1, col2 = st.columns([1,3])
    with col1.form("my_form"):
        st.write(f'<h3>Parameters</h3>', unsafe_allow_html=True)
        # data.min_r = st.number_input("Minimum real value", value=-2.23845, min_value=-2.23845, max_value=0.83845, step=0.1, format="%.10f", key="min_r")
        # data.max_r = st.number_input("Maximum real value", value=0.83845, min_value=-2.23845, max_value=0.83845, step=0.1, format="%.10f", key="max_r")
        # data.min_i = st.number_input("Minimum imaginary value", value=-1.53845, min_value=-1.53845, max_value=1.53845, step=0.1, format="%.10f", key="min_i")
        # data.max_i = st.number_input("Maximum imaginary value", value=1.53845, min_value=-1.53845, max_value=1.53845, step=0.1, format="%.10f", key="max_i")

        center_r = st.number_input("Center Re(c)", value=-.7, min_value=-3., max_value=3., step=0.1, format="%.15g", key="center_r")
        center_i = st.number_input("Center Im(c)", value=.0, min_value=-3., max_value=3., step=0.1, format="%.15g", key="center_i")
        diameter = st.number_input("Diameter", value=3., min_value=0., max_value=6., step=0.1, format="%.15g", key="diameter")
        data.max_iteration = st.number_input("Iterations", value=512, min_value=1, key="max_iteration")
        resolution = st.number_input("Resolution", value=1024, min_value=256, key="resolution")
        data.colormap_name = st.text_input("Color map", value="mycolormap", key="colormap")

        data.min_r = center_r - diameter / 2
        data.max_r = center_r + diameter / 2
        data.min_i = center_i - diameter / 2
        data.max_i = center_i + diameter / 2

        data.width = resolution
        data.height = resolution

        submitted = st.form_submit_button("Submit")

        if submitted:
            url = 'http://localhost:5000/compute'
            print(data)
            x = requests.post(url, data = json.dumps(dataclasses.asdict(data)))
            image_data = x.content.decode("utf-8")
            html_img = f'<img src="{image_data}" width="900" height="900"/>'
            col2.write(html_img, unsafe_allow_html=True)
            file_content=base64.b64decode(image_data.split(',')[1])
            col2.download_button("💾 Download", file_content, "mandelbrot.png")

    col1.write(f'<h3>Some interesting coordinates</h3>', unsafe_allow_html=True)
    col1.write(f'<h4>🌊🐎 Seahorse:</h4>', unsafe_allow_html=True)
    col1.button("Try me!", key="seahorse_button", on_click=update_coordinates, args=(-0.74303, 0.126433, 0.01611))
    col1.write(f'<h4>👑 Crowns:</h4>', unsafe_allow_html=True)
    col1.button("Try me!", key="crown_button", on_click=update_coordinates, args=(-0.743643135, 0.131825963, 0.000014628))
    col1.write(f'<h4>🦎 Tail:</h4>', unsafe_allow_html=True)
    col1.button("Try me!", key="tail_button", on_click=update_coordinates, args=(-0.7436499, 0.13188204, 0.00073801))

def update_coordinates(center_r, center_i, diameter):
    st.session_state["center_r"] = center_r
    st.session_state["center_i"] = center_i
    st.session_state["diameter"] = diameter    

if __name__ == '__main__':
    main()