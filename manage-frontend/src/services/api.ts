// filepath: src/services/api.ts
import axios from 'axios';

const API_BASE_URL = 'http://manage-backend-rec-twitcasting:8888';

export const fetchCastingUsers = async () => {
    try {
        const response = await axios.get(`${API_BASE_URL}/list-casting-users`);
        return response.data;
    } catch (error) {
        return { error: (error as any).response?.data?.error || 'An error occurred' };
    }
};

export const addCastingUser = async (username: string) => {
    try {
        const response = await axios.get(`${API_BASE_URL}/add-casting-user`, {
            params: { username }
        });
        return response.data;
    } catch (error: any) {
        return { error: error.response?.data?.error || 'An error occurred' };
    }
};