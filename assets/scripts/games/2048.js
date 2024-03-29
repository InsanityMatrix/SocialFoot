var gameboard;
var canvas = document.getElementById('game');
var context = canvas.getContext('2d');
var score = 0;
var gameState = 0;
var borderSize = 20;
var cellSize = 100;
var userid = 0;
var leaderboard = null;
// gameState
// 0 = Playing
// 1 = Lost
function setUserID(uid) {
  userid = uid;
}
function populateScores() {
  $.ajax({
    url: '/games/2048/scores',
    success: function(data) {
      leaderboard = data;
    }
  });
}
function startGame() {
  gameboard = [
    [0,0,0,0],
    [0,0,0,0],
    [0,0,0,0],
    [0,0,0,0]
  ];
  let cell = chooseEmptyCell();
  gameboard[cell[0]][cell[1]] = 2;
  cell = chooseEmptyCell();
  score = 0;
  gameboard[cell[0]][cell[1]] = 2;
  gameState = 0;
}

function chooseEmptyCell() {
  let tries = 0;
  while(tries < 125) {
    let rRow = Math.floor(Math.random() * 4);
    let rCol = Math.floor(Math.random() * 4);

    let cell = gameboard[rRow][rCol];
    if(cell === 0) {
      return [rRow, rCol];
    }
    tries++;
  }
  for(var i = 0; i < 4; i++) {
    for(var j = 0; j < 4; j++) {
      if(gameboard[i][j] === 0) {
        return [i, j];
      }
    }
  }
  return [-1, -1];
}
function moveLeft() {
  //Start combine from left to right
  for(var i = 0; i < 4; i++) {
    for(var j = 0; j < 3; j++) {
      if(gameboard[i][j] === gameboard[i][j + 1]) {
        gameboard[i][j] *= 2;
        gameboard[i][j + 1] = 0;
        score += gameboard[i][j];
      } else if(gameboard[i][j + 1] == 0 &&
                j + 2 < 4 &&
                gameboard[i][j + 2] === gameboard[i][j]) {
        gameboard[i][j] *= 2;
        gameboard[i][j + 2] = 0;
        score += gameboard[i][j];
      } else if(gameboard[i][j + 1] == 0 &&
                j + 2 < 4 &&
                gameboard[i][j + 2] == 0 &&
                j + 3 < 4 &&
                gameboard[i][j + 3] === gameboard[i][j]) {
        gameboard[i][j] *= 2;
        gameboard[i][j + 3] = 0;
        score += gameboard[i][j];
      }
    }
  }
  for(var i = 0; i < 4; i++) {
    for(var j = 0; j < 3; j++) {
      if(gameboard[i][j + 1] != 0 && gameboard[i][j] === 0) {
        //Shift this cell left
        gameboard[i][j] = gameboard[i][j + 1];
        gameboard[i][j + 1] = 0;
      } else if(gameboard[i][j+1] === 0 &&
                j + 2 <= 3 &&
                gameboard[i][j+2] != 0 &&
                gameboard[i][j] === 0) {
        gameboard[i][j] = gameboard[i][j+2];
        gameboard[i][j+2] = 0;
      } else if(gameboard[i][j+1] === 0 &&
                j + 2 <= 3 &&
                gameboard[i][j+2] === 0 &&
                j + 3 <= 3 &&
                gameboard[i][j+3] != 0 &&
                gameboard[i][j] === 0) {
        gameboard[i][j] = gameboard[i][j+3];
        gameboard[i][j+3] = 0;
      }
    }
  }
  spawnNewCell();
}
function moveRight() {
  for(var i = 0; i < 4; i++) {
    for(var j = 3; j > 0; j--) {
      if(gameboard[i][j] === gameboard[i][j - 1]) {
        gameboard[i][j] *= 2;
        gameboard[i][j - 1] = 0;
        score += gameboard[i][j];
      }else if(gameboard[i][j - 1] == 0 &&
                j - 2 >= 0 &&
                gameboard[i][j - 2] === gameboard[i][j]) {
        gameboard[i][j] *= 2;
        gameboard[i][j - 2] = 0;
        score += gameboard[i][j];
      } else if(gameboard[i][j - 1] == 0 &&
                j - 2 >= 0 &&
                gameboard[i][j - 2] == 0 &&
                j - 3 >= 0 &&
                gameboard[i][j - 3] === gameboard[i][j]) {
        gameboard[i][j] *= 2;
        gameboard[i][j - 3] = 0;
        score += gameboard[i][j];
      }
    }
  }
  for(var i = 0; i < 4; i++) {
    for(var j = 3; j > 0; j--) {
      if(gameboard[i][j - 1] != 0 && gameboard[i][j] === 0) {
        //Shift this cell right
        gameboard[i][j] = gameboard[i][j-1];
        gameboard[i][j-1] = 0;
      } else if(gameboard[i][j-1] === 0 &&
                j - 2 >= 0 &&
                gameboard[i][j-2] != 0 &&
                gameboard[i][j] === 0) {
        gameboard[i][j] = gameboard[i][j-2];
        gameboard[i][j-2] = 0;
      } else if(gameboard[i][j-1] === 0 &&
                j - 2 >= 0 &&
                gameboard[i][j-2] === 0 &&
                j - 3 >= 0 &&
                gameboard[i][j-3] != 0 &&
                gameboard[i][j] === 0) {
        gameboard[i][j] = gameboard[i][j-3];
        gameboard[i][j-3] = 0;
      }
    }
  }
  spawnNewCell();
}
function moveUp() {
  for(var j = 0; j < 4; j++) {
    for(var i = 0; i < 3; i++) {
      if(gameboard[i][j] === gameboard[i + 1][j]) {
        gameboard[i][j] *= 2;
        gameboard[i+1][j] = 0;
        score += gameboard[i][j];
      } else if(gameboard[i + 1][j] == 0 &&
                i + 2 < 4 &&
                gameboard[i + 2][j] === gameboard[i][j]) {
        gameboard[i][j] *= 2;
        gameboard[i + 2][j] = 0;
        score += gameboard[i][j];
      } else if(gameboard[i + 1][j] == 0 &&
                i + 2 < 4 &&
                gameboard[i + 2][j] == 0 &&
                i + 3 < 4 &&
                gameboard[i + 3][j] === gameboard[i][j]) {
        gameboard[i][j] *= 2;
        gameboard[i+ 3][j] = 0;
        score += gameboard[i][j];
      }
    }
  }
  for(var j = 0; j < 4; j++) {
    for(var i = 0; i < 3; i++) {
      if(gameboard[i+1][j] != 0 && gameboard[i][j] === 0) {
        gameboard[i][j] = gameboard[i+1][j];
        gameboard[i+1][j] = 0;
      } else if(gameboard[i+1][j] === 0 &&
                i + 2 <= 3 &&
                gameboard[i+2][j] != 0 &&
                gameboard[i][j] === 0) {
        gameboard[i][j] = gameboard[i+2][j];
        gameboard[i+2][j] = 0;
      } else if(gameboard[i+1][j] === 0 &&
                i + 2 <= 3 &&
                gameboard[i+2][j] === 0 &&
                i + 3 <= 3 &&
                gameboard[i+3][j] != 0 &&
                gameboard[i][j] === 0) {
        gameboard[i][j] = gameboard[i+3][j];
        gameboard[i+3][j] = 0;
      }
    }
  }
  spawnNewCell();
}
function moveDown() {
  for(var j = 0; j < 4; j++) {
    for(var i = 3; i > 0; i--) {
      if(gameboard[i][j] === gameboard[i - 1][j]) {
        gameboard[i][j] *= 2;
        gameboard[i-1][j] = 0;
        score += gameboard[i][j];
      } else if(gameboard[i - 1][j] == 0 &&
                i - 2 >= 0 &&
                gameboard[i - 2][j] === gameboard[i][j]) {
        gameboard[i][j] *= 2;
        gameboard[i - 2][j] = 0;
        score += gameboard[i][j];
      } else if(gameboard[i - 1][j] == 0 &&
                i - 2 >= 4 &&
                gameboard[i - 2][j] == 0 &&
                i - 3 >= 0 &&
                gameboard[i - 3][j] === gameboard[i][j]) {
        gameboard[i][j] *= 2;
        gameboard[i-3][j] = 0;
        score += gameboard[i][j];
      }
    }
  }
  for(var j = 0; j < 4; j++) {
    for(var i = 3; i > 0; i--) {
      if(gameboard[i-1][j] != 0 && gameboard[i][j] === 0) {
        gameboard[i][j] = gameboard[i-1][j];
        gameboard[i-1][j] = 0;
      } else if(gameboard[i-1][j] === 0 &&
                i - 2 >= 0 &&
                gameboard[i-2][j] != 0 &&
                gameboard[i][j] === 0) {
        gameboard[i][j] = gameboard[i-2][j];
        gameboard[i-2][j] = 0;
      } else if(gameboard[i-1][j] === 0 &&
                i - 2 >= 0 &&
                gameboard[i-2][j] === 0 &&
                i - 3 >= 0 &&
                gameboard[i-3][j] != 0 &&
                gameboard[i][j] === 0) {
        gameboard[i][j] = gameboard[i-3][j];
        gameboard[i-3][j] = 0;
      }
    }
  }
  spawnNewCell();
}
function spawnNewCell() {
  let value = 2;
  if(Math.random() > 0.7) {
    value = 4;
  }
  cell = chooseEmptyCell();
  if (cell[0] === -1) {
    gameState = 1;
    $.ajax({
      url: '/games/2048/update',
      method: 'POST',
      data: {
        "userid": userid,
        "score": score
      },
      success: populateScores
    });
    return;
  }
  gameboard[cell[0]][cell[1]] = value;
}
document.addEventListener('keydown', function(e) {
  //LEFT
  if(e.which === 37) {
    moveLeft();
  }
  //RIGHT
  if(e.which === 39) {
    moveRight();
  }
  //UP
  if(e.which === 38) {
    moveUp();
  }
  //DOWN
  if(e.which === 40) {
    moveDown();
  }

  if(e.which === 32) {
    if(gameState === 1) {
      startGame()
    }
  }
});
//Precondition width and height = 500
function loop() {
  requestAnimationFrame(loop);
  //60 FPS 2048 Babyyy
  $("#score").html(score);
  if(gameState === 0) {
    context.clearRect(0,0,canvas.width, canvas.height);
    //Borders = 20
    context.fillStyle = "#c5c8d3";
    //Vertical lines
    context.fillRect(0,0,borderSize,canvas.height);
    context.fillRect((cellSize) + (borderSize),0,borderSize,canvas.height);
    context.fillRect((cellSize * 2) + (borderSize * 2), 0, borderSize, canvas.height);
    context.fillRect((cellSize * 3) + (borderSize * 3),0,borderSize,canvas.height);
    context.fillRect((cellSize * 4) + (borderSize * 4),0,borderSize,canvas.height);
    //Horizontal Lines
    context.fillRect(0,0,canvas.width, borderSize);
    context.fillRect(0,(cellSize) + (borderSize),canvas.width, borderSize);
    context.fillRect(0,(cellSize * 2) + (borderSize * 2),canvas.width, borderSize);
    context.fillRect(0,(cellSize * 3) + (borderSize * 3),canvas.width, borderSize);
    context.fillRect(0,(cellSize * 4) + (borderSize * 4),canvas.width,borderSize);


    //Draw Game Grid
    drawGrid(borderSize,borderSize, gameboard[0][0]);
    drawGrid((cellSize) + (borderSize * 2),borderSize, gameboard[0][1]);
    drawGrid((cellSize * 2) + (borderSize * 3),borderSize, gameboard[0][2]);
    drawGrid((cellSize * 3) + (borderSize * 4),borderSize, gameboard[0][3]);
    drawGrid(borderSize,(cellSize) + (borderSize * 2), gameboard[1][0]);
    drawGrid((cellSize) + (borderSize * 2),(cellSize) + (borderSize * 2), gameboard[1][1]);
    drawGrid((cellSize * 2) + (borderSize * 3),(cellSize) + (borderSize * 2), gameboard[1][2]);
    drawGrid((cellSize * 3) + (borderSize * 4),(cellSize) + (borderSize * 2), gameboard[1][3]);
    drawGrid(borderSize,(cellSize * 2) + (borderSize * 3), gameboard[2][0]);
    drawGrid((cellSize) + (borderSize * 2),(cellSize * 2) + (borderSize * 3), gameboard[2][1]);
    drawGrid((cellSize * 2) + (borderSize * 3),(cellSize * 2) + (borderSize * 3), gameboard[2][2]);
    drawGrid((cellSize * 3) + (borderSize * 4),(cellSize * 2) + (borderSize * 3), gameboard[2][3]);
    drawGrid(borderSize,(cellSize * 3) + (borderSize * 4), gameboard[3][0]);
    drawGrid((cellSize) + (borderSize * 2),(cellSize * 3) + (borderSize * 4), gameboard[3][1]);
    drawGrid((cellSize * 2) + (borderSize * 3),(cellSize * 3) + (borderSize * 4), gameboard[3][2]);
    drawGrid((cellSize * 3) + (borderSize * 4),(cellSize * 3) + (borderSize * 4), gameboard[3][3]);
    return;
  }
  //Lost
  let data = leaderboard;
  context.clearRect(0,0,canvas.width, canvas.height);
  context.fillStyle = "#b58251";
  context.fillRect(0,0,canvas.width, 40);
  context.font = "20px Comic Sans MS";
  context.textAlign = "center";
  context.fillStyle = "white";
  context.fillText("Your Score: " + score, canvas.width/2, 30);
  context.fillStyle = "#5fc4b6";
  context.fillRect(0,40,canvas.width, 130);
  context.fillStyle = "black";
  context.fillText("1. " + data[0].Username, canvas.width/4 * 1.5, 60);
  context.fillText("" + data[0].Score, canvas.width/4 * 3, 60)
  if(data[1] != undefined) {
    context.fillText("2. " + data[1].Username, canvas.width/4 * 1.5, 90)
    context.fillText("" + data[1].Score, canvas.width/4 * 3, 90)
    if(data[2] != undefined) {
      context.fillText("3. " + data[2].Username, canvas.width/4 * 1.5, 120)
      context.fillText("" + data[2].Score, canvas.width/4 * 3, 120)
      if(data[3] != undefined) {
        context.fillText("4. " + data[3].Username, canvas.width/4 * 1.5, 150)
        context.fillText("" + data[3].Score, canvas.width/4 * 3, 150)
        if(data[4] != undefined) {
          context.fillText("5. " + data[4].Username, canvas.width/4 * 1.5, 180)
          context.fillText("" + data[4].Score, canvas.width/4 * 3, 180)
        }
      }
    }
  }
  context.font = "15px Comic Sans MS";
  context.fillText("(Press Space or Tap to Continue)", canvas.width/2, canvas.height - 10);
}
function drawGrid(x,y, value) {
  context.fillStyle = "#eaebef";
  let fontColor = "black";
  if(value >= 1028) {
    fontColor = "white";
    context.fillStyle = "#e68a00";
  } else if(value >= 256) {
    fontColor = "white";
    context.fillStyle = "#ff4d4d";
  } else if(value >= 64) {
    fontColor = "white";
    context.fillStyle = "#ff8080";
  } else if(value >= 8) {
    context.fillStyle = "#ff9900";
    fontColor = "white";
  } else if (value > 0) {
    context.fillStyle = "#ffcc80";
  }
  context.fillRect(x,y,cellSize,cellSize);
  if(value != 0) {
    context.font = "30px Comic Sans MS";
    context.fillStyle = fontColor;
    context.textAlign = "center";
    context.fillText(value, x + (cellSize/2), y + (cellSize/2));
  }
}
//Prevent page scrolling
document.touchmove = function(e)
{
    e.preventDefault();
};
//Add touches for phones
document.addEventListener('touchstart', handleTouchStart, false);
document.addEventListener('touchmove', handleTouchMove, false);
var touchX = null;
var touchY = null;
function getTouches(e) {
  return e.touches || e.originalEvent.touches;
}
function handleTouchStart(e) {
  const firstTouch = getTouches(e)[0];
  touchX = firstTouch.clientX;
  touchY = firstTouch.clientY;

  setTimeout(function() {
    if(touchX != null  && touchY != null) {
      if(gameState === 1) {
        startGame()
      }
    }
  }, 150);
}
function handleTouchMove(e) {
  if (! touchX || ! touchY) {
    return;
  }

  let xUp = e.touches[0].clientX;
  let yUp = e.touches[0].clientY;

  let xDiff = touchX - xUp;
  let yDiff = touchY - yUp;

  if(Math.abs(xDiff) > Math.abs(yDiff)) {
    if (xDiff > 0) {
      //Right to Left swipe
      moveLeft();
    } else {
      //Left to Right Swipe
      moveRight();
    }
  } else {
    if ( yDiff > 0) {
      //Bottom to Top
      moveUp();
    } else {
      //Top to Bottom
      moveDown();
    }
  }
  touchX = null;
  touchY = null;
}
$(document).ready(function() {
  if($(window).width() < 500) {
    context.canvas.width = window.innerWidth;
    context.canvas.height = context.canvas.width;
  } else {
    context.canvas.width = 500;
    context.canvas.height = 500;
  }
  borderSize = context.canvas.width / 25;
  cellSize = context.canvas.width / 5;
});
startGame();
requestAnimationFrame(loop);
