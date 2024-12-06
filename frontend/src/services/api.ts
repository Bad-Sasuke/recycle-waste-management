interface ResponseAPI {
    message: string;
    data?: string[] | { [key: string]: unknown };
}

const fetchData = async (url: string): Promise<ResponseAPI> => {
    try {
        // Make the GET request using Fetch API
        const response = await fetch(url);

        if (!response.ok) {
            // Throw an error if the response is not okay
            throw new Error(`HTTP error! Status: ${response.status}`);
        }

        // Parse the response as JSON
        const data: ResponseAPI = await response.json();

        return data;
    } catch (error) {
        // Handle errors (e.g., network issues, invalid JSON, etc.)
        console.error("Error fetching data:", error);

        // Return a default ResponseAPI object in case of an error
        return {
            message: 'An error occurred while fetching data.',
            data: []  // Empty array as a fallback value
        };
    }
};

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const postData = async (url: string, data: unknown): Promise<ResponseAPI | string> => {
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
            // Throw an error if the response is not okay
            throw new Error(`HTTP error! Status: ${response.status}`);
        }

        // Parse the response as JSON
        const responseData: ResponseAPI = await response.json();

        return responseData;
    } catch (error: unknown) {
        // Handle errors (e.g., network issues, invalid JSON, etc.)
        console.error("Error fetching data:", error);

        // Check if the error is an instance of Error and return the message
        if (error instanceof Error) {
            return error.message;
        }

        // Return a default error message if the error type is unknown
        return 'An unknown error occurred.';
    }
};


export { fetchData, postData }