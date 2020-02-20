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




<h4>Database Tables + Explanations</h4>
<p> users:
    <ul>
        <li>id - serial - PRIMARY KEY</li>
        <li>username - VARCHAR(26)</li>
        <li>gender - BOOL</li>
        <li>age - INT</li>
        <li>password - VARCHAR(355)</li>
        <li>email - VARCHAR(55)</li>
    </ul>
</p>
<p>user_settings:
    <ul>
        <li>userid - INT - PRIMARY KEY</li>
        <li>bio - TEXT</li>
        <li>website - TEXT</li>
        <li>location - TEXT</li>
        <li>publicity - BOOL</li>
    </ul>
</p>
<p>posts:
    <ul>
        <li>postid - SERIAL - PRIMARY KEY</li>
        <li>userid - INT</li>
        <li>tags - TEXT</li>
        <li>caption - TEXT</li>
        <li>type - TEXT</li>
        <li>posted - date</li>
        <li>extension - TEXT</li>
    </ul>
</p>
<p>private_conversations:
    <ul>
        <li>convoID - SERIAL - PRIMARY KEY</li>
        <li>userOne - VARCHAR(26)</li>
        <li>userTwo - VARCHAR(26)</li>
        <li>created - DATE</li>
    </ul>
</p>
