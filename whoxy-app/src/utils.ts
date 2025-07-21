export const forwardRequestTo = async (url: string, request: Request) => {
  return fetch(url, {
    method: request.method,
    headers: request.headers,
    body: request.body,
  });
};
