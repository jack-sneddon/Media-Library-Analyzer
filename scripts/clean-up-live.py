import os
import shutil

def clean_up_mov_files(source_folder):
    """
    Cleans up MOV files associated with Live Photos from the source folder.

    Args:
        source_folder: Path to the source folder containing the exported images.
    """

    for root, dirs, files in os.walk(source_folder):
        for file in files:
            if file.lower().endswith(".heic"):
                heic_name = os.path.splitext(file)[0] 
                mov_path = os.path.join(root, heic_name + ".mov")

                if os.path.exists(mov_path):
                    try:
                        os.remove(mov_path)
                        print(f"Deleted {mov_path}")
                    except OSError as e:
                        print(f"Error deleting {mov_path}: {e}")

if __name__ == "__main__":
    source_folder = "/Volumes/Media Drive/Photos/Outdoors/Goblin_Valley_2024"  # Replace with your import path

    clean_up_mov_files(source_folder)
    print("Cleanup complete!")
