import { Hono } from "hono";
import { zValidator } from "@hono/zod-validator";
import {
  CreateEventSinkPayloadSchema,
  EventSchema,
  EventSinkSchema,
} from "./schemas";
import { config } from "./config";

type Bindings = {
  event_sinks: KVNamespace;
};

const webhooks = new Hono<{ Bindings: Bindings }>().post(
  "/:service",
  zValidator("json", EventSchema),
  async (c) => {
    const service = c.req.param("service");
    const eventPayload = c.req.valid("json");

    c.req.raw.clone();

    const listResult = await c.env.event_sinks.list();
    const eventSinks = await Promise.all(
      listResult.keys
        .map((k) => k.name)
        .map((id) => c.env.event_sinks.get(id))
        .map((sink) => EventSinkSchema.safeParse(sink))
        .filter((sink) => sink.success)
        .map((sink) => sink.data)
    );

    for (const sink of eventSinks) {
      for (const f of config.eventFilters) {
        const url = f.shouldDeliverTo(eventPayload, service, sink);
        if (url) {
          console.log("delivering event to: " + url);
          fetch(url, {
            headers: c.req.raw.headers,
            body: JSON.stringify(eventPayload),
          }).catch((err) =>
            console.log("error forwarding to " + url + ": " + err)
          );
        }
      }
    }

    return c.json({ message: "Event processed successfully" }, 200);
  }
);

const api = new Hono<{ Bindings: Bindings }>()
  .post(
    "/event_sinks",
    zValidator("json", CreateEventSinkPayloadSchema),
    async (c) => {
      console.log(await c.req.text());
      const payload = c.req.valid("json");
      const id = crypto.randomUUID();
      await c.env.event_sinks.put(id, JSON.stringify({ ...payload, id }));
      return c.json({ id });
    }
  )
  .delete("/event_sinks/:id", async (c) => {
    const id = c.req.param("id");
    await c.env.event_sinks.delete(id);
    return c.json("ok");
  });

const app = new Hono<{ Bindings: Bindings }>()
  .get("/", (c) => c.json("whoxy"))
  .route("/api", api)
  .route("/webhooks", webhooks);

export default app;
