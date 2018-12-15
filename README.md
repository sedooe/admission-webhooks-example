# admission-webhooks-example

Injects **organization** and **project** labels with **MutatingWebhook** and validates that there is no change on these labels with **ValidatingWebhook**.

Use [ngrok](https://ngrok.com/) for the ease of development.

- `dep ensure`
- `ngrok http 8080`
- Replace the ngrok url in the resource files with yours.
- `kubectl apply -f .`
- `go build && go run *.go`
