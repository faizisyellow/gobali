import axios from "axios";

const axiosPublic = axios.create({
  baseURL: import.meta.env.VITE_BASE_URL_DEV,
  headers: {
    "Content-Type": "application/json",
  },
});

class AxiosQueryPublic {
  constructor(a) {
    this.axios = a;
  }

  async login(payload) {
    try {
      const response = await this?.axios?.post(
        "/v1/authentication/login",
        payload
      );

      return response;
    } catch (error) {
      throw new Error(error); // send the error to the caller
    }
  }

  async register(email, password, username) {
    try {
      const response = await this.axios.post("/v1/authentication/register", {
        email,
        password,
        username,
      });

      return response.data;
    } catch (error) {
      console.error("Unexpected error while login:", error);

      throw new Error(error.message);
    }
  }
}

const axiosQueryPublic = new AxiosQueryPublic(axiosPublic);

export { axiosQueryPublic };
