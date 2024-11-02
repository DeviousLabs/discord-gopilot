# Discord Gopilot
Meet Gopilot: a cutting-edge bot crafted in Go, designed to integrate AI models seamlessly into Discord channels. This bot allows users to effortlessly utilize AI tools to generate code snippets, solve algorithmic problems, and assist with project tasks, directly within their Discord communities. Enhance your coding and collaboration efficiency with Discord Gopilot, bringing advanced AI expertise right to your fingertips.

### Docker Build Example
```
docker build -t gopilot:latest-amd64
```

### Docker Run Example
```
docker run -e CLOUDFLARE_ACCOUNT_ID='ACCOUNT_ID' CLOUDFLARE_API_TOKEN='TOKEN' DISCORD_TOKEN='TOKEN' MODEL='MODEL_ID' PERSONA='PERSONA' gopilot:latest-amd64
```

### Dock-Compose.yml Example
```
version: '3.7'
  services:
    gopilot:
      image: "gopilot:latest-amd64"
      container_name: "gopilot"
      environment:
        CLOUDFLARE_ACCOUNT_ID: "ACCOUNT_ID"
        CLOUDFLARE_API_TOKEN: "TOKEN"
        DISCORD_TOKEN: "TOKEN"
        MODEL: "MODEL_ID"
        PERSONA: "PERSONA"
      restart: "unless-stopped"
```
### Models Available
|  ID  | Model Name                        | Available |
|------|-----------------------------------|-----------|
|1     | llama-3.1-70b-instruct            | ✅        |
|2     | deepseek-coder-6.7b-instruct-awq  | ✅        |
|3     | gemma-7b-it                       | ✅        |
|4     | mistral-7b-instruct-v0.2          | ✅        |
|5     | qwen1.5-14b-chat-awq              | ✅        |
|6     | phi-2                             | ✅        |
|7     | stable-diffusion-xl-base-1.0      | ❌        |


### Persona's Available
- [x] default
- [x] developer
