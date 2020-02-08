package main

import (
	"io/ioutil"
	"path"
	"strings"
)

type Html interface {
	initIndexHTML() string
	initProfileHTML() string
	initProfileSettingsHTML() string
}

func initIndexHTML() string {
//	header := "<!DOCTYPE html><html lang='en'>    <head>        <title>Home</title>		<meta charset='utf-8'><meta name='description' content='{MF_PLUGIN_SETTING:HOME_DESCRIPTION}'/><meta name='viewport' content='width=device-width,initial-scale=1,maximum-scale=1,minimum-scale=1'/><meta name='apple-mobile-web-app-capable' content='yes'><meta name='mobile-web-app-capable' content='yes'>		<!-- Bootstrap CSS -->    <link rel='stylesheet' href='https://maxcdn.bootstrapcdn.com/bootstrap/3.4.1/css/bootstrap.min.css'>    <script src='https://ajax.googleapis.com/ajax/libs/jquery/3.4.1/jquery.min.js'></script>    <script src='https://maxcdn.bootstrapcdn.com/bootstrap/3.4.1/js/bootstrap.min.js'></script>    	<script src='https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js' integrity='sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM' crossorigin='anonymous'></script>	<script src='https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js' integrity='sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1' crossorigin='anonymous'></script> <style>html,body {  margin: 0;  padding: 0;}  </style>  <link rel='stylesheet' href='./assets/css/SocialFoot.css'>    </head>    <body id='body'>        <nav class='navbar navbar-inverse bg-primary'>            <div class='container-fluid'>                <div class='navbar-header'>                    <a class='navbar-brand' href='/live'>SocialFoot</a>               </div>                <ul class='nav navbar-nav'>                  <li class='active'><a href='/live'>Home</a></li>                  <li class='dropdown'>                    <a class='dropdown-toggle' data-toggle='dropdown' href='#'>Profile                     <span class='caret'></span></a>                    <ul class='dropdown-menu'>                      <li><a href='/live/profile'>Settings</a></li>                      <li><a href='#'>Messages</a></li>                      <li><a href='#'>Your Feed</a></li>                   </ul>                  </li>                  <li><a href='#'>Search</a></li>                 <li><a href='#'>Report</a></li>                </ul>                <ul class='nav navbar-nav navbar-right'>                  <li><a href='/live/profile'><span class='glyphicon glyphicon-user'></span> {{.username}}</a></li>                </ul>           </div>        </nav>"
//	body := "<div class='net container'>		<div class='col-md-9 col-sm-12'>			<!-- Feed Here -->			<div class='feed'>				<center>					<h2 class='bg-primary'>Your Feed</h2>				</center>				<!-- TODO: Insert data for feed -->			</div>		</div>	</div>"
	//bottom := "<script>      var body = document.getElementById('body');      body.removeChild(body.childNodes[0]);   </script>		<script src='https://code.jquery.com/jquery-3.3.1.slim.min.js' integrity='sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo' crossorigin='anonymous'></script>    </body></html>"
	//ht := header + body + bottom
	fp := path.Join("templates","index.html")
	content, err := ioutil.ReadFile(fp)
	if err != nil {
		return ""
	}
	ht := string(content)
	return ht
}

func initProfileHTML() string {
	header := "<!DOCTYPE html><html lang='en'>    <head>        <title>Home</title>		<meta charset='utf-8'>		<meta name='viewport' content='width=device-width, initial-scale=1, shrink-to-fit=no'>    <meta http-equiv='Content-Type' content='text/html; charset=utf-8' />    <meta name='apple-mobile-web-app-capable' content='yes'>    <meta name='mobile-web-app-capable' content='yes'>		<!-- Bootstrap CSS -->    <link rel='stylesheet' href='https://maxcdn.bootstrapcdn.com/bootstrap/3.4.1/css/bootstrap.min.css'>    <script src='https://ajax.googleapis.com/ajax/libs/jquery/3.4.1/jquery.min.js'></script>    <script src='https://maxcdn.bootstrapcdn.com/bootstrap/3.4.1/js/bootstrap.min.js'></script>    	<script src='https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js' integrity='sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM' crossorigin='anonymous'></script>	<script src='https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js' integrity='sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1' crossorigin='anonymous'></script>  <style>html,body {  margin: 0;  padding: 0;}  </style>  <link rel='stylesheet' href='../assets/css/SocialFoot.css'>    </head>    <body id='body'>        <nav class='navbar navbar-inverse bg-primary'>            <div class='container-fluid'>                <div class='navbar-header'>                    <a class='navbar-brand' href='/live'>SocialFoot</a>                </div>                <ul class='nav navbar-nav'>          <li class='active'><a href='/live'>Home</a></li>                  <li class='dropdown'>                    <a class='dropdown-toggle' data-toggle='dropdown' href='#'>Profile                    <span class='caret'></span></a>                    <ul class='dropdown-menu'>                      <li><a href='#'>Settings</a></li>                      <li><a href='#'>Messages</a></li>                     <li><a href='#'>Your Feed</a></li>                   </ul>                  </li>                  <li><a href='#'>Search</a></li>                  <li><a href='#'>Report</a></li>                </ul>                <ul class='nav navbar-nav navbar-right'>                 <li><a href='#'><span class='glyphicon glyphicon-user'></span> {{.username}}</a></li>                </ul>            </div>        </nav>   <script>      var body = document.getElementById('body');      body.removeChild(body.childNodes[0]);   </script>"
	body := "	<div class='net container'>    <nav aria-label='breadcrumb'>        <ol class='breadcrumb'>          <li class='breadcrumb-item'><a href='/live'>Home</a></li>          <li class='breadcrumb-item active'><a href='http://socialfoot.herokuapp.com/live/profile'>Profile</a></li>        </ol>    </nav>    <form class='container paper' method='post' action='/live/profile/settings'>      <center>      <h3>Enter your password, {{.username}}</h3>      <div class='row justify-content-md-center'>          <div class='col-sm-9 col-xs-10' style='display:block;margin: 5px auto;'>              <input type='password' class='form-control' name='password' id='password' placeholder='Password'>          </div>          <div class='col-sm-3 col-xs-10' style='display:block;margin: 5px auto;'>           <input type='submit' style='height:100%;display:block;width:100%;height:100%;padding:5px' class='btn-primary' name='submit' value='Verify'></div>        </div>        </center>      </form>    </div>"

	bottom := "<script src='https://code.jquery.com/jquery-3.3.1.slim.min.js' integrity='sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo' crossorigin='anonymous'></script>    </body></html>"
	ht := header + body + bottom
	return ht
}

func initProfileSettingsHTML() string {
	fp := path.Join("templates","profile.html")
	content, err := ioutil.ReadFile(fp)
	if err != nil {
		return ""
	}
	text := string(content)
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&lt;","<")
	return text
}
