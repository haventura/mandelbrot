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
    colormap_name: str   ="BWY"

# run with python -m streamlit run main.py

def main():
    data = Data()
    st.set_page_config(page_title="Mandelbrot Generator", page_icon="🚀", layout="wide")
    st.title("Mandelbrot Generator")
    col1, col2, col3 = st.columns([1,2,1])
    if 'center_r' not in st.session_state:
        st.session_state['center_r'] = -.75
        st.session_state["center_i"] = .0
        st.session_state["diameter"] = 2.5
    with col1.form("my_form"):
        st.header("🧪 Parameters")
        center_r = st.number_input("Center Re(c)", min_value=-3., max_value=3., step=0.1, format="%.15g", key="center_r")
        center_i = st.number_input("Center Im(c)", min_value=-3., max_value=3., step=0.1, format="%.15g", key="center_i")
        diameter = st.number_input("Diameter", min_value=0., max_value=6., step=0.1, format="%.15g", key="diameter")
        data.max_iteration = st.number_input("Iterations", value=512, min_value=1, key="max_iteration")
        resolution = st.number_input("Resolution", value=1024, min_value=256, key="resolution")
        data.colormap_name = st.text_input("Color map", value="BWY", key="colormap")

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
            file_content=base64.b64decode(image_data.split(',')[1])
            col2.image(file_content, use_column_width="always")
            col1.download_button("💾 Download", file_content, "mandelbrot.png", key="download_button")

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

if __name__ == '__main__':
    main()