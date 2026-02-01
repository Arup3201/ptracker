export const API_ROOT = "http://localhost:8081/api/v1";

type URLMethod = "GET" | "POST" | "PUT" | "DELETE";

export async function ApiRequest<TDomain>(
  endpoint: string,
  method: URLMethod,
  body: Record<string, any> | null,
): Promise<TDomain | null> {
  if (method == "GET") {
    const response = await fetch(API_ROOT + endpoint, {
      credentials: "include",
    });
    const json = await response.json();
    if (response.status != 200) {
      if (json.status == "error") {
        if (response.status == 401) {
          // refresh
          const response = await fetch(API_ROOT + "/auth/refresh", {
            method: "POST",
            credentials: "include",
          });

          if (response.status == 200) {
            return await ApiRequest(endpoint, method, body);
          }
        }

        if (response.status === 401) {
          window.location.href = "/login";
        }

        throw new Error(json.error.message);
      } else {
        throw new Error("Something went wrong there. Please try again.");
      }
    }

    if (json.data) {
      return json.data;
    }
  } else {
    const response = await fetch(API_ROOT + endpoint, {
      method: method,
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: body ? JSON.stringify(body) : null,
    });
    const json = await response.json();
    if (response.status >= 400) {
      if (json.status == "error") {
        if (response.status == 401) {
          const response = await fetch(API_ROOT + "/auth/refresh", {
            method: "POST",
            credentials: "include",
          });

          if (response.status == 200) {
            return await ApiRequest(endpoint, method, body);
          }
        }

        throw new Error(json.error.message);
      } else {
        throw new Error("Something went wrong there. Please try again.");
      }
    }

    if (json.data) {
      return json.data;
    }
  }

  return null;
}
