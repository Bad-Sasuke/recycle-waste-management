interface ResponseAPI {
    message: string;
    data?: any;
}

const fetchData = async (url: string): Promise<ResponseAPI | string> => {
    try {
        // Make the GET request using Fetch API
        const response = await fetch(url);

        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
            const data: ResponseAPI = await response.json();
            return data.message
        }

        // Parse the response as JSON
        const data: ResponseAPI = await response.json();

        return data;
    } catch (error: any) {
        // Handle errors (e.g., network issues, invalid JSON, etc.)
        console.error("Error fetching data:", error);
        return error.message;
    }
};

const postData = async (url: string, data: any): Promise<ResponseAPI | string> => {
    try {
        // Make the POST request using Fetch API
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });

        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
            const data: ResponseAPI = await response.json();
            return data.message
        }

        // Parse the response as JSON
        const data: ResponseAPI = await response.json();

        return data;
    } catch (error: any) {
        // Handle errors (e.g., network issues, invalid JSON, etc.)
        console.error("Error fetching data:", error);
        return error.message;
    }
};

export { fetchData, postData }