import localtunnel from "localtunnel";
import { api } from "./api";

const TUNNEL_PORT =
  Bun.env.TUNNEL_PORT ||
  (() => {
    throw new Error("Whoxy URL is required");
  })();

let SINK_ID = "";

// todo: why does process.on("SIGINT", ...) case this code to fire twice?
process.once("SIGINT", async () => {
  if (SINK_ID) {
    console.log("cleaning up event sink");
    await api.deleteEventSink(SINK_ID);
  }
  process.exit();
});

async function main() {
  const lt = await localtunnel({ port: parseInt(TUNNEL_PORT) });
  const resp = await api.createEventSink(lt.url, { ipv4: [] });
  SINK_ID = ((await resp.json()) as { id: string }).id;
  console.log("Tunnel at: " + lt.url);
  console.log("Forwarding port: " + TUNNEL_PORT);
}

main();
