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
			config.headers['Authorization'] = `Bearer ${token}`;
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
		if (error.response && error.response.status === 401) {
			// Token is invalid or expired
			localStorage.removeItem('token');
			window.location.href = '/#/session';
		}
		return Promise.reject(error);
	}
);

export default instance;
