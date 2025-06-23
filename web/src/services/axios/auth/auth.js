import axios from "axios";
import { ErrorData } from "../../error/error";

const token = localStorage.getItem("auth_token") || null;

const axiosAuthenticated = axios.create({
  baseURL: import.meta.env.VITE_BASE_URL_DEV,
  headers: {
    "Content-Type": "application/json",
    Authorization: "Bearer " + token,
  },
});

class AxiosQueryWithAuth {
  constructor(axios) {
    this.axios = axios;
  }

  async getAllLocations() {
    try {
      const response = await this?.axios?.get("/v1/locations");
      return response;
    } catch (error) {
      throw new ErrorData(
        error.message,
        error?.response?.status,
        "fetching",
        "error while fetching locations"
      );
    }
  }
  async getAllAmenties() {
    try {
      const response = await this?.axios?.get("/v1/amenities");
      return response;
    } catch (error) {
      throw new ErrorData(
        error.message,
        error?.response?.status,
        "fetching",
        "error while fetching amenities"
      );
    }
  }

  async getAllCategories() {
    try {
      const response = await this?.axios?.get("/v1/categories");
      return response;
    } catch (error) {
      throw new ErrorData(
        error.message,
        error?.response?.status,
        "fetching",
        "error while fetching categories"
      );
    }
  }
}

const axiosQueryWithAuth = new AxiosQueryWithAuth(axiosAuthenticated);

export { axiosQueryWithAuth };
