## Generate the HTMX and GoLang code for a single-page weather application.

Title: Current Weather Worldwide

**Functionality:**

1.  Display an input field prominently labeled "Enter City Name" where the user can type a city.
2.  Display a button labeled "Get Weather".
3.  When the button is clicked:
    - Display a simple loading indicator (e.g., text like "Loading...") while the data is being fetched.
    - Fetch the current weather data for the entered city using the WeatherAPI `current.json` endpoint (`https://api.weatherapi.com/v1/current.json`).
    - Assume the necessary WeatherAPI API key is securely stored and accessible within the application's environment, represented in the code by a constant named `WEATHER_API_KEY`. **If possible on the platform, do not hardcode the key directly in the script.**
    - Once the data is received or an error occurs, hide the loading indicator.
    - If the API call is successful, display the following information clearly in a designated results area:
      - City Name (as returned by the API for confirmation)
      - Current Temperature (in Celsius °C)
      - Feels Like Temperature (in Celsius °C)
      - Humidity (in %)
      - Wind Speed (in kilometers per hour, kph)
    - Implement robust error handling: If the city is not found, the API key is invalid, or any other network/API error occurs, display a user-friendly error message (e.g., "Error: Could not retrieve weather data for [City Name]. Please verify the city name or check API key configuration."). Clear any previously displayed weather data when an error occurs.
    - Also add local form validation for the City input field that requires at least three letters (a-z, A-Z), but allows other characters. Validate after button click. Do not submit to the server from button click if validation fails. Use AlpineJS to perform the local validation and to display a validation error.

**Styling:**

- Apply a clean, modern, and visually appealing style using DaisyUI themes.
- Ensure the layout is responsive and looks good on both desktop and mobile screens.
- Center the core elements (input field, button, results area) on the page.
- Use clear typography, appropriate spacing, and padding for readability.
- Style the results area distinctly (e.g., using a card-like background or border).
- Use stylish icons for the app, the report card and the results as shown on the design image.

* Create a light mode design that matches docs/Screenshot.png in colors styles and icons (theme: fantasy). Also implement a dark mode capabale design (theme: aqua). Light and dark mode should switch automatically based on the browser settings.
* the app should be centered on the page and be mobile capable

You must create a UI similar to the attached image docs/Screenshot.png
If you cannot see the screenshot abort the task.

**Technical Requirements & Best Practices:**

- use the best practices for the Tech Stack specified in the AI Rules

You must use this key: WEATHER_API_KEY from the `.env` environment variables file.

API Documentation which you can access using the *fetch* tool:
https://www.weatherapi.com/docs/

Example call:
https://api.weatherapi.com/v1/current.json?key=WEATHER_API_KEY&q=London&aqi=no

Example Response Body:

```
{
    "location": {
        "name": "London",
        "region": "City of London, Greater London",
        "country": "United Kingdom",
        "lat": 51.5171,
        "lon": -0.1062,
        "tz_id": "Europe/London",
        "localtime_epoch": 1744266568,
        "localtime": "2025-04-10 07:29"
    },
    "current": {
        "last_updated_epoch": 1744265700,
        "last_updated": "2025-04-10 07:15",
        "temp_c": 7.3,
        "temp_f": 45.1,
        "is_day": 1,
        "condition": {
            "text": "Overcast",
            "icon": "//cdn.weatherapi.com/weather/64x64/day/122.png",
            "code": 1009
        },
        "wind_mph": 4.3,
        "wind_kph": 6.8,
        "wind_degree": 23,
        "wind_dir": "NNE",
        "pressure_mb": 1032.0,
        "pressure_in": 30.47,
        "precip_mm": 0.0,
        "precip_in": 0.0,
        "humidity": 70,
        "cloud": 100,
        "feelslike_c": 6.1,
        "feelslike_f": 43.0,
        "windchill_c": 4.6,
        "windchill_f": 40.3,
        "heatindex_c": 6.0,
        "heatindex_f": 42.8,
        "dewpoint_c": 2.6,
        "dewpoint_f": 36.8,
        "vis_km": 10.0,
        "vis_miles": 6.0,
        "uv": 0.1,
        "gust_mph": 5.6,
        "gust_kph": 9.0
    }
}
```
