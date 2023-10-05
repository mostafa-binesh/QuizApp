import os

# Define the path to the directory you want to start from
start_path = r"C:\Users\UbiRo\Desktop\next generation 2\New Gneerations"
# Function to print folder and subfolder names
def print_folders_and_subfolders(path, prefix=""):
    if prefix == "":
        print('new folder')
    for item in os.listdir(path):
        item_path = os.path.join(path, item)
        if os.path.isdir(item_path):
            print(f'"{prefix}{item}",')
            print_folders_and_subfolders(item_path, prefix + "---")

# Start printing from the specified path
print_folders_and_subfolders(start_path)
