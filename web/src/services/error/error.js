class ErrorData extends Error {
  constructor(message, status, type, details) {
    super(message);
    this.status = status;
    this.type = type;
    this.details = details;
  }

}

export { ErrorData };
