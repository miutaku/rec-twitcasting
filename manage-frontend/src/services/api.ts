import axios from 'axios';

const API_BASE_URL = 'http://localhost:8888';

export const addCastingUser = async (username) => {
    try {
        const response = await axios.get(`${API_BASE_URL}/add-casting-user`, {
            params: { username }
        });
        return response.data;
    } catch (error) {
        return { error: error.response?.data?.error || 'An error occurred' };
    }
};

// 他のAPI関数をここに追加できます。