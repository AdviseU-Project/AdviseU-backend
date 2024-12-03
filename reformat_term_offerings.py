import os
import json

# Define the mapping for terms
term_mapping = {
    'fall_2024': 'Fall 2024',
    'winter_2025': 'Winter 2025',
    'spring_2025': 'Spring 2025'
}

# Path to the term_offerings directory
term_offerings_path = 'data/term_offerings'

# Iterate through each term directory (fall_2024, winter_2025, etc.)
for term_dir in os.listdir(term_offerings_path):
    term_dir_path = os.path.join(term_offerings_path, term_dir)

    # Skip if it's not a directory
    if not os.path.isdir(term_dir_path):
        continue

    # Get the term name from the mapping
    term_name = term_mapping.get(term_dir, None)

    if term_name:
        department_data = []
        curr_department = ""

        # Iterate through the JSON files in the current term directory
        for file_name in os.listdir(term_dir_path):
            if file_name.endswith('.json'):
                file_path = os.path.join(term_dir_path, file_name)

                # Determine the department from the filename
                department = file_name.split("_")[0].upper()
                
                # Save all data for a department during a specific term into a JSON file
                if curr_department != department and curr_department != "":
                    # Create the combined file name (e.g., WSE_fall_2024.json)
                    output_file_name = f'{curr_department}_{term_dir}.json'
                    output_file_path = os.path.join(term_offerings_path, output_file_name)

                    # Create the output JSON structure
                    output_data = [{
                        "department": curr_department,
                        "term": term_name,
                        "courses": department_data
                    }]

                    # Write the combined courses to the new file
                    with open(output_file_path, 'w') as file:
                        json.dump(output_data, file, indent=4)

                    # Reset department_data for the next department
                    department_data = []

                    print(f'Created combined file for {curr_department} in {term_name}: {output_file_name}')

                # Update current department
                if curr_department != department:
                    curr_department = department

                # Determine the code from the filename
                code = file_name.split("_")[0].upper() + ' ' + file_name.split("_")[1].upper()
    
                # Read the JSON data
                with open(file_path, 'r') as file:
                    data = json.load(file)

                # Add the term field to each offering and get the course's title
                title = ''
                for course in data.get('results', []):
                    course['term'] = term_name
                    title = course.get('title', '')

                # Remove srcdb from the outer JSON
                data.pop('srcdb', None)

                # Add term, code, and title to the outer JSON
                data['term'] = term_name
                data['code'] = code
                data['title'] = title

                # Wrap the entire data in a list
                department_data.append(data)

        # Handle the last department after the loop
        if curr_department and department_data:
            output_file_name = f'{curr_department}_{term_dir}.json'
            output_file_path = os.path.join(term_offerings_path, output_file_name)

            # Create the output JSON structure
            output_data = [{
                "department": curr_department,
                "term": term_name,
                "courses": department_data
            }]

            # Write the combined courses to the new file
            with open(output_file_path, 'w') as file:
                json.dump(output_data, file, indent=4)

            print(f'Created combined file for {curr_department} in {term_name}: {output_file_name}')
    else:
        print(f'No mapping found for directory: {term_dir}')
