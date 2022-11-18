# Shared Recruiting Co. (SRC)

## Project Layout

### `/app`

The SRC web app ([https://sharedrecruiting.co](https://sharedrecruiting.co)) built via Sveltkit + Supbase and deployed with Vercel

### `/cloudfunctions`

The SRC Google Cloud Functions. The cloud functions are responsible for managing and reacting to user emails. Each cloud function is an independent, deployable  `go` module. 

### `/libs`

Shared `go` libraries
