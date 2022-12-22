import streamlit as st
import requests
import json
import dataclasses

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
    st.set_page_config(page_title="Mandelbrot renderer", page_icon="🤖")
    st.title("Mandelbrot")
    with st.form("my_form"):
        data.min_r = st.number_input("Minimum real value", min_value=-2, max_value=2, key="min_r")
        data.max_r = st.number_input("Maximum real value", min_value=-2, max_value=2, key="max_r")
        data.min_i = st.number_input("Minimum imaginary value", min_value=-2, max_value=2, key="min_i")
        data.max_i = st.number_input("Maximum imaginary value", min_value=-2, max_value=2, key="max_i")
        data.max_iteration = st.number_input("Iterations", min_value=1, key="max_iteration")
        data.colormap_name = st.text_input("Color map", key="colormap")
        data.width = 1024
        data.height = 1024

        submitted = st.form_submit_button("Submit")

        if submitted:
            st.write("Result")
            url = 'http://localhost:5000/compute'
            print(data)
            x = requests.post(url, json = json.dumps(dataclasses.asdict(data)))
            st.write(x)


if __name__ == '__main__':
    main()