# GameScores.info

GameScores.info is a web based system, to keep track of game results and player
stats for small office games. It is initially designed to keep track of
[Foosball][foosball] scores.


## Google App Engine Edition

This branch contains a rewrite of the Java+Wicket+MySQL implementation in
[frankbille/ScoreBoard], into a [Google App Engine][gae] application, using
[Go language][golang] on the backend, to provide data to the front-end, through
a REST api. On the front-end [AngularJS][angularjs] is used, so most of the user
interface is rendered in the browser, to minimize the hit on the server.


### Why a rewrite?

Basically because the hassle of deploying Java applications, when you want
cheap/free hosting. This app is not meant for any commercial use but only for
fun/recreational spare time in office spaces.


### Why [App Engine][gae]?

I think the idea of having a platform that automatically scales up and down,
and handles the infrastructure for you, is very welcoming. It has some very
strict rules that must be followed, but when you do a lot of best practices
will also follow, making your application faster and more scalable.

It is still not my ideal backend, as there is still a lot of boilerplate to be
done, when all I really want is to enforce what people can do server-side.
Firebase is so far the best option, but doesn't allow server side code execution.


### Why [Go][golang]?

With the rewrite I had two criterias: Google App Engine and a `JavaScript` frontend.
A third criteria quickly popped up, which was that new App Engine instances should
start *fast*. With `Java` and `Python`, a cold boot is actually very slow, so the first
request from a user will take 10-30 seconds, depending on how much needing to be
done on startup. `Go` on the other hand, is a fully compiled binary, so a new instance
is available in <500ms (I haven't seen it so slow yet, normally around 50-100ms).

So by choosing `Go`, I can deliver a page with data, within a second even if no
instances is running. This is pretty amazing and is worth a try-out. I say this
because I haven't decided 100% that I will use Go. It is a newer language, with
syntax similar to `C` and with very few mature frameworks compared to `Java` or
`Python`.

### How to develop on the code?

This is still new technologies for me, so I haven't figured out what the best
stack is to develop this. I basically use a text editor ([Atom][atom]
in this case) to edit the `JavaScript` and `Go` files, and then run the Google
App Engine development server to build/host the files, so I can continuously
develop and test in a browser, without having to restart all the time.


## Todo

See the [issue list][issues] for details of what is needed to be done.


## License

Licensed under the [GNU General Public License, Version 3.0][license]


[foosball]: http://en.wikipedia.org/wiki/Table_football
[issues]: https://github.com/frankbille/gamescores.info/issues
[license]: http://www.gnu.org/licenses/gpl.html
[gae]: http://developers.google.com/appengine
[gaeinstall]: https://developers.google.com/appengine/downloads#Google_App_Engine_SDK_for_Go
[golang]: http://golang.org
[goinstalldoc]: http://golang.org/doc/install
[angularjs]: http://angularjs.org
[atom]: http://atom.io
