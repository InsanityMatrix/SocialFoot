# SocialFoot
A Web App written with full backend on Golang.

Static Content in /assets/ (html/css/javascript)
<h1>Pages (Not including assets):</h1>
<ul>
    <li><h4>/live</h4> - Logged in homepage</li>
    <li><h4>/live/profile</h4> - Verify before editing profile settings</li>
    <li><h4>/live/profile/settings</h4> - Edit profile settings</li>
</ul>
<h2>main.go</h2>
<p>This file sets up the router and starts a database connection, handling the requests and some functions</p>

<h2>store.go</h2>
<h5>This is the database interaction file.</h5>

<h2>tools.go</h2>
<p>Contains tools like password hashing, etc.</p>
