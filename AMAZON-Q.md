### 💡 **Implementation Prompt for Amazon Q**

Implement an SDK initialization flow for ENBUILD that performs authentication and client creation based on admin settings from a remote API.

#### ✅ **Step-by-Step Behavior**

1. **Fetch Admin Settings**

   * Construct the Admin Settings URL:

   ENBUILD_BASE_URL="https://enbuild-dev.vivplatform.io/" then append the path:
    `/enbuild-user/api/v1/adminSettings` so it becomes:

     `https://enbuild-dev.vivplatform.io/enbuild-user/api/v1/adminSettings`
   * Make a GET request to this endpoint.
   * If the response is not received or fails (non-2xx), terminate the process and log a meaningful error message such as:
     `"Failed to fetch authMechanism from ENBUILD. Please check ENBUILD_BASE_URL or network connectivity."`

2. **Extract `authMechanism`**

   * From the JSON response, extract the value of the key `authMechanism`.

3. **Authenticate Based on `authMechanism`**

   * **If `authMechanism` is `keycloak`:**

     * Extract the following Keycloak details from the response:

       ```json
       "adminConfigs": {
         "keycloak": {
           "KEYCLOAK_BACKEND_URL": "...",
           "KEYCLOAK_CLIENT_ID": "...",
           "KEYCLOAK_REALM": "..."
         }
       }
       ```
     * Use the following credentials:

       ```
       ENBUILD_USERNAME="juned"
       ENBUILD_PASSWORD="juned"
       ```
     * Perform a password grant token request to Keycloak:

       * Token URL format:
         `{KEYCLOAK_BACKEND_URL}/realms/{KEYCLOAK_REALM}/protocol/openid-connect/token`
       * Request `client_id`, `username`, `password`, and `grant_type=password` as form parameters.
     * If token acquisition fails, error out with a clear message like:
       `"Authentication with Keycloak failed. Check credentials or Keycloak settings."`
     * On success, extract the `access_token`.

   * **If `authMechanism` is `local`:**

     * Use a default static token:
       `enbuild_local_admin_token`

4. **Initialize ENBUILD Client**

   * Use the acquired token (either from Keycloak or local static token) to instantiate and authenticate the ENBUILD client.
   * Proceed with the client usage from this point.

---

### 🧪 Notes for Implementation

* Ensure all logs and errors are precise and help in debugging failures.
* Abstract the token-fetching logic so it can support additional auth mechanisms in the future if needed.
* Avoid hardcoding URLs or tokens—use environment variables where applicable.

