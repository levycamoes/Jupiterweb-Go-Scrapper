# 🎓 USP PUB Crawler (Scholarship Scraper)
This repository contains Web Scraping tools designed to extract public data from the Unified Scholarship Program (PUB) of the University of São Paulo (USP), hosted on the JupiterWeb system.

The goal is to create a structured dataset (CSV) containing information about calls for applications, research projects, advisors, and scholarship distribution, enabling academic data analysis.

# 🚀 Project Approaches

This project is a second version of a script that I've made in the past to scrape the same data, but using python and selenium.

## About this version
Built using the Colly framework. This version runs headless (no browser UI).

Technique: Reverse Engineering of the DWR (Direct Web Remoting) protocol.

Key Feature: It simulates the generation of security tokens (scriptSessionId and DWRSESSIONID) and Java session cookies programmatically, communicating directly with the server's API.

Advantage: Extremely lightweight, fast, and resource-efficient.

## 🧠 Technical Details (DWR Reverse Engineering)

The JupiterWeb system uses a legacy technology called DWR (Direct Web Remoting), which allows JavaScript to invoke Java methods on the server.

The challenge for the Go crawler was to replicate the DWR security handshake without executing JavaScript. The script performs the following steps:

1. Login: Authenticates via a standard POST request and stores the JSESSIONID cookie.

2. Handshake: Locally generates an alphanumeric DWRSESSIONID and injects it into the request cookies.

3. Tokenization: Calculates the scriptSessionId by combining the session ID with a random suffix, as expected by the USP Java backend.

4. Parsing: The server response is not JSON, but executable JavaScript. The script uses Regex (Regular Expressions) to clean the response and extract structured data.

## ⚠️ Owner disclaimer

This software was developed for strictly educational purposes and public data analysis. Do not use this script to overload the system. This script DOES NOT get the full name of the project advisor.

## ⚙️ Function Breakdown

Here is a breakdown of the main functions in the `main.go` file:

- **`main`**: The entry point of the application. It initializes the Colly collector, sets up the CSV writer for storing the scraped data, handles the login process to the USP digital system, and orchestrates the entire web scraping workflow by calling the other functions in sequence.

- **`prepararSessaoDWR`**: This function prepares the session for DWR (Direct Web Remoting) communication. It mimics the behavior of a real browser by generating a fake `DWRSESSIONID` and a `scriptSessionId`. These are crucial for making authenticated requests to the DWR API and retrieving the scholarship data.

- **`dispararDWR`**: This function is responsible for sending the main data request to the DWR API. It constructs the appropriate request headers and payload, then sends the request to fetch the scholarship information for the specified year. It also handles the server's response, passing it to `parseDWRResponse` for processing.

- **`parseDWRResponse`**: This function parses the raw DWR response, which is a JavaScript-like string rather than a standard JSON format. It uses regular expressions to meticulously extract the relevant scholarship data, such as the year, university unit, project title, and number of scholarships. The extracted data is then written to the CSV file.

- **`extract`**: A utility function that simplifies the process of data extraction. It takes a regular expression and a block of text as input, and returns the first matching substring. This is used throughout `parseDWRResponse` to pull out specific pieces of information.

- **`unquoteUnicode`**: A helper function that decodes Unicode escape sequences (e.g., `\u00E3`) into their proper characters. This ensures that special characters and accents in the scraped data are correctly represented in the final CSV file.

- **`min`**: A simple utility function that returns the smaller of two integers. It is used to prevent out-of-bounds errors when logging server error messages.
