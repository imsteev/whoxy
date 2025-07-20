import { z } from "zod";

// Downstream implementations should define their own event schema.
export const EventSchema = z.any();

export const CreateEventSinkPayloadSchema = z.object({
  url: z.string(),
  attribution: z.record(z.string(), z.any()),
});

export const EventSinkSchema = z.object({
  id: z.string(),
  url: z.string(),
  attribution: z.record(z.string(), z.any()),
});

export type Event = z.infer<typeof EventSchema>;
export type EventSink = z.infer<typeof EventSinkSchema>;
export type CreateEventSink = z.infer<typeof CreateEventSinkPayloadSchema>;
