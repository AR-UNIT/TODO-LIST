import requests
import time

# Step 1: Get the JWT token
def get_token():
    url = "http://localhost:8080/api/authenticate"
    headers = {
        "client_id": "admin",
        "client_secret": "password",
        "Content-Type": "application/json"
    }

    response = requests.get(url, headers=headers)
    if response.status_code == 200:
        token = response.json().get("token")
        print(f"Token retrieved: {token}")
        return token
    else:
        print(f"Failed to retrieve token. Status Code: {response.status_code}, Response: {response.text}")
        return None

# Step 2: Test rate limiting
def test_rate_limiting(token):
    url = "http://localhost:8080/tasks"
    headers = {
        "Authorization": f"Bearer {token}",
        "client_id": "admin",
        "client_secret": "password",
        "Content-Type": "application/json"
    }

    number_of_requests = 20
    # delay_between_requests = 0.2  # Delay in seconds
    delay_between_requests = 0.0  # Delay in seconds

    for i in range(1, number_of_requests + 1):
        response = requests.get(url, headers=headers)
        print("Request response: ", response.text)
        if response.status_code == 200:
            print(f"Request {i}: Success. Status Code: {response.status_code}")
        elif response.status_code == 429:
            print(f"Request {i}: Rate limit exceeded! Status Code: {response.status_code}")
        else:
            print(f"Request {i}: Failed. Status Code: {response.status_code}, Response: {response.text}")
        time.sleep(delay_between_requests)

# Main script execution
if __name__ == "__main__":
    token = get_token()
    if token:
        test_rate_limiting(token)
