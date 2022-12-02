
A demo showcasing how to build a real-time multi-user experience with minimal custom javascript, using the excellent [htmx](https://htmx.org/) and the built-in [sse extension](https://htmx.org/attributes/hx-sse/).

This is a quiz you can run for a group of friends, where everyone's position in the quiz is synchronized, and all players see any guesses made by any player in a chat-like sidebar.

Technical users that can view source and make HTTP requests on their own are still unable to cheat :tm: since the state is on the server-side.

The quiz runner can reveal the answer and/or take everyone to the next question.

# Setup

Set a `SESSION_SECRET` environment variable to authenticate users

Use `go run cmd/token/main.go <username>` to get access tokens for every user you want to add to the quiz.

Edit `questions.go` if you want to change or add questions. Note that since it's code you can edit the `check` function for each question, to allow more creative answers if you wish.

# Starting the server

```
go run .
```
	
You can use [ngrok](https://ngrok.com/) or a similar tool to generate a public facing base URL for your quiz.

Then you can generate a unique access URL for every user, by appending the access token you generated above. For example:

https://2196-2403-etc.au.ngrok.io?access_token=tokenfromsetupstep

# Running the quiz

Go to `http://localhost:5000` and you will see the same interface as everyone else, but only you have the `Reveal` and `Next` buttons.

I recommend dialing into a video chat or something with your friends or co-workers, and then using those buttons as you see fit.
