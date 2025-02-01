export const getKey = async (key) => {
    try {
      const response = await fetch(`http://localhost:9090/keys/${key}`);
      return await response.json();
    } catch (error) {
      return { error: "Failed to fetch key" };
    }
  };
  
  export const setKey = async (key, value) => {
    try {
      const response = await fetch(`http://localhost:9090/keys`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ key, value }),
      });
      return await response.json();
    } catch (error) {
      return { error: "Failed to set key" };
    }
  };
  
  export const deleteKey = async (key) => {
    try {
      const response = await fetch(`http://localhost:9090/keys/${key}`, {
        method: "DELETE",
      });
      return await response.json();
    } catch (error) {
      return { error: "Failed to delete key" };
    }
  };