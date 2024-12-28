const apiURL = import.meta.env.DEV
  ? 'http://127.0.0.1:8001/api'
  : window.location.origin + '/api';

export { apiURL };
