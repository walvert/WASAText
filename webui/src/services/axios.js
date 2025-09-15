import axios from "axios";

const instance = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5
});

// Request interceptor
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

// Response interceptor to handle authentication errors
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

			localStorage.removeItem('token');
			localStorage.removeItem('user');

			if (window.location.hash !== '#/session') {
				window.location.href = '/#/session';
			}
		}
		return Promise.reject(error);
	}
);

export default instance;
