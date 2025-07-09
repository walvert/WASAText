import axios from "axios";

const instance = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5
});

// Add a request interceptor to include auth token in all requests
instance.interceptors.request.use(
	config => {
		const token = localStorage.getItem('token');
		if (token) {
			config.headers['Authorization'] = `${token}`;
		}
		return config;
	},
	error => {
		return Promise.reject(error);
	}
);

// Add a response interceptor to handle authentication errors
instance.interceptors.response.use(
	response => {
		return response;
	},
	error => {
		console.error('API Error:', {
			url: error.config?.url,
			status: error.response?.status,
			data: error.response?.data,
			headers: error.response?.headers
		});

		if (error.response && error.response.status === 401) {
			console.log('401 error detected, removing token and redirecting to login');

			// Token is invalid or expired
			localStorage.removeItem('token');
			localStorage.removeItem('user');

			// Only redirect if we're not already on the session page
			if (window.location.hash !== '#/session') {
				window.location.href = '/#/session';
			}
		}
		return Promise.reject(error);
	}
);

export default instance;
