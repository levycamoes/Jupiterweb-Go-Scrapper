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
