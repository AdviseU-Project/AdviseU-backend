import os
import json

catalog_dir = "data/catalogs"
for filename in os.listdir(catalog_dir):
    if filename.endswith("_catalog.json"):
        filepath = os.path.join(catalog_dir, filename)
        with open(filepath, "r", encoding="utf-8") as f: # Specify UTF-8 encoding
            courses = json.load(f)

        # Determine the department from the filename
        department = filename.split("_")[0].upper()

        # Add the department field to each course
        for course in courses:
            course["department"] = department

        # Save the updated data back to the file
        with open(filepath, "w", encoding="utf-8") as f:
            json.dump(courses, f, indent=2)
