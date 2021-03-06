/**
 * @file static.go
 * @author Mikhail Klementyev jollheef<AT>riseup.net
 * @license GNU AGPLv3
 * @date September, 2015
 * @brief non-dynamic html results
 *
 * Generate static html page with scoreboard
 */

package scoreboard

import (
	"fmt"
	"net/http"
)

func staticScoreboard(w http.ResponseWriter, r *http.Request) {

	// TODO: rewrite with templates

	fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>IBST.PSU CTF Scoreboard</title>

    <link rel="stylesheet" href="css/bootstrap.min.css">
    <link rel="stylesheet" href="css/style.css">

    <script type="text/javascript">
      var scoreboard = new WebSocket("ws://" + location.host + "/scoreboard");

      scoreboard.onmessage = function(e) {
        document.getElementById('scoreboard-table').innerHTML = e.data
      }

      var info = new WebSocket("ws://" + location.host + "/info");

      info.onmessage = function(e) {
        document.getElementById('info').innerHTML = e.data
      }
    </script>
  </head>
  <body class="full">
    <ul class="nav nav-tabs">
      <li class="active">
        <a href="#">Scoreboard</a>
      </li>
      <!-- <li><a href="advisory.html">Advisory</a></li> -->
      <li><a href="info.html">Information</a></li>
    </ul>
    <div class="page-header"><center><h1>IBST.PSU CTF Scoreboard</h1></center></div>
    <div style="padding: 15px;">
      <div id="info">%s</div>
      <br>
      <table id="scoreboard-table" class="table table-hover">%s</table>
      <script src="js/bootstrap.min.js"></script>
    </div>
  </body>
</html>`, getInfo(), currentResult)
}
