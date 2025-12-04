const URL_ROOT = window.location.origin + "/api/v1";

const HttpGet = async (resource: string) => {
  try {
    const response = await fetch(URL_ROOT + resource);

    if (!response.ok) {
      throw new Error(`Response status: ${response.status}`);
    }

    return response.json();
  } catch (err) {
    if (err instanceof Error) {
      console.error(`HttpGet fetch call error: ${err.message}`);
      throw new Error(err.message);
    }
  }
};
const HttpPost = async (resource: string, payload: any) => {
  try {
    const response = await fetch(URL_ROOT + resource, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });

    if (!response.ok) {
      throw new Error(`Response status: ${response.status}`);
    }

    return response.json();
  } catch (err) {
    if (err instanceof Error) {
      console.error(`HttpPost fetch call error: ${err.message}`);
      throw new Error(err.message);
    }
  }
};
const HttpPut = async (resource: string, payload: any) => {
  try {
    const response = await fetch(URL_ROOT + resource, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });

    if (!response.ok) {
      throw new Error(`Response status: ${response.status}`);
    }

    return response.json();
  } catch (err) {
    if (err instanceof Error) {
      console.error(`HttpPost fetch call error: ${err.message}`);
      throw new Error(err.message);
    }
  }
};
const HttpDelete = async (resource: string) => {
  try {
    const response = await fetch(URL_ROOT + resource, {
      method: "DELETE",
    });

    if (!response.ok) {
      throw new Error(`Response status: ${response.status}`);
    }

    return response.json();
  } catch (err) {
    if (err instanceof Error) {
      console.error(`HttpPost fetch call error: ${err.message}`);
      throw new Error(err.message);
    }
  }
};

export { HttpGet, HttpPost, HttpPut, HttpDelete };
