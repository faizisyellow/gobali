import axios from "axios";
import { ErrorData } from "../../error/error";

const axiosAuthenticated = axios.create({
  baseURL: import.meta.env.VITE_BASE_URL_DEV,
  headers: {
    "Content-Type": "application/json",
  },
});

axiosAuthenticated.interceptors.request.use(
  function (config) {
    const token = localStorage.getItem("auth_token") || null;
    config.headers.Authorization = "Bearer " + token;

    return config;
  },
  function (error) {
    return Promise.reject(error);
  }
);

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

  async CreateNewVilla(payloads) {
    try {
      const response = await this?.axios?.post("/v1/villas", payloads, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      });

      return response;
    } catch (error) {
      throw new ErrorData(
        error.message,
        error?.response?.status,
        "mutation",
        "error while create new villa"
      );
    }
  }

  /* GET ALL VILLAS */
  async GetAllVillas(query) {
    try {
      const response = await this.axios.get(
        `/v1/villas?location=${query.location ?? ""}&category=${
          query.category ?? ""
        }&min_guest=${query.minGuest ?? ""}&bedrooms=${query.bedrooms ?? ""}`
      );
      return response;
    } catch (error) {
      throw new Error(error);
    }
  }
}

const axiosQueryWithAuth = new AxiosQueryWithAuth(axiosAuthenticated);

export { axiosQueryWithAuth };
