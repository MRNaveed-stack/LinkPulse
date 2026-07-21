import axios from 'axios';
import { useAuthStore } from '../store/authStore';

const apiClient = axios.create({
    baseURL: import.meta.env.VITE_API_URL || '/api',
    headers: {
        'Content-Type': 'application/json',
    },
});

apiClient.interceptors.request.use((config) => {
    const token = useAuthStore.getState().accessToken;
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

let isRefreshing = false;
let failedQueue: Array<{
    resolve: (value: unknown) => void;
    reject: (reason?: unknown) => void;
}> = [];

const processQueue = (error: unknown, token: string | null = null) => {
    failedQueue.forEach((prom) => {
        if (error) {
            prom.reject(error);
        } else {
            prom.resolve(token);
        }
    });
    failedQueue = [];
};

apiClient.interceptors.response.use(
    (response) => response,
    async (error) => {
        const originalRequest = error.config;

        // Check if the error is an unauthenticated error (401) and hasn't been retried yet
        if (error.response?.status === 401 && originalRequest && !originalRequest._retry) {
            
            // Scenario A: A token refresh is ALREADY happening. Put this request in the waiting room.
            if (isRefreshing) {
                return new Promise((resolve, reject) => {
                    failedQueue.push({ resolve, reject });
                })
                .then((token) => {
                    if (originalRequest.headers) {
                        originalRequest.headers.Authorization = `Bearer ${token}`;
                    }
                    return apiClient(originalRequest);
                })
                .catch((err) => Promise.reject(err));
            }

            // Scenario B: This is the first request to hit the 401 error. Start the refresh process.
            originalRequest._retry = true;
            isRefreshing = true;

            const refreshToken = useAuthStore.getState().refreshToken;
            if (!refreshToken) {
                useAuthStore.getState().logout();
                return Promise.reject(error);
            }

            try {
                // Use standard global axios here to prevent an infinite loop!
                const { data } = await axios.post(`${apiClient.defaults.baseURL}/auth/refresh`, {
                    refresh_token: refreshToken,
                });
                
                const newAccessToken = data.access_token;
                const newRefreshToken = data.refresh_token;

                // Save new tokens to global state
                useAuthStore.getState().setTokens({
                    accessToken: newAccessToken,
                    refreshToken: newRefreshToken,
                });

                // Wake up all pending requests in the waiting room and feed them the new token
                processQueue(null, newAccessToken);

                // Retry the original request that kicked off this whole process
                if (originalRequest.headers) {
                    originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
                }
                return apiClient(originalRequest);

            } catch (refreshError) {
                // If the refresh token itself is expired/bad, fail everything and log out
                processQueue(refreshError, null);
                useAuthStore.getState().logout();
                return Promise.reject(refreshError);
            } finally {
                isRefreshing = false;
            }
        }
        
        return Promise.reject(error);
    }
);

export default apiClient;
