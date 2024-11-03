# Discord Gopilot

Gopilot is a Discord bot written in Go that integrates AI models via the Cloudflare AI API. It allows users to generate code snippets, solve algorithmic problems, and assist with project tasks directly within their Discord channels, making AI tools easily accessible to developers and project teams.

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

### Contributing

Contributions are welcome! Please follow these security best practices:

Pull Requests: Submit changes through pull requests for code review.

Security Testing: Run all tests, including static analysis and dependency checks.

Secrets Management: Do not hardcode any sensitive information; use environment variables.

### License

This project is open source under the MIT License.
