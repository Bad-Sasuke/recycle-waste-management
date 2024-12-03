interface ResponseAPI {
    message: string;
    data?: any;
}

const fetchData = async (url: string): Promise<ResponseAPI | null> => {
    try {
        // Make the GET request using Fetch API
        const response = await fetch(url);

        // Check if the response status is OK (status code 200-299)
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }

        // Parse the response as JSON
        const data: ResponseAPI = await response.json();

        return data;
    } catch (error: any) {
        // Handle errors (e.g., network issues, invalid JSON, etc.)
        console.error("Error fetching data:", error);
        return null; // Return null or handle the error as needed
    }
};


export { fetchData }