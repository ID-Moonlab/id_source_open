#! -*- coding: utf-8 -*-
import os
import json
import time
import pandas as pd
from docx import Document

SECRET_TITLES = ["Title-A", "Title-B"]

def get_script_dir():
    return os.path.dirname(os.path.abspath(__file__))

def get_file_path(inp_param):
    return os.path.join(get_script_dir(), inp_param)

def jumbled_process(inp_param):
    doc = Document(get_file_path(inp_param))
    temp_result = {title: [] for title in SECRET_TITLES}
    final_result = {}

    for table in doc.tables:
        title_found = None
        group_list = []

        for row in table.rows:
            cell_data = []
            counter = 1

            for cell in row.cells:
                cell_text = cell.text
                if not cell_text:
                    break
                if cell_text in SECRET_TITLES:
                    title_found = cell_text
                    break
                if cell_text == "tax":
                    if counter == 1:
                        cell_data.append(f"{cell_text}-{counter}")
                        counter += 1
                    elif counter == 2:
                        cell_data.append(f"{cell_text}-{counter}")
                        counter = 1
                else:
                    cell_data.append(cell_text)

            if cell_data:
                group_list.append(cell_data)

            if "sum" in cell_data:
                temp_result[title_found].append(group_list)
                group_list = []

    for key, value in temp_result.items():
        processed_list = []
        for matrix in value:
            headers = matrix[0]
            data_rows = matrix[1:]
            df = pd.DataFrame(data_rows, columns=headers)
            for _, row in df.iterrows():
                key_value_dict = row.to_dict()
                processed_list.append(key_value_dict)
        final_result[key] = processed_list

    return final_result

def main_jumbled(input_file, output_file):
    result = jumbled_process(input_file)
    if os.path.exists(output_file):
        os.remove(output_file)
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(result, f, ensure_ascii=False, indent=4)


if __name__ == "__main__":
    user_input = input("Please input your file name: ")
    print(f"Input file is: {user_input}")
    output_file = user_input.split(".")[0] + ".json"
    print(f"Output file is: {output_file}")
    try:
        main_jumbled(user_input, output_file)
        time.sleep(2)
        print()
        print("...transform success...OK!")
    except Exception as e:
        time.sleep(2)
        print()
        print("...transform error...Failed!")
        print(e)
