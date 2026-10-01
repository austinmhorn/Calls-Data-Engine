import gspread
from oauth2client.service_account import ServiceAccountCredentials
import csv
import os

# Define scope and credentials
scope = ["https://spreadsheets.google.com/feeds", "https://www.googleapis.com/auth/drive"]
creds = ServiceAccountCredentials.from_json_keyfile_name("credentials.json", scope)
client = gspread.authorize(creds)

spreadsheet = client.open("Calls Data Feed")
csv_files = [f for f in os.listdir(".") if f.endswith(".csv")]

for csv_file in csv_files:
    sheet_name = os.path.splitext(csv_file)[0]
    print(f"Uploading {csv_file} to Google Sheet: {sheet_name}")

    try:
        sheet = spreadsheet.worksheet(sheet_name)
    except gspread.exceptions.WorksheetNotFound:
        print(f"Worksheet '{sheet_name}' not found. Creating a new one.")
        sheet = spreadsheet.add_worksheet(title=sheet_name, rows=1000, cols=26)

    # Read CSV
    with open(csv_file, "r", encoding="utf-8") as f:
        data = list(csv.reader(f))

    if not data:
        sheet.clear()
        continue

    rows_needed = len(data)
    cols_needed = max(len(r) for r in data)

    # Ensure the grid is big enough (handles existing sheets too)
    target_rows = max(sheet.row_count, rows_needed, 1000)
    target_cols = max(sheet.col_count, cols_needed, 26)
    if target_rows != sheet.row_count or target_cols != sheet.col_count:
        sheet.resize(rows=target_rows, cols=target_cols)

    # Clear then write the whole matrix at once (formulas still work)
    sheet.clear()
    sheet.update("A1", data, value_input_option="USER_ENTERED")

    print(f"{csv_file} synced to Google Sheet '{sheet_name}' successfully.")