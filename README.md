# three-thirteen

A simple game card to try out golang.

Current status: work in progress.

To setup a server take a look at the `BOOTSRAP` script and run it if you want.

Then you should call something like:
```
go build and ./three-thirteen
```

Then in your browser you:
```
https://localhost:8080/3-13/?create=test&players=a,b,c
https://localhost:8080/static/test.htm
https://localhost:8080/3-13/test/a/?marshal=0
https://localhost:8080/static/test.htm
https://localhost:8080/3-13/test/
https://localhost:8080/3-13/?list
```
and some other calls.

After marshal call above (which is "DEAL"), you can try some GUI interactions.
Go to `https://localhost:8080/static/test.htm?name=c` (or replace c with current player).
You can drag a card from pile or deck to your hand.
Refresh after this, and drag a card from hand to the pile.
Refresh, and a new prompt message should be visible.

