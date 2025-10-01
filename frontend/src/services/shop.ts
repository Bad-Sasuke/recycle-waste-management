import { getCookie } from '@/stores/cookie';

interface ResponseAPI {
  message: string;
  data?: any;
  page?: number;
  limit?: number;
  total_pages?: number;
  total_items?: number;
}

// Get all shops with pagination
const fetchShops = async (page: number = 1, limit: number = 12): Promise<ResponseAPI> => {
  try {
    const apiUrl = import.meta.env.VITE_WEB_API;
    if (!apiUrl) {
      console.error('API configuration error. Please check environment settings.');
      return { message: 'API configuration error' };
    }

    const response = await fetch(`${apiUrl}/api/shops/get-shops?page=${page}&limit=${limit}`);

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

    const data: ResponseAPI = await response.json();
    return data;
  } catch (error) {
    console.error('Error fetching shops:', error);
    return {
      message: 'An error occurred while fetching shops.',
      data: []
    };
  }
};

// Get shop by shop ID
const fetchShopById = async (shopId: string): Promise<ResponseAPI> => {
  try {
    const apiUrl = import.meta.env.VITE_WEB_API;
    if (!apiUrl) {
      console.error('API configuration error. Please check environment settings.');
      return { message: 'API configuration error' };
    }

    const response = await fetch(`${apiUrl}/api/shops/get-shop/${shopId}`);

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

    const data: ResponseAPI = await response.json();
    return data;
  } catch (error) {
    console.error('Error fetching shop by ID:', error);
    return {
      message: 'An error occurred while fetching shop.',
      data: null
    };
  }
};

// Get current user's shop
const fetchMyShop = async (): Promise<ResponseAPI> => {
  try {
    const apiUrl = import.meta.env.VITE_WEB_API;
    const token = getCookie('token');
    if (!apiUrl || !token) {
      console.error('API configuration error or missing token. Please check environment settings.');
      return { message: 'API configuration error or missing token' };
    }

    const response = await fetch(`${apiUrl}/api/shops/my-shop`, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

    const data: ResponseAPI = await response.json();
    return data;
  } catch (error) {
    console.error('Error fetching my shop:', error);
    return {
      message: 'An error occurred while fetching your shop.',
      data: null
    };
  }
};

// Create a new shop
const createShop = async (formData: FormData): Promise<ResponseAPI> => {
  try {
    const apiUrl = import.meta.env.VITE_WEB_API;
    const token = getCookie('token');
    if (!apiUrl || !token) {
      console.error('API configuration error or missing token. Please check environment settings.');
      return { message: 'API configuration error or missing token' };
    }

    const response = await fetch(`${apiUrl}/api/shops/create-shop`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
      body: formData
    });

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

    const data: ResponseAPI = await response.json();
    return data;
  } catch (error) {
    console.error('Error creating shop:', error);
    return {
      message: 'An error occurred while creating shop.',
      data: null
    };
  }
};

// Update a shop
const updateShop = async (shopId: string, formData: FormData): Promise<ResponseAPI> => {
  try {
    const apiUrl = import.meta.env.VITE_WEB_API;
    const token = getCookie('token');
    if (!apiUrl || !token) {
      console.error('API configuration error or missing token. Please check environment settings.');
      return { message: 'API configuration error or missing token' };
    }

    const response = await fetch(`${apiUrl}/api/shops/update-shop/${shopId}`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
      body: formData
    });

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

    const data: ResponseAPI = await response.json();
    return data;
  } catch (error) {
    console.error('Error updating shop:', error);
    return {
      message: 'An error occurred while updating shop.',
      data: null
    };
  }
};

// Delete a shop
const deleteShop = async (shopId: string): Promise<ResponseAPI> => {
  try {
    const apiUrl = import.meta.env.VITE_WEB_API;
    const token = getCookie('token');
    if (!apiUrl || !token) {
      console.error('API configuration error or missing token. Please check environment settings.');
      return { message: 'API configuration error or missing token' };
    }

    const response = await fetch(`${apiUrl}/api/shops/delete-shop/${shopId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

    const data: ResponseAPI = await response.json();
    return data;
  } catch (error) {
    console.error('Error deleting shop:', error);
    return {
      message: 'An error occurred while deleting shop.',
      data: null
    };
  }
};

export {
  fetchShops,
  fetchShopById,
  fetchMyShop,
  createShop,
  updateShop,
  deleteShop
};