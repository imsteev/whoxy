import localtunnel from "localtunnel";
import { api } from "./api";

const TUNNEL_PORT =
  Bun.env.TUNNEL_PORT ||
  (() => {
    throw new Error("Whoxy URL is required");
  })();

async function main() {
  const lt = await localtunnel({ port: parseInt(TUNNEL_PORT) });
  console.log("Tunnel at: " + lt.url);
  const resp = await api.createEventSink(lt.url, { ipv4: [] });
  const sinkId = ((await resp.json()) as { id: string }).id;
  await api.deleteEventSink(sinkId);
}

main();
