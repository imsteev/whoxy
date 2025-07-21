problem / motivation

i was working on a side project that involved webhooks. for the uninitiated, webhooks are a data delivery paradigm where an external service would attempt to deliver information about events that happened in that service to a set of destination urls that is typically configured by the application developers

one frustration i kept running into when testing webhooks during development was constantly having to update the local tunnel url in the webhook provider's dashboard of whitelisted url destinations

local tunnel tools generally randomize the url they generate any time you create a new tunnel (e.g, `https://early-pots-brake.loca.lt`, `https://small-rivers-write.loca.lt`)

that got me thinking: wouldn't it be cool to just set a single, static URL for the webhook destination, and have the endpoint handling that URL figure out how to fan-out to dev machines? from a DevEx standpoint this is a productivity win, at the cost of a bit of infrastructure

using cloudflare primitives, we can build something fairly reliable and available:

architecture
<img width="751" height="531" alt="Screenshot 2025-07-20 at 11 39 21 PM" src="https://github.com/user-attachments/assets/4dfaa1ac-0ef8-4d00-a121-0b0aead3ae69" />
