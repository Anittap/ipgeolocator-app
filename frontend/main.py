from flask import Flask, render_template, make_response, request
import requests
import re
import os

app = Flask(__name__)

# Define routes
@app.route('/', strict_slashes=False)
@app.route('/ip', strict_slashes=False)
@app.route('/ip/<ip>', strict_slashes=False)
def index(ip=None):
    hostname = os.getenv("HOSTNAME", "none")
    api_server = os.getenv("API_SERVER", None)
    api_server_port = os.getenv("API_SERVER_PORT", None)

    # Validate API server and port
    if not api_server or not api_server_port:
        return render_template('error.html', error_message="API server configuration is missing.")

    if ip:
        # Validate IP address
        pattern = re.compile(r"^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$")
        valid_ip = pattern.match(ip.strip())
        
        if valid_ip:
            # Fetch data from the backend API
            api_url = f"http://{api_server}:{api_server_port}/ip/{ip}"
            try:
                response = requests.get(url=api_url, timeout=5)
                if response.status_code == 200:
                    geodata = response.json()
                    return render_template(
                        'index.html',
                        continent_name=geodata.get('continent_name', 'Unknown'),
                        continent_code=geodata.get('continent', 'Unknown'),
                        country_name=geodata.get('country_name', 'Unknown'),
                        isp=geodata.get('isp', 'Unknown'),
                        cached=geodata.get('cached', 'Unknown'),
                        host=hostname,
                        apiServer=geodata.get('apiServer', 'Unknown'),
                        apiServerVersion=geodata.get('version', 'Unknown'),
                        frontendServerVersion="version1",
                        hostname=hostname
                    )
                else:
                    return render_template('error.html', error_message=f"Error fetching data from API. Status code: {response.status_code}")
            except requests.RequestException as e:
                return render_template('error.html', error_message=f"Error connecting to API: {e}")
        else:
            return render_template('error.html', error_message="Invalid IP address format.")
    else:
        return render_template('error.html', error_message="No IP address provided.")

@app.route('/status', strict_slashes=False)
def check_status():
    return make_response("", 200)

if __name__ == "__main__":
    hostname = os.getenv("HOSTNAME", "none")
    api_server = os.getenv("API_SERVER", None)
    api_server_port = os.getenv("API_SERVER_PORT", None)
    app_port = os.getenv('APP_PORT', "8080")
    app.run(port=int(app_port), host="0.0.0.0", debug=True)

