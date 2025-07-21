const WHOXY_URL =
  Bun.env.WHOXY_URL ||
  (() => {
    throw new Error("Whoxy URL is required");
  })();

export const api = {
  createEventSink(url: string, attribution: any) {
    return fetch(`${WHOXY_URL}/api/event_sinks`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        url,
        attribution,
      }),
    });
  },
  deleteEventSink(id: string) {
    return fetch(`${WHOXY_URL}/api/event_sinks/${id}`, {
      method: "DELETE",
    });
  },
};
