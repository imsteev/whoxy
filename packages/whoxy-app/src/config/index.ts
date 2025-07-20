import { EventSink } from "../schemas";

interface EventFilter {
  /**
   * Returns the url to deliver the event to, or false.
   */
  shouldDeliverTo(
    event: unknown,
    service: string,
    eventSink: EventSink
  ): string | false;
}

export const config: { eventFilters: EventFilter[] } = {
  // Add your config here.
  eventFilters: [],
};
