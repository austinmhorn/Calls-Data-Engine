import pandas as pd

def insert_vlookup_column(
    csv_path,
    output_path,
    new_column_name="Regional Manager",
    lookup_column_letter="A",
    table_array="property_data!A:CP",
    col_index=20,
    start_row=2
):
    # Load the CSV
    df = pd.read_csv(csv_path)

    # Generate VLOOKUP formulas dynamically
    formulas = [
        f'=IFERROR(VLOOKUP({lookup_column_letter}{i}, {table_array}, {col_index}, FALSE), "Not Found")'
        for i in range(start_row, start_row + len(df))
    ]

    # Insert the new column
    df[new_column_name] = formulas

    # Save the updated CSV
    df.to_csv(output_path, index=False)
    print(f"Saved with VLOOKUP column: {output_path}")

# Example usage:
insert_vlookup_column(
    csv_path="notion_data.csv",
    output_path="notion_data.csv"
)
