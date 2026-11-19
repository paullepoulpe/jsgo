class DataService { async fetchData(url) { return url ? { status: 200 } : null; } }
const svc = new DataService();
