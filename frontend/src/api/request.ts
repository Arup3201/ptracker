const API_ROOT = "http://localhost:8081/api";

type URLMethod = "GET" | "POST" | "PUT" | "DELETE";

export async function ApiRequest<TApi, TDomain>(
  endpoint: string,
  method: URLMethod,
  body: Record<string, any> | null,
  mapper: (data: TApi) => TDomain | null
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
            return await ApiRequest(endpoint, method, body, mapper);
          }
        }

        throw new Error(json.error.message);
      } else {
        throw new Error("Something went wrong there. Please try again.");
      }
    }

    if (json.data) {
      return mapper ? mapper(json.data) : json.data;
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
        const response = await fetch(API_ROOT + "/auth/refresh", {
          method: "POST",
          credentials: "include",
        });

        if (response.status == 200) {
          return await ApiRequest(endpoint, method, body, mapper);
        }

        throw new Error(json.error.message);
      } else {
        throw new Error("Something went wrong there. Please try again.");
      }
    }

    if (json.data) {
      return mapper ? mapper(json.data) : json.data;
    }
  }

  return null;
}
